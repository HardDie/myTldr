package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	source, language, err := handleFlags()
	if err != nil {
		return
	}

	command := flag.Args()[0]
	page, err := checkLocal(source, language, command)
	if err == nil {
		fmt.Println(output(page))
		return
	}

	fmt.Printf("`%s` documentation is not available. Consider contributing Pull Request to https://github.com/tldr-pages/tldr\n", command)
}

func handleFlags() (source, language string, err error) {
	fVersion := flag.Bool("version", false, "show program's version number and exit")
	fUpdateCache := flag.Bool("update_cache", false, "Update the local cache of pages and exit")
	fPlatform := flag.String("platform", "linux", "Override the operating system [linux, osx, sunos, windows, common]")
	fList := flag.Bool("list", false, "List all available commands for operating system")
	fSource := flag.String("source", getLocalPath(), "Override the default page source")
	fLanguage := flag.String("language", "en", "Override the default language")
	flag.Usage = func() {
		fmt.Printf("usage: %s [options] command\n\n", os.Args[0])
		fmt.Println("Go command line client for tldr\n")
		fmt.Println("optional arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()

	_ = fUpdateCache
	_ = fPlatform
	_ = fList

	switch {
	case *fVersion:
		fmt.Println(getVersion())
		os.Exit(0)
	case *fList:
		fmt.Printf("%q\n", printList(*fSource, *fLanguage))
		os.Exit(0)
	}

	source = *fSource
	language = *fLanguage
	if len(flag.Args()) != 1 {
		flag.Usage()
		err = errors.New("not enough arguments")
		return
	}
	return
}
