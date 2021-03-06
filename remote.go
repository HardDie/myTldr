package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	RemoteBaseURL = "https://raw.githubusercontent.com/tldr-pages/tldr/master"
)

func buildRemotePath(cfg *Config) string {
	folder := "pages"
	if cfg.Language != "en" {
		folder += "." + cfg.Language
	}
	return RemoteBaseURL + "/" + folder + "/" + cfg.Platform
}

func checkRemote(cfg *Config, name string) (page []string, err error) {
	// Build url to possible tldr page
	url := buildRemotePath(cfg)
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
	err = putCache(cfg, name, data)
	if err != nil {
		return
	}
	// Split text on lines
	page = strings.Split(string(data), "\n")
	return
}
