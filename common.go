package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
)

func isFileExists(path string) (isExist bool) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func httpGet(url string) (result []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	// Close the body
	defer func() { _ = resp.Body.Close() }()

	// If response not OK, it means page not exists
	if resp.StatusCode != 200 {
		return
	}

	// Read data from body
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func downloadZip(url string) (zipReader *zip.Reader, err error) {
	// Download the ZIP file
	zipFile, err := httpGet(url)
	if err != nil {
		return
	}

	// Turn this array into a zip reader
	zipReader, err = zip.NewReader(
		bytes.NewReader(zipFile),
		int64(len(zipFile)),
	)
	if err != nil {
		return
	}
	return
}

func getHomeDir() (homeDir string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	homeDir = usr.HomeDir
	return
}