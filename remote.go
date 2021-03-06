package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	RemoteBaseURL = "https://raw.githubusercontent.com/tldr-pages/tldr/master"
)

func buildRemotePath(platform, language string) string {
	folder := "pages"
	if language != "en" {
		folder += "." + language
	}
	return RemoteBaseURL + "/" + folder + "/" + platform
}

func checkRemote(platform, language, name, dbSource string) (page []string, err error) {
	// Build url to possible tldr page
	url := buildRemotePath(platform, language)
	// Get page from official repository
	res, err := http.Get(url + "/" + name + ".md")
	if err != nil {
		return
	}
	// If github response not OK, it means page not exists
	if res.StatusCode != 200 {
		return
	}
	// Read page text
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	// Put new page to cache or update old
	err = putCache(dbSource, platform, language, name, data)
	if err != nil {
		return
	}
	// Split text on lines
	page = strings.Split(string(data), "\n")
	return
}
