package debuginfo

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mkacmar/crack/internal/model"
)

func EnhanceWithDebugInfo(info *model.ParsedBinary, debugPath string, logger *slog.Logger) error {
	logger = logger.With(slog.String("component", "dwarf"))

	debugFile, err := elf.Open(debugPath)
	if err != nil {
		return fmt.Errorf("failed to open debug file: %w", err)
	}
	defer debugFile.Close()

	if !hasDwarfSections(debugFile) {
		return fmt.Errorf("debug file contains no DWARF data")
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

	newToolchain := model.ParseToolchain(compilerInfo)
	if info.Build.Toolchain.Compiler == model.CompilerUnknown && newToolchain.Compiler != model.CompilerUnknown {
		info.Build.Toolchain = newToolchain
		logger.Debug("updated toolchain from DWARF", slog.String("compiler", newToolchain.Compiler.String()), slog.String("version", newToolchain.Version.String()))
	}

	return nil
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
