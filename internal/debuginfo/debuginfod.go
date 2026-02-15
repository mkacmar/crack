package debuginfo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"log/slog"

	"golang.org/x/sync/singleflight"

	"github.com/mkacmar/crack/internal/version"
)

const (
	DefaultServerURL   = "https://debuginfod.elfutils.org"
	DefaultTimeout     = 30 * time.Second
	DefaultRetries     = 3
	DefaultMaxFileSize = 2 * 1024 * 1024 * 1024 // 2GB
)

type nonRetryableError struct{ err error }

func (e *nonRetryableError) Error() string { return e.err.Error() }
func (e *nonRetryableError) Unwrap() error { return e.err }

func isNonRetryable(statusCode int) bool {
	switch statusCode {
	case 400, 403, 404, 405, 410:
		return true
	default:
		return false
	}
}

func DefaultCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user cache directory: %w", err)
	}
	return filepath.Join(cacheDir, "crack", "debuginfo"), nil
}

func userAgent() string {
	return "crack/" + version.Version + " (Compiler Hardening Checker; +https://github.com/mkacmar/crack)"
}

type Client struct {
	serverURLs  []string
	cacheDir    string
	httpClient  *http.Client
	maxRetries  int
	maxFileSize int64
	logger      *slog.Logger
	downloads   singleflight.Group
}

type Options struct {
	ServerURLs  []string
	CacheDir    string
	Timeout     time.Duration
	MaxRetries  int
	MaxFileSize int64
	Logger      *slog.Logger
}

func NewClient(opts Options) (*Client, error) {
	if len(opts.ServerURLs) == 0 {
		return nil, fmt.Errorf("no debuginfod servers configured")
	}

	cacheDir := opts.CacheDir
	if cacheDir == "" {
		var err error
		cacheDir, err = DefaultCacheDir()
		if err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(cacheDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	maxRetries := opts.MaxRetries
	if maxRetries < 1 {
		maxRetries = DefaultRetries
	}

	maxFileSize := opts.MaxFileSize
	if maxFileSize <= 0 {
		maxFileSize = DefaultMaxFileSize
	}

	client := &Client{
		serverURLs:  opts.ServerURLs,
		cacheDir:    cacheDir,
		httpClient:  &http.Client{Timeout: timeout},
		logger:      opts.Logger.With(slog.String("component", "debuginfod")),
		maxRetries:  maxRetries,
		maxFileSize: maxFileSize,
	}

	return client, nil
}

func (c *Client) FetchDebugInfo(ctx context.Context, buildID string) (string, error) {
	if buildID == "" {
		return "", fmt.Errorf("build-id is empty")
	}

	result, err, _ := c.downloads.Do(buildID, func() (interface{}, error) {
		return c.fetchDebugInfo(ctx, buildID)
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (c *Client) fetchDebugInfo(ctx context.Context, buildID string) (string, error) {
	cachedPath := c.getCachePath(buildID)

	if _, err := os.Stat(cachedPath); err == nil {
		c.logger.Debug("using cached debug symbols", slog.String("build_id", buildID))
		return cachedPath, nil
	}

	for _, serverURL := range c.serverURLs {
		path, err := c.fetchFromServerWithRetry(ctx, serverURL, buildID, cachedPath)
		if err == nil {
			c.logger.Debug("fetched debug symbols", slog.String("build_id", buildID), slog.String("server", serverURL))
			return path, nil
		}
		c.logger.Debug("server failed", slog.String("server", serverURL), slog.Any("error", err))
	}

	return "", fmt.Errorf("debug symbols not found on any server for build-id %s", buildID)
}

func (c *Client) fetchFromServerWithRetry(ctx context.Context, serverURL, buildID, destPath string) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		path, err := c.fetchFromServer(ctx, serverURL, buildID, destPath)
		if err == nil {
			return path, nil
		}

		lastErr = err

		var nonRetryable *nonRetryableError
		if errors.As(err, &nonRetryable) {
			return "", err
		}

		if attempt < c.maxRetries {
			backoff := c.calculateBackoff(attempt)
			c.logger.Debug("retrying after backoff", slog.Int("attempt", attempt), slog.Duration("backoff", backoff), slog.Any("error", err))

			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(backoff):
			}
		}
	}

	return "", lastErr
}

func (c *Client) calculateBackoff(attempt int) time.Duration {
	// Base delay: 1s, 2s, 4s... (exponential)
	baseDelay := time.Duration(1<<uint(attempt-1)) * time.Second // #nosec G115 -- bounded by maxRetries

	// Randomize between 100% and 200% of base delay
	return baseDelay + time.Duration(rand.Int64N(int64(baseDelay))) // #nosec G404 -- jitter, not security
}

func (c *Client) fetchFromServer(ctx context.Context, serverURL, buildID, destPath string) (string, error) {
	url := fmt.Sprintf("%s/buildid/%s/debuginfo", strings.TrimSuffix(serverURL, "/"), buildID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent())

	resp, err := c.httpClient.Do(req) // #nosec G704 -- URL from user-configured --debuginfod-servers
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("server returned %d", resp.StatusCode)
		if isNonRetryable(resp.StatusCode) {
			return "", &nonRetryableError{err: err}
		}
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0750); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	if err := downloadToFile(resp.Body, destPath, c.maxFileSize); err != nil {
		return "", err
	}

	return destPath, nil
}

func downloadToFile(r io.Reader, destPath string, maxSize int64) error {
	tmpPath := destPath + ".tmp"

	tmpFile, err := os.Create(tmpPath) // #nosec G304 -- path derived from validated cache directory
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	_, copyErr := io.Copy(tmpFile, io.LimitReader(r, maxSize))
	closeErr := tmpFile.Close()

	if copyErr != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("download failed: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", closeErr)
	}

	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to move file: %w", err)
	}

	return nil
}

func (c *Client) getCachePath(buildID string) string {
	if len(buildID) < 2 {
		return filepath.Join(c.cacheDir, buildID+".debug")
	}
	subdir := buildID[:2]
	return filepath.Join(c.cacheDir, subdir, buildID[2:]+".debug")
}
