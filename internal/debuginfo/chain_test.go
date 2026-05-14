package debuginfo

import (
	"errors"
	"testing"

	"go.kacmar.sk/crack/binary/elf"
)

type stubResolver struct {
	data []byte
	err  error
	hits int
}

func (s *stubResolver) FetchSection(name string) ([]byte, error) {
	_ = name
	s.hits++
	return s.data, s.err
}

func TestChainEmpty(t *testing.T) {
	var c Chain
	_, err := c.FetchSection(".symtab")
	if !errors.Is(err, elf.ErrSectionMissing) {
		t.Fatalf("empty chain should return ErrSectionMissing, got %v", err)
	}
}

func TestChainSingleSuccess(t *testing.T) {
	r := &stubResolver{data: []byte("hello")}
	c := Chain{r}
	data, err := c.FetchSection(".symtab")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q, want %q", data, "hello")
	}
}

func TestChainFallthroughOnMissing(t *testing.T) {
	first := &stubResolver{err: elf.ErrSectionMissing}
	second := &stubResolver{data: []byte("found")}
	c := Chain{first, second}

	data, err := c.FetchSection(".symtab")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "found" {
		t.Fatalf("got %q, want %q", data, "found")
	}
	if first.hits != 1 || second.hits != 1 {
		t.Fatalf("hit counts: first=%d second=%d, want 1/1", first.hits, second.hits)
	}
}

func TestChainShortCircuitsOnRealError(t *testing.T) {
	boom := errors.New("network down")
	first := &stubResolver{err: boom}
	second := &stubResolver{data: []byte("never reached")}
	c := Chain{first, second}

	_, err := c.FetchSection(".symtab")
	if !errors.Is(err, boom) {
		t.Fatalf("got %v, want wrapped boom", err)
	}
	if second.hits != 0 {
		t.Fatalf("second resolver should not be consulted, hits=%d", second.hits)
	}
}

func TestChainAllMissingReturnsMissing(t *testing.T) {
	first := &stubResolver{err: elf.ErrSectionMissing}
	second := &stubResolver{err: elf.ErrSectionMissing}
	c := Chain{first, second}

	_, err := c.FetchSection(".symtab")
	if !errors.Is(err, elf.ErrSectionMissing) {
		t.Fatalf("got %v, want ErrSectionMissing", err)
	}
}

func TestChainSkipsNilEntries(t *testing.T) {
	r := &stubResolver{data: []byte("ok")}
	c := Chain{nil, r, nil}

	data, err := c.FetchSection(".symtab")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "ok" {
		t.Fatalf("got %q, want %q", data, "ok")
	}
}
