package main

import (
	"errors"
	"io/ioutil"
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

func buildBucketName(cfg *Config) string {
	if *cfg.Language == "en" {
		return *cfg.Platform
	}
	return *cfg.Platform + "." + *cfg.Language
}

func checkCache(cfg *Config, name string) (page []string, err error) {
	// If DB not exist, return
	if !isFileExists(*cfg.DBSource) {
		err = ErrorCacheNotExists
		return
	}

	// // Open DB
	bw, err := NewBoltWrapper(*cfg.DBSource + "/" + DBDefaultName)
	if err != nil {
		return
	}
	defer func() { _ = bw.Close() }()

	// Build bucket name from input values
	bucketName := buildBucketName(cfg)
	data, err := bw.GetDataFromBucket(bucketName, name)
	if err != nil {
		return
	}

	// Split data to lines
	page = strings.Split(data, "\n")
	return
}

func putCache(cfg *Config, name string, data []byte) (err error) {
	bw, err := NewBoltWrapper(*cfg.DBSource + "/" + DBDefaultName)
	if err != nil {
		return
	}
	defer func() { _ = bw.Close() }()

	bucketName := buildBucketName(cfg)
	if err = bw.CreateBucketIfNotExists(bucketName); err != nil {
		return
	}
	if err = bw.PutDataToBucket(bucketName, name, data); err != nil {
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
			if err = putCache(tmpCfg, command, data); err != nil {
				return err
			}
		}
	}
	return nil
}
