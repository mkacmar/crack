package debuginfo

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mkacmar/crack/binary"
	"github.com/mkacmar/crack/toolchain"
)

func EnhanceWithDebugInfo(bin *binary.ELFBinary, debugPath string, logger *slog.Logger) error {
	logger = logger.With(slog.String("component", "dwarf"))

	debugFile, err := elf.Open(debugPath)
	if err != nil {
		return fmt.Errorf("failed to open debug file: %w", err)
	}
	defer debugFile.Close()

	if debugSymbols, err := debugFile.Symbols(); err == nil && len(debugSymbols) > 0 {
		bin.Symbols = mergeSymbols(bin.Symbols, debugSymbols)
		logger.Debug("merged symbols", slog.Int("total", len(bin.Symbols)))
	}

	if !hasDwarfSections(debugFile) {
		return nil
	}

	dwarfData, err := debugFile.DWARF()
	if err != nil {
		return fmt.Errorf("failed to extract DWARF data: %w", err)
	}

	compilerInfo := extractCompilerFromDWARF(dwarfData, logger)
	if compilerInfo == "" {
		logger.Debug("no DW_AT_producer found in DWARF data", slog.String("debug_path", debugPath))
		return nil
	}

	detector := toolchain.ELFCommentDetector{}
	newCompiler, newVersion := detector.Detect(compilerInfo)
	if bin.Build.Compiler == toolchain.Unknown && newCompiler != toolchain.Unknown {
		bin.Build.Compiler = newCompiler
		bin.Build.Version = newVersion
		logger.Debug("updated toolchain from DWARF", slog.String("compiler", newCompiler.String()), slog.String("version", newVersion.String()))
	}

	return nil
}

func mergeSymbols(binarySymbols, debugSymbols []elf.Symbol) []elf.Symbol {
	if len(binarySymbols) == 0 {
		return debugSymbols
	}
	if len(debugSymbols) == 0 {
		return binarySymbols
	}

	existing := make(map[string]struct{}, len(binarySymbols))
	for _, sym := range binarySymbols {
		existing[sym.Name] = struct{}{}
	}

	merged := make([]elf.Symbol, len(binarySymbols), len(binarySymbols)+len(debugSymbols))
	copy(merged, binarySymbols)

	for _, sym := range debugSymbols {
		if _, exists := existing[sym.Name]; !exists {
			merged = append(merged, sym)
			existing[sym.Name] = struct{}{}
		}
	}

	return merged
}

func hasDwarfSections(f *elf.File) bool {
	for _, section := range f.Sections {
		if strings.HasPrefix(section.Name, ".debug_") {
			return true
		}
	}
	return false
}

func extractCompilerFromDWARF(d *dwarf.Data, logger *slog.Logger) string {
	reader := d.Reader()

	for {
		entry, err := reader.Next()
		if err != nil || entry == nil {
			break
		}

		if entry.Tag == dwarf.TagCompileUnit {
			if producer, ok := entry.Val(dwarf.AttrProducer).(string); ok && producer != "" {
				logger.Debug("found DW_AT_producer in DWARF", slog.String("producer", producer))
				return producer
			}
		}
	}

	return ""
}
