package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func printLocalList(cfg *Config) (commands []string) {
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

func printGlobalList(cfg *Config) (commands []string, err error) {
	if !isFileExists(*cfg.DBSource) {
		return
	}

	// Open DB
	bw, err := NewBoltWrapper(*cfg.DBSource + "/" + DBDefaultName)
	if err != nil {
		return
	}
	defer func() { _ = bw.Close() }()

	// Build bucket name from input values
	bucketName := buildBucketName(cfg)

	if commands, err = bw.GetKeysFromBucket(bucketName); err != nil {
		return
	}
	return
}
