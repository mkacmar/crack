package debuginfo

import (
	"context"

	"go.kacmar.sk/crack/binary/elf"
)

// Source produces a Resolver bound to a single binary's build ID.
// Implementations encapsulate a kind of debug-information backend.
// The analyzer composes one or more Sources and chains their Resolvers when fetching missing sections.
type Source interface {
	ResolverFor(ctx context.Context, buildID string) elf.Resolver
}
