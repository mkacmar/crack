//go:build ignore

package main

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <binary>\n", os.Args[0])
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stat file: %v\n", err)
		os.Exit(1)
	}

	needle := make([]byte, 8)
	binary.LittleEndian.PutUint64(needle, uint64(elf.DT_DEBUG))

	idx := bytes.Index(data, needle)
	if idx == -1 {
		fmt.Fprintln(os.Stderr, "DT_DEBUG not found")
		os.Exit(1)
	}

	binary.LittleEndian.PutUint64(data[idx:], uint64(elf.DT_TEXTREL))

	if err := os.WriteFile(os.Args[1], data, info.Mode()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
		os.Exit(1)
	}
}
