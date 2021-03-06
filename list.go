package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func printList(cfg *Config) (commands []string) {
	path := buildLocalPath(cfg)
	if !isFileExists(path) {
		return
	}
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		commands = append(commands, strings.ReplaceAll(info.Name(), ".md", ""))
		return nil
	})
	sort.Strings(commands)
	return
}
