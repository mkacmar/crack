package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseURLList(s string) []string {
	var urls []string
	for _, url := range strings.Split(s, ",") {
		if url = strings.TrimSpace(url); url != "" {
			urls = append(urls, url)
		}
	}
	return urls
}

func readPathsFromInput(input string) ([]string, error) {
	var reader *bufio.Scanner
	if input == "-" {
		reader = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(input)
		if err != nil {
			return nil, fmt.Errorf("failed to open input file: %w", err)
		}
		defer f.Close()
		reader = bufio.NewScanner(f)
	}

	var paths []string
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if line != "" {
			paths = append(paths, line)
		}
	}
	if err := reader.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}
	return paths, nil
}
