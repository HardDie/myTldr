package main

import (
	"io/ioutil"
	"os/user"
	"strings"
)

const (
	FilesDefaultPath = ".my_scripts/.tldr"
)

func getLocalPath() (path string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	path = usr.HomeDir + "/" + FilesDefaultPath
	return
}

func buildLocalPath(cfg *Config) string {
	folder := "pages"
	if cfg.Language != "en" {
		folder += "." + cfg.Language
	}
	return cfg.Source + "/" + folder + "/" + cfg.Platform
}

func checkLocal(cfg *Config, name string) (page []string, err error) {
	// Build path to the local pages
	fileName := buildLocalPath(cfg) + "/" + name + ".md"
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
	return
}
