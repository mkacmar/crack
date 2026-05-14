package debuginfo

import (
	"errors"

	"go.kacmar.sk/crack/binary/elf"
)

// Chain composes multiple resolvers and queries them in order.
// A resolver that returns ErrSectionMissing causes the chain to try the next entry.
// Any other error short-circuits and is returned to the caller.
type Chain []elf.Resolver

// FetchSection walks the chain and returns the first successful section fetch.
// Returns ErrSectionMissing when every resolver reports the section as missing.
func (c Chain) FetchSection(name string) ([]byte, error) {
	for _, r := range c {
		if r == nil {
			continue
		}
		data, err := r.FetchSection(name)
		if err == nil {
			return data, nil
		}
		if errors.Is(err, elf.ErrSectionMissing) {
			continue
		}
		return nil, err
	}
	return nil, elf.ErrSectionMissing
}
