package debuginfo

import (
	"context"
	stdelf "debug/elf"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"go.kacmar.sk/crack/binary/elf"
)

// DefaultBuildIDDir is the conventional root for separate debug files on Linux distributions.
// Matches GDB's default debug-file-directory.
const DefaultBuildIDDir = "/usr/lib/debug"

// BuildIDDirSource resolves ELF sections from local separate-debug files laid out under
// <root>/.build-id/<xx>/<rest>.debug, the convention shared by GDB, debuginfod, and distro packagers.
type BuildIDDirSource struct {
	root   string
	logger *slog.Logger
}

// NewBuildIDDirSource constructs a Source rooted at the given directory.
// An empty root falls back to DefaultBuildIDDir.
func NewBuildIDDirSource(root string, logger *slog.Logger) *BuildIDDirSource {
	if root == "" {
		root = DefaultBuildIDDir
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &BuildIDDirSource{root: root, logger: logger}
}

// ResolverFor returns a Resolver scoped to the given build ID.
// The context is accepted for interface conformance but local I/O does not consult it.
func (s *BuildIDDirSource) ResolverFor(_ context.Context, buildID string) elf.Resolver {
	if buildID == "" {
		panic("debuginfo.BuildIDDirSource.ResolverFor: empty buildID")
	}
	return &buildIDDirResolver{root: s.root, buildID: buildID, logger: s.logger}
}

type buildIDDirResolver struct {
	root    string
	buildID string
	logger  *slog.Logger
}

// FetchSection opens the build-id-indexed .debug file and returns the named section's bytes.
// Returns ErrSectionMissing when the .debug file or the requested section is absent.
func (r *buildIDDirResolver) FetchSection(name string) ([]byte, error) {
	path := r.debugFilePath()
	r.logger.Debug("build-id dir source fetching section",
		slog.String("build_id", r.buildID),
		slog.String("section", name),
		slog.String("path", path))

	f, err := stdelf.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, os.ErrNotExist) {
			return nil, elf.ErrSectionMissing
		}
		return nil, fmt.Errorf("build-id dir open %s: %w", path, err)
	}
	defer f.Close()

	section := f.Section(name)
	if section == nil {
		return nil, elf.ErrSectionMissing
	}
	data, err := section.Data()
	if err != nil {
		return nil, fmt.Errorf("build-id dir read %s from %s: %w", name, path, err)
	}
	return data, nil
}

// debugFilePath returns the on-disk debug file location under the build-ID directory layout.
// The layout is <root>/.build-id/<xx>/<rest>.debug, where <xx> is the first two hex characters of the build ID and <rest> is the remainder.
// See GDB's "Separate Debug Files" manual section: https://sourceware.org/gdb/current/onlinedocs/gdb.html/Separate-Debug-Files.html
func (r *buildIDDirResolver) debugFilePath() string {
	return filepath.Join(r.root, ".build-id", r.buildID[:2], r.buildID[2:]+".debug")
}
