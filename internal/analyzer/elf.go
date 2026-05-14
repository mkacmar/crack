package analyzer

import (
	"context"
	"io"
	"log/slog"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/internal/debuginfo"
	"go.kacmar.sk/crack/rule"
)

// ELFAnalyzer runs ELF-specific analysis and returns findings.
type ELFAnalyzer struct {
	rules    []rule.ELFRule
	sources  []debuginfo.Source
	detector elf.ToolchainDetector
	logger   *slog.Logger
}

// ELFAnalyzerOptions configures ELFAnalyzer creation.
type ELFAnalyzerOptions struct {
	Rules    []rule.ELFRule
	Sources  []debuginfo.Source
	Detector elf.ToolchainDetector
	Logger   *slog.Logger
}

// NewELFAnalyzer creates an ELF analyzer with the given options.
func NewELFAnalyzer(opts ELFAnalyzerOptions) *ELFAnalyzer {
	detector := opts.Detector
	if detector == nil {
		detector = elf.DefaultToolchainDetector{}
	}
	return &ELFAnalyzer{
		rules:    opts.Rules,
		sources:  opts.Sources,
		detector: detector,
		logger:   opts.Logger.With(slog.String("component", "elf-analyzer")),
	}
}

// Analyze opens r as an ELF binary, runs ELF-specific rules, and returns the composed Profile alongside findings.
// Returns binary.ErrUnsupportedFormat when r isn't an ELF file.
func (a *ELFAnalyzer) Analyze(ctx context.Context, r io.ReaderAt) (binary.Profile, string, []rule.Finding, error) {
	bin, err := elf.Open(r, elf.WithResolverFactory(a.resolverFactory(ctx)))
	if err != nil {
		return binary.Profile{}, "", nil, err
	}

	profile := binary.Profile{
		Architecture: elf.DetectArchitecture(bin),
		LibC:         elf.DetectLibC(bin),
		Toolchain:    a.detector.Detect(bin),
	}

	findings := rule.Check(a.rules, profile, func(r rule.ELFRule) rule.Result {
		return r.Execute(bin)
	})
	return profile, bin.BuildID(), findings, nil
}

// resolverFactory builds a Resolver for a given build ID, scoped to ctx.
// Composes the configured sources in declared order, chaining them when more than one applies.
// Returns nil for binaries without a build ID or when no source is configured.
func (a *ELFAnalyzer) resolverFactory(ctx context.Context) func(buildID string) elf.Resolver {
	return func(buildID string) elf.Resolver {
		if buildID == "" || len(a.sources) == 0 {
			return nil
		}

		resolvers := make([]elf.Resolver, 0, len(a.sources))
		for _, src := range a.sources {
			resolvers = append(resolvers, src.ResolverFor(ctx, buildID))
		}

		if len(resolvers) == 1 {
			a.logger.Debug("resolver attached", slog.String("build_id", buildID))
			return resolvers[0]
		}
		a.logger.Debug("resolver chain attached", slog.String("build_id", buildID), slog.Int("sources", len(resolvers)))
		return debuginfo.Chain(resolvers)
	}
}
