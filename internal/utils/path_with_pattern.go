package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func PathToPathWithPattern(input string) string {
	// let's check if this is already a file by any chance
	if _, err := os.Stat(input); err == nil {
		return input
	}
	ext := filepath.Ext(input)
	base := strings.Replace(filepath.Base(input), ext, "", 1)
	dir := filepath.Dir(input)

	fmt.Printf("input: %s, ext: %s; base: %s; dir: %s\n", input, ext, base, dir)

	if strings.HasPrefix(base, "*") {
		return filepath.Join(dir, base+ext)
	}

	if ext == "" {
		return filepath.Join(input, "*")
	}

	return input
}
