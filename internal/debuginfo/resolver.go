package debuginfo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	binelf "go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/debuginfod"
)

// Resolver fetches ELF sections via a debuginfod client, scoped to one build ID and one request context.
type Resolver struct {
	ctx     context.Context
	buildID string
	client  *debuginfod.Client
	logger  *slog.Logger
}

// NewResolver constructs a Resolver scoped to the given context and build ID.
// buildID must be non-empty and client must be non-nil. Callers are expected to filter out the empty or no-client case before invoking.
func NewResolver(ctx context.Context, buildID string, client *debuginfod.Client, logger *slog.Logger) *Resolver {
	if buildID == "" {
		panic("debuginfo.NewResolver: empty buildID")
	}
	if client == nil {
		panic("debuginfo.NewResolver: nil client")
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Resolver{ctx: ctx, buildID: buildID, client: client, logger: logger}
}

// FetchSection retrieves the named section's raw bytes via debuginfod.
// Returns binelf.ErrSectionMissing when the server reports the artifact as not found.
func (r *Resolver) FetchSection(name string) ([]byte, error) {
	r.logger.Debug("resolver fetching section", slog.String("build_id", r.buildID), slog.String("section", name))

	rc, err := r.client.FetchSection(r.ctx, r.buildID, name)
	if err != nil {
		if errors.Is(err, debuginfod.ErrNotFound) {
			return nil, binelf.ErrSectionMissing
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
