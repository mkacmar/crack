package debuginfo

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/debuginfod"
	"go.kacmar.sk/debuginfod/cache"
	"go.kacmar.sk/debuginfod/key"
)

// fakeCache writes artifacts to files, so entries expose the real *os.File that cache.Entry embeds.
type fakeCache struct {
	dir  string
	data map[key.Key][]byte
	errs map[key.Key]error
	gets []key.Key
}

func newFakeCache(t *testing.T) *fakeCache {
	t.Helper()
	return &fakeCache{
		dir:  t.TempDir(),
		data: make(map[key.Key][]byte),
		errs: make(map[key.Key]error),
	}
}

func (f *fakeCache) Get(_ context.Context, k key.Key) (cache.Entry, error) {
	f.gets = append(f.gets, k)
	if err, ok := f.errs[k]; ok {
		return cache.Entry{}, err
	}
	data, ok := f.data[k]
	if !ok {
		return cache.Entry{}, debuginfod.ErrNotFound
	}
	file, err := os.CreateTemp(f.dir, "artifact-*")
	if err != nil {
		return cache.Entry{}, err
	}
	if _, err := file.Write(data); err != nil {
		return cache.Entry{}, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return cache.Entry{}, err
	}
	return cache.Entry{File: file}, nil
}

func newTestSource(t *testing.T, c *fakeCache) elf.Resolver {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewDebuginfodSource(c, logger).ResolverFor(context.Background(), testBuildID)
}

func TestDebuginfodFetchSection(t *testing.T) {
	c := newFakeCache(t)
	c.data[key.Section(testBuildID, ".text")] = []byte("payload")

	got, err := newTestSource(t, c).FetchSection(".text")
	if err != nil {
		t.Fatalf("FetchSection: %v", err)
	}
	if string(got) != "payload" {
		t.Errorf("FetchSection = %q, want %q", got, "payload")
	}
	if len(c.gets) != 1 {
		t.Errorf("cache lookups = %d, want 1 (no fallback expected)", len(c.gets))
	}
}

func TestDebuginfodFetchSectionError(t *testing.T) {
	sentinel := errors.New("transport failed")
	c := newFakeCache(t)
	c.errs[key.Section(testBuildID, ".text")] = sentinel

	_, err := newTestSource(t, c).FetchSection(".text")
	if !errors.Is(err, sentinel) {
		t.Errorf("error = %v, want it to wrap %v", err, sentinel)
	}
	if len(c.gets) != 1 {
		t.Errorf("cache lookups = %d, want 1 (a real error must not fall back)", len(c.gets))
	}
}

func TestDebuginfodFallsBackToDebugInfo(t *testing.T) {
	c := newFakeCache(t)
	c.data[key.DebugInfo(testBuildID)] = makeMinimalELF(".debug_info", []byte("dwarf"))

	got, err := newTestSource(t, c).FetchSection(".debug_info")
	if err != nil {
		t.Fatalf("FetchSection: %v", err)
	}
	if string(got) != "dwarf" {
		t.Errorf("FetchSection = %q, want %q", got, "dwarf")
	}
	if len(c.gets) != 2 {
		t.Errorf("cache lookups = %d, want 2 (section then debuginfo)", len(c.gets))
	}
}

func TestDebuginfodSectionAbsentFromDebugInfo(t *testing.T) {
	c := newFakeCache(t)
	c.data[key.DebugInfo(testBuildID)] = makeMinimalELF(".debug_info", []byte("dwarf"))

	_, err := newTestSource(t, c).FetchSection(".debug_line")
	if !errors.Is(err, elf.ErrSectionMissing) {
		t.Errorf("error = %v, want %v", err, elf.ErrSectionMissing)
	}
}

func TestDebuginfodBothNotFound(t *testing.T) {
	c := newFakeCache(t)

	_, err := newTestSource(t, c).FetchSection(".text")
	if !errors.Is(err, elf.ErrSectionMissing) {
		t.Errorf("error = %v, want %v", err, elf.ErrSectionMissing)
	}
}

func TestDebuginfodDebugInfoError(t *testing.T) {
	sentinel := errors.New("transport failed")
	c := newFakeCache(t)
	c.errs[key.DebugInfo(testBuildID)] = sentinel

	_, err := newTestSource(t, c).FetchSection(".text")
	if !errors.Is(err, sentinel) {
		t.Errorf("error = %v, want it to wrap %v", err, sentinel)
	}
}

func TestDebuginfodDebugInfoNotELF(t *testing.T) {
	c := newFakeCache(t)
	c.data[key.DebugInfo(testBuildID)] = []byte("not an ELF file at all")

	_, err := newTestSource(t, c).FetchSection(".text")
	if err == nil {
		t.Fatal("FetchSection: expected an error for a non-ELF debuginfo artifact")
	}
	if errors.Is(err, elf.ErrSectionMissing) {
		t.Errorf("error = %v, want a parse failure rather than ErrSectionMissing", err)
	}
}

func TestNewDebuginfodSourceNilCachePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewDebuginfodSource(nil, ...): expected a panic")
		}
	}()
	NewDebuginfodSource(nil, nil)
}

func TestDebuginfodResolverForEmptyBuildIDPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("ResolverFor(\"\"): expected a panic")
		}
	}()
	NewDebuginfodSource(newFakeCache(t), nil).ResolverFor(context.Background(), "")
}
