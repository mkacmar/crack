package debuginfo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/debuginfod"
)

// DebuginfodSource resolves ELF sections via a debuginfod client.
type DebuginfodSource struct {
	client *debuginfod.Client
	logger *slog.Logger
}

// NewDebuginfodSource constructs a Source backed by the given debuginfod client.
func NewDebuginfodSource(client *debuginfod.Client, logger *slog.Logger) *DebuginfodSource {
	if client == nil {
		panic("debuginfo.NewDebuginfodSource: nil client")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &DebuginfodSource{client: client, logger: logger}
}

// ResolverFor returns a Resolver scoped to the given context and build ID.
func (s *DebuginfodSource) ResolverFor(ctx context.Context, buildID string) elf.Resolver {
	if buildID == "" {
		panic("debuginfo.DebuginfodSource.ResolverFor: empty buildID")
	}
	return &debuginfodResolver{ctx: ctx, buildID: buildID, client: s.client, logger: s.logger}
}

type debuginfodResolver struct {
	ctx     context.Context
	buildID string
	client  *debuginfod.Client
	logger  *slog.Logger
}

// FetchSection retrieves the named section's raw bytes via debuginfod.
// Returns ErrSectionMissing when the server reports the artifact as not found.
func (r *debuginfodResolver) FetchSection(name string) ([]byte, error) {
	r.logger.Debug("debuginfod source fetching section", slog.String("build_id", r.buildID), slog.String("section", name))

	rc, err := r.client.FetchSection(r.ctx, r.buildID, name)
	if err != nil {
		if errors.Is(err, debuginfod.ErrNotFound) {
			return nil, elf.ErrSectionMissing
		}
		return nil, fmt.Errorf("debuginfod fetch %s: %w", name, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("debuginfod read %s: %w", name, err)
	}
	return data, nil
}
