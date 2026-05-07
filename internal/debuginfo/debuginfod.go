package debuginfo

import (
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"go.kacmar.sk/crack/internal/version"
	"go.kacmar.sk/debuginfod"
)

func userAgent() string {
	return "crack/" + version.Version + " (Compiler Hardening Checker; +https://github.com/mkacmar/crack)"
}

// DefaultMaxRetries is the default number of additional retry rounds after the initial attempt.
const DefaultMaxRetries = 2 // 3 total attempts (initial + 2 retries)

type Options struct {
	ServerURLs []string
	CacheDir   string
	Timeout    time.Duration
	MaxRetries int
	Logger     *slog.Logger
}

func NewClient(opts Options) (*debuginfod.Client, error) {
	cacheDir := opts.CacheDir
	if cacheDir == "" {
		var err error
		cacheDir, err = debuginfod.DefaultCacheDir()
		if err != nil {
			return nil, err
		}
	}

	cache, err := debuginfod.NewDiskCache(debuginfod.DiskCacheOptions{
		Dir: cacheDir,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	maxRetries := opts.MaxRetries
	if maxRetries < 1 {
		maxRetries = DefaultMaxRetries
	}

	return debuginfod.NewClient(debuginfod.Options{
		ServerURLs: opts.ServerURLs,
		Cache:      cache,
		HTTP: debuginfod.HTTPOptions{
			Client:     &http.Client{Timeout: timeout},
			MaxRetries: maxRetries,
			UserAgent:  userAgent(),
		},
		Logger: opts.Logger,
	})
}
