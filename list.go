package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/boltdb/bolt"
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

func printGlobalList(cfg *Config) (commands []string) {
	if !isFileExists(*cfg.DBSource) {
		return
	}

	// Open DB
	db, clean, err := openBoldDB(*cfg.DBSource)
	if err != nil {
		return
	}
	defer clean()

	// Build bucket name from input values
	bucketName := buildBucketName(cfg)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			commands = append(commands, string(k))
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}
	return
}
