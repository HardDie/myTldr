package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func printList(source, platform, language string) (commands []string) {
	path := buildLocalPath(source, platform, language)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		commands = append(commands, strings.ReplaceAll(info.Name(), ".md", ""))
		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	sort.Strings(commands)
	return
}
