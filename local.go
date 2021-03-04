package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

const (
	FilesDefaultPath = ".my_scripts/.tldr"
)

func getLocalPath() string {
	user, err := user.Current()
	if err != nil {
		// Application can't continue
		os.Exit(1)
	}
	return user.HomeDir + "/" + FilesDefaultPath
}

func buildLocalPath(source, language string) string {
	folder := "pages"
	if language != "en" {
		folder += "." + language
	}
	return source + "/" + folder
}

func checkLocal(source, language string, name string) (page []string, err error) {
	fileName := buildLocalPath(source, language) + "/" + name + ".md"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	page = strings.Split(string(data), "\n")
	return
}
