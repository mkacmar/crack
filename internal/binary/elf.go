package binary

import "debug/elf"

type ELFBinary struct {
	Binary
	File       *elf.File
	Symbols    []elf.Symbol
	DynSymbols []elf.Symbol
}
