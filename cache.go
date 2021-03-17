package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	DBDefaultPath = ".cache/tldr"
	DBDefaultName = "pages.db"
	PagesSource   = "https://tldr.sh/assets/tldr.zip"
)

var (
	ErrorCacheNotExists = errors.New("cache doesn't exists")
	ErrorDataNotExists = errors.New("data not exists")
)

func getDBPath(homeDir string) (path string) {
	return homeDir + "/" + DBDefaultPath
}

func buildBucketName(cfg *Config) string {
	if *cfg.Language == "en" {
		return *cfg.Platform
	} else {
		return *cfg.Platform + "." + *cfg.Language
	}
}

func openBoldDB(source string) (db *bolt.DB, clean func(), err error) {
	// Create folder if not exists
	if !isFileExists(source) {
		err = os.MkdirAll(source, 0777)
		if err != nil {
			return
		}
	}

	db, err = bolt.Open(source+"/"+DBDefaultName, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return
	}

	clean = func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func checkCache(cfg *Config, name string) (page []string, err error) {
	// If DB not exist, return
	if !isFileExists(*cfg.Source) {
		err = ErrorCacheNotExists
		return
	}

	// Open DB
	db, clean, err := openBoldDB(*cfg.DBSource)
	if err != nil {
		return
	}
	defer clean()

	var data []byte
	// Build bucket name from input values
	bucketName := buildBucketName(cfg)
	err = db.View(func(tx *bolt.Tx) error {
		// Check if bucket exists
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return ErrorDataNotExists
		}
		// Check if data exists
		data = b.Get([]byte(name))
		if data == nil {
			return ErrorDataNotExists
		}
		return nil
	})
	if err != nil {
		return
	}

	// Split data to lines
	page = strings.Split(string(data), "\n")
	return
}

func putCache(cfg *Config, name string, data []byte) (err error) {
	db, clean, err := openBoldDB(*cfg.DBSource)
	if err != nil {
		return
	}
	defer clean()

	bucketName := buildBucketName(cfg)
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(name), data)
	})
	if err != nil {
		return
	}
	return
}

func updateCache(cfg *Config) (err error) {
	// Delete DB if exists
	if isFileExists(*cfg.DBSource + "/" + DBDefaultName) {
		if err = os.Remove(*cfg.DBSource + "/" + DBDefaultName); err != nil {
			return
		}
	}

	// Download archive with all pages
	zipReader, err := downloadZip(PagesSource)
	if err != nil {
		return
	}

	for _, file := range zipReader.File {
		// If started from "pages" and end with ".md", its page with info file
		if strings.HasSuffix(file.Name, ".md") && strings.HasPrefix(file.Name, "pages") {
			// Split by "/" symbol, to get every name separated
			path := strings.Split(strings.TrimSuffix(file.Name, ".md"), "/")

			// get language
			language := "en"
			folders := strings.Split(path[0], ".")
			if len(folders) == 2 {
				language = folders[1]
			}

			// get platform
			platform := path[1]

			// get command
			command := path[2]

			f, err := file.Open()
			if err != nil {
				return err
			}

			var data []byte
			if data, err = ioutil.ReadAll(f); err != nil {
				return err
			}

			tmpCfg := &Config{
				Source:   cfg.Source,
				DBSource: cfg.DBSource,
				Platform: &platform,
				Language: &language,
			}
			if err = putCache(tmpCfg, command, data); err != nil {
				return err
			}
		}
	}
	return
}

func printAllCache(source string) (err error) {
	db, clean, err := openBoldDB(source)
	if err != nil {
		return
	}
	defer clean()

	err = db.View(func(tx *bolt.Tx) error {
		_ = tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println("bucket:", string(name))
			return b.ForEach(func(k, v []byte) error {
				fmt.Println("   key:", string(k), "value:", strings.Split(string(v), "\n")[0])
				return nil
			})
		})
		return nil
	})
	if err != nil {
		return
	}

	return
}
