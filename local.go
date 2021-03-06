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

func buildLocalPath(source, platform, language string) string {
	folder := "pages"
	if language != "en" {
		folder += "." + language
	}
	return source + "/" + folder + "/" + platform
}

func checkLocal(source, platform, language, name string) (page []string, err error) {
	// Build path to the local pages
	fileName := buildLocalPath(source, platform, language) + "/" + name + ".md"
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
