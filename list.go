package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func printList(source, platform, language string) (commands []string) {
	path := buildLocalPath(source, platform, language)
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
