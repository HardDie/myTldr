package main

import (
	"os"
)

func isFileExists(path string) (isExist bool) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
