package debuginfo

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"log/slog"

	"github.com/mkacmar/crack/internal/version"
)

const (
	DefaultServerURL = "https://debuginfod.elfutils.org"
	DefaultTimeout   = 30 * time.Second
	DefaultRetries   = 3
)

func userAgent() string {
	return "crack/" + version.Version + " (Compiler Hardening Checker; +https://github.com/mkacmar/crack)"
}

type Client struct {
	serverURLs []string
	cacheDir   string
	httpClient *http.Client
	maxRetries int
	logger     *slog.Logger
}

type Options struct {
	ServerURLs []string
	CacheDir   string
	Timeout    time.Duration
	MaxRetries int
	Logger     *slog.Logger
}

func NewClient(opts Options) (*Client, error) {
	if len(opts.ServerURLs) == 0 {
		return nil, fmt.Errorf("no debuginfod servers configured")
	}

	if opts.CacheDir == "" {
		return nil, fmt.Errorf("cache directory not configured")
	}

	if err := os.MkdirAll(opts.CacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	maxRetries := opts.MaxRetries
	if maxRetries < 1 {
		maxRetries = 3
	}

	client := &Client{
		serverURLs: opts.ServerURLs,
		cacheDir:   opts.CacheDir,
		httpClient: &http.Client{Timeout: timeout},
		logger:     opts.Logger.With(slog.String("component", "debuginfod")),
		maxRetries: maxRetries,
	}

	opts.Logger.Debug("debuginfod client initialized",
		slog.Any("servers", opts.ServerURLs),
		slog.String("cache", opts.CacheDir))

	return client, nil
}

func (c *Client) FetchDebugInfo(ctx context.Context, buildID string) (string, error) {
	if buildID == "" {
		return "", fmt.Errorf("build-id is empty")
	}

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

	return "", fmt.Errorf("debug symbols not found on any server")
}

func (c *Client) fetchFromServerWithRetry(ctx context.Context, serverURL, buildID, destPath string) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		path, err := c.fetchFromServer(ctx, serverURL, buildID, destPath)
		if err == nil {
			return path, nil
		}

		lastErr = err

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
	baseDelay := time.Duration(1<<uint(attempt-1)) * time.Second

	// Add jitter: +-25% of base delay
	jitter := time.Duration(rand.Int64N(int64(baseDelay / 2)))
	if rand.IntN(2) == 0 {
		return baseDelay + jitter
	}
	return baseDelay - jitter/2
}

func (c *Client) fetchFromServer(ctx context.Context, serverURL, buildID, destPath string) (string, error) {
	url := fmt.Sprintf("%s/buildid/%s/debuginfo", strings.TrimSuffix(serverURL, "/"), buildID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	tmpPath := destPath + ".tmp"
	cleanup := true
	defer func() {
		if cleanup {
			os.Remove(tmpPath)
		}
	}()

	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	written, err := io.Copy(tmpFile, resp.Body)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("download failed: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, destPath); err != nil {
		return "", fmt.Errorf("failed to move file: %w", err)
	}

	cleanup = false
	c.logger.Debug("downloaded debug symbols", slog.Int64("size", written), slog.String("path", destPath))

	return destPath, nil
}

func (c *Client) getCachePath(buildID string) string {
	if len(buildID) < 2 {
		return filepath.Join(c.cacheDir, buildID+".debug")
	}
	subdir := buildID[:2]
	return filepath.Join(c.cacheDir, subdir, buildID[2:]+".debug")
}
