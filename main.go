package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	source, platform, language, global, err := handleFlags()
	if err != nil {
		return
	}

	// Get command name
	command := flag.Args()[0]

	if global == false {
		// Get page from local folder
		page, err := checkLocal(source, platform, language, command)
		if err == nil {
			fmt.Println(GrayString("[local]"))
			fmt.Println(output(page))
			return
		}
	}

	// Try to find page in official repository
	page, err := checkRemote(platform, language, command)
	if err == nil {
		fmt.Println(GrayString("[global]"))
		fmt.Println(output(page))
		return
	}

	fmt.Printf("`%s` documentation is not available. Consider contributing Pull Request to https://github.com/tldr-pages/tldr\n", command)
}

func handleFlags() (source, platform, language string, global bool, err error) {
	fVersion := flag.Bool("version", false, "show program's version number and exit")
	fUpdateCache := flag.Bool("update_cache", false, "Update the local cache of pages and exit")
	fPlatform := flag.String("platform", "linux", "Override the operating system [linux, osx, sunos, windows, common]")
	fList := flag.Bool("list", false, "List all available commands for operating system")
	fSource := flag.String("source", getLocalPath(), "Override the default page source")
	fLanguage := flag.String("language", "en", "Override the default language")
	fGlobal := flag.Bool("global", false, "Force to get info from official repository")
	flag.Usage = func() {
		fmt.Printf("usage: %s [options] command\n\n", os.Args[0])
		fmt.Println("Go command line client for tldr\n")
		fmt.Println("optional arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()
	_ = fUpdateCache

	switch {
	case *fVersion:
		fmt.Println(getVersion())
		os.Exit(0)
	case *fList:
		fmt.Printf("%q\n", printList(*fSource, *fPlatform, *fLanguage))
		os.Exit(0)
	}

	source = *fSource
	platform = *fPlatform
	language = *fLanguage
	global = *fGlobal
	if len(flag.Args()) != 1 {
		flag.Usage()
		err = errors.New("not enough arguments")
		return
	}
	return
}
