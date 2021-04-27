package main

import (
	"strings"
)

const (
	RemoteBaseURL = "https://raw.githubusercontent.com/tldr-pages/tldr/master"
)

func buildRemotePath(cfg *Config, platform string) string {
	folder := "pages"
	if *cfg.Language != "en" {
		folder += "." + *cfg.Language
	}
	return RemoteBaseURL + "/" + folder + "/" + platform
}

func checkRemote(cfg *Config, name string) (page []string, err error) {
	platforms := []string{PlatformCommon, *cfg.Platform}

	var data []byte
	for _, platform := range platforms {
		// Build url to possible tldr page
		url := buildRemotePath(cfg, platform)

		// Get page from official repository
		data, err = httpGet(url + "/" + name + ".md")
		if err == ErrorPageNotExists {
			// If such page not exists, just check another platform
			continue
		}
		if err != nil {
			return
		}
		break
	}
	if len(data) == 0 {
		return
	}

	bw, err := openCache(cfg)
	if err != nil {
		return
	}
	defer func() { _ = bw.Close() }()

	// Put new page to cache or update old
	err = putCache(cfg, bw, name, data)
	if err != nil {
		return
	}
	// Split text on lines
	page = strings.Split(string(data), "\n")
	return
}
