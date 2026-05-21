package debuginfo

import (
	"context"
	stdelf "debug/elf"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/internal/version"
	"go.kacmar.sk/debuginfod"
	"go.kacmar.sk/debuginfod/cache"
	"go.kacmar.sk/debuginfod/key"
)

const DefaultMaxRetries = 2

func DefaultCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("debuginfod: failed to determine user cache directory: %w", err)
	}
	return filepath.Join(cacheDir, "debuginfod"), nil
}

func userAgent() string {
	return "crack/" + version.Version + " (Compiler Hardening Checker; +https://github.com/mkacmar/crack)"
}

type Options struct {
	ServerURLs []string
	CacheDir   string
	Timeout    time.Duration
	MaxRetries int
	Logger     *slog.Logger
}

// NewCache constructs a disk-backed debuginfod cache configured from opts.
func NewCache(opts Options) (*cache.DiskCache, error) {
	cacheDir := opts.CacheDir
	if cacheDir == "" {
		var err error
		cacheDir, err = DefaultCacheDir()
		if err != nil {
			return nil, err
		}
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	maxRetries := opts.MaxRetries
	if maxRetries < 1 {
		maxRetries = DefaultMaxRetries
	}

	client, err := debuginfod.NewClient(debuginfod.Options{
		ServerURLs: opts.ServerURLs,
		HTTP: debuginfod.HTTPOptions{
			Client:     &http.Client{Timeout: timeout},
			MaxRetries: maxRetries,
			UserAgent:  userAgent(),
		},
		Logger: opts.Logger,
	})
	if err != nil {
		return nil, err
	}

	disk, err := cache.NewDiskCache(cache.DiskCacheOptions{
		Client: client,
		Dir:    cacheDir,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}
	return disk, nil
}

// artifactCache fetches cached debuginfod artifacts as random-access files.
type artifactCache interface {
	Get(ctx context.Context, k key.Key) (*os.File, error)
}

// DebuginfodSource resolves ELF sections via a debuginfod cache.
type DebuginfodSource struct {
	cache  artifactCache
	logger *slog.Logger
}

// NewDebuginfodSource constructs a Source backed by the given debuginfod cache.
func NewDebuginfodSource(c artifactCache, logger *slog.Logger) *DebuginfodSource {
	if c == nil {
		panic("debuginfo.NewDebuginfodSource: nil cache")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &DebuginfodSource{cache: c, logger: logger}
}

// ResolverFor returns a Resolver scoped to the given context and build ID.
func (s *DebuginfodSource) ResolverFor(ctx context.Context, buildID string) elf.Resolver {
	if buildID == "" {
		panic("debuginfo.DebuginfodSource.ResolverFor: empty buildID")
	}
	return &debuginfodResolver{ctx: ctx, buildID: buildID, cache: s.cache, logger: s.logger}
}

type debuginfodResolver struct {
	ctx     context.Context
	buildID string
	cache   artifactCache
	logger  *slog.Logger
}

// FetchSection retrieves the named section's raw bytes via debuginfod.
// Returns ErrSectionMissing when the server reports the artifact as not found.
// On a not-found response we fall back to fetching the full debuginfo file and extracting the section locally,
// since older servers may not expose the /section/ endpoint at all.
func (r *debuginfodResolver) FetchSection(name string) ([]byte, error) {
	r.logger.Debug("debuginfod source fetching section", slog.String("build_id", r.buildID), slog.String("section", name))

	f, err := r.cache.Get(r.ctx, key.Section(r.buildID, name))
	if err == nil {
		defer f.Close()
		data, readErr := io.ReadAll(f)
		if readErr != nil {
			return nil, fmt.Errorf("debuginfod read %s: %w", name, readErr)
		}
		return data, nil
	}
	if !errors.Is(err, debuginfod.ErrNotFound) {
		return nil, fmt.Errorf("debuginfod fetch %s: %w", name, err)
	}

	return r.fetchSectionViaDebugInfo(name)
}

func (r *debuginfodResolver) fetchSectionViaDebugInfo(name string) ([]byte, error) {
	r.logger.Debug("debuginfod section endpoint returned not-found, falling back to full debuginfo",
		slog.String("build_id", r.buildID), slog.String("section", name))

	f, err := r.cache.Get(r.ctx, key.DebugInfo(r.buildID))
	if err != nil {
		if errors.Is(err, debuginfod.ErrNotFound) {
			return nil, elf.ErrSectionMissing
		}
		return nil, fmt.Errorf("debuginfod fetch debuginfo: %w", err)
	}
	defer f.Close()

	return extractSection(f, name)
}

func extractSection(f *os.File, name string) ([]byte, error) {
	ef, err := stdelf.NewFile(f)
	if err != nil {
		return nil, fmt.Errorf("debuginfod parse debuginfo: %w", err)
	}
	defer ef.Close()

	sect := ef.Section(name)
	if sect == nil {
		return nil, elf.ErrSectionMissing
	}
	data, err := sect.Data()
	if err != nil {
		return nil, fmt.Errorf("debuginfod read section %s from debuginfo: %w", name, err)
	}
	return data, nil
}
