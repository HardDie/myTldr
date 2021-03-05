package main

import (
	"errors"
	"io/ioutil"
	"log"
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

func checkRemote(platform, language, name string) (page []string, err error) {
	url := buildRemotePath(platform, language)
	res, err := http.Get(url + "/" + name + ".md")
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = errors.New("file not found")
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	page = strings.Split(string(data), "\n")
	return
}
