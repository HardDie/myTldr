package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	DBDefaultPath = ".cache/tldr"
	DBDefaultName = "pages.db"
	PagesSource   = "https://tldr.sh/assets/tldr.zip"
)

var (
	ErrorCacheNotExists = errors.New("cache doesn't exists")
)

func getDBPath(homeDir string) (path string) {
	return homeDir + "/" + DBDefaultPath
}

func buildBucketName(cfg *Config, platform string) string {
	if *cfg.Language == "en" {
		return platform
	}
	return platform + "." + *cfg.Language
}

func checkCache(cfg *Config, name string) (page []string, err error) {
	// If DB not exist, return
	if !isFileExists(*cfg.DBSource + "/" + DBDefaultName) {
		err = ErrorCacheNotExists
		return
	}

	// // Open DB
	bw, err := NewBoltWrapper(*cfg.DBSource + "/" + DBDefaultName)
	if err != nil {
		return
	}
	defer func() { _ = bw.Close() }()

	platforms := []string{PlatformCommon, *cfg.Platform}

	// Check file for all platforms
	var data string
	for _, platform := range platforms {
		// Build bucket name from input values
		bucketName := buildBucketName(cfg, platform)
		data, err = bw.GetDataFromBucket(bucketName, name)
		if err == ErrorInvalidKey || err == ErrorCannotFindBucket {
			// If such data not exists just check next platform
			continue
		}
		if err != nil {
			return
		}
		break
	}

	if err != nil {
		return
	}
	// Split data to lines
	page = strings.Split(data, "\n")
	return
}

func openCache(cfg *Config) (bw *BoltWrapper, err error) {
	return NewBoltWrapper(*cfg.DBSource + "/" + DBDefaultName)
}

func putCache(cfg *Config, bw *BoltWrapper, name string, data []byte) (err error) {
	bucketName := buildBucketName(cfg, *cfg.Platform)
	if err = bw.CreateBucketIfNotExists(bucketName); err != nil {
		return
	}
	if err = bw.PutDataToBucket(bucketName, name, data); err != nil {
		return
	}
	return
}

func updateCache(cfg *Config) (err error) {
	// Check if cache folder exists
	if !isFileExists(*cfg.DBSource) {
		if err = os.MkdirAll(*cfg.DBSource, os.ModePerm); err != nil {
			return
		}
		fmt.Printf("Cache folder %s were created!\n", *cfg.DBSource)
	}

	// Delete DB if exists
	if isFileExists(*cfg.DBSource + "/" + DBDefaultName) {
		if err = os.Remove(*cfg.DBSource + "/" + DBDefaultName); err != nil {
			return
		}
		fmt.Printf("Cache file %s were removed!\n", *cfg.DBSource+"/"+DBDefaultName)
	}

	// Download archive with all pages
	zipReader, err := downloadZip(PagesSource)
	if err != nil {
		return
	}

	bw, err := openCache(cfg)
	if err != nil {
		return
	}
	defer func() { _ = bw.Close() }()
	fmt.Printf("Cache file %s were created!\n", *cfg.DBSource+"/"+DBDefaultName)

	fmt.Println("Start caching...")

	num := 0
	for _, file := range zipReader.File {
		num++
		printProgress(num, len(zipReader.File))

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
				_ = f.Close()
				return err
			}
			_ = f.Close()

			tmpCfg := &Config{
				Source:   cfg.Source,
				DBSource: cfg.DBSource,
				Platform: &platform,
				Language: &language,
			}
			if err = putCache(tmpCfg, bw, command, data); err != nil {
				return err
			}
		}
	}
	return nil
}

func getCacheInfo(cfg *Config) string {
	cachePath := *cfg.DBSource + "/" + DBDefaultName

	// Check if cache exists
	if !isFileExists(cachePath) {
		return "Cache not exists"
	}

	// Get file stats
	info, err := os.Stat(cachePath)
	if err != nil {
		log.Fatal("Can't get file stats:", err.Error())
	}

	s := Size(info.Size())
	date := info.ModTime()
	return fmt.Sprintf("last update(%d-%s-%02d %02d:%02d) size(%s)", date.Year(), date.Month().String(), date.Day(), date.Hour(), date.Minute(), s)
}
