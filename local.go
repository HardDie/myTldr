package main

import (
	"io/ioutil"
	"strings"
)

const (
	FilesDefaultPath = ".my_scripts/.tldr"
)

func getLocalPath(homeDir string) (path string) {
	return homeDir + "/" + FilesDefaultPath
}

func buildLocalPath(cfg *Config, platform string) string {
	folder := "pages"
	if *cfg.Language != "en" {
		folder += "." + *cfg.Language
	}
	return *cfg.Source + "/" + folder + "/" + platform
}

func checkLocal(cfg *Config, name string) (page []string, err error) {
	platforms := []string{PlatformCommon, *cfg.Platform}

	// Check file for all platforms
	for _, platform := range platforms {
		// Build path to the local pages
		fileName := buildLocalPath(cfg, platform) + "/" + name + ".md"
		// If page not exist, just return
		if isFileExists(fileName) {
			var data []byte
			// Read page data
			data, err = ioutil.ReadFile(fileName)
			if err != nil {
				return
			}
			// Split text to lines
			page = strings.Split(string(data), "\n")
		}
	}
	return
}
