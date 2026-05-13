package analyzer

import (
	"context"
	"io"
	"log/slog"

	"go.kacmar.sk/crack/binary"
	"go.kacmar.sk/crack/binary/elf"
	"go.kacmar.sk/crack/internal/debuginfo"
	"go.kacmar.sk/crack/rule"
	"go.kacmar.sk/crack/toolchain"
	"go.kacmar.sk/debuginfod"
)

// ELFAnalyzer runs ELF-specific analysis and returns findings.
type ELFAnalyzer struct {
	rules            []rule.ELFRule
	debuginfodClient *debuginfod.Client
	detector         toolchain.ELFDetector
	logger           *slog.Logger
}

// ELFAnalyzerOptions configures ELFAnalyzer creation.
type ELFAnalyzerOptions struct {
	Rules            []rule.ELFRule
	DebuginfodClient *debuginfod.Client
	Detector         toolchain.ELFDetector
	Logger           *slog.Logger
}

// NewELFAnalyzer creates an ELF analyzer with the given options.
func NewELFAnalyzer(opts ELFAnalyzerOptions) *ELFAnalyzer {
	detector := opts.Detector
	if detector == nil {
		detector = toolchain.ELFCommentDetector{}
	}
	return &ELFAnalyzer{
		rules:            opts.Rules,
		debuginfodClient: opts.DebuginfodClient,
		detector:         detector,
		logger:           opts.Logger.With(slog.String("component", "elf-analyzer")),
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
		Toolchain:    elf.DetectToolchain(bin, a.detector),
	}
	if profile.Toolchain.Compiler == toolchain.Unknown {
		if tc := elf.DetectToolchainFromDWARF(bin, a.detector); tc.Compiler != toolchain.Unknown {
			profile.Toolchain = tc
			a.logger.Debug("detected toolchain from DWARF", slog.String("compiler", tc.Compiler.String()), slog.String("version", tc.Version.String()))
		}
	}

	findings := rule.Check(a.rules, profile, func(r rule.ELFRule) rule.Result {
		return r.Execute(bin)
	})
	return profile, bin.BuildID(), findings, nil
}

// resolverFactory builds a Resolver for a given build ID, scoped to ctx and the analyzer's debuginfod client.
// Returns nil for binaries without a build ID or when no client is configured. The binary then operates without remote section fetching.
func (a *ELFAnalyzer) resolverFactory(ctx context.Context) func(buildID string) elf.Resolver {
	return func(buildID string) elf.Resolver {
		if buildID == "" || a.debuginfodClient == nil {
			return nil
		}
		a.logger.Debug("resolver attached", slog.String("build_id", buildID))
		return debuginfo.NewResolver(ctx, buildID, a.debuginfodClient, a.logger)
	}
}
