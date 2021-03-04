package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	RemoteBaseURL = "https://raw.githubusercontent.com/tldr-pages/tldr/master/"
)

func buildRemotePath(language string) string {
	folder := "pages"
	if language != "en" {
		folder += "." + language
	}
	return RemoteBaseURL + folder
}

func checkCategory(path, category, name string) (page []string, err error) {
	res, err := http.Get(path + "/" + category + "/" + name + ".md")
	if err != nil {
		return
	}
	if res.StatusCode == 404 {
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

func checkRemote(language, name string) (page []string, err error) {
	path := buildRemotePath(language)
	if page, err = checkCategory(path, "common", name); err == nil {
		return
	}
	if page, err = checkCategory(path, "linux", name); err == nil {
		return
	}

	err = errors.New("no required command")
	return
}
