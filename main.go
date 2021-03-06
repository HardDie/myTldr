package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Process flags or print information if required
	source, platform, language, dbSource, global, done, err := handleFlags()
	if err != nil {
		log.Fatal(err)
	}
	if done {
		return
	}

	// Get command name
	command := flag.Args()[0]

	// If flag not set, first search result in local files
	if global == false {
		// Get page from local folder
		var page []string
		page, err = checkLocal(source, platform, language, command)
		if err != nil {
			log.Fatal(err)
		}
		if len(page) > 0 {
			fmt.Println(GrayString("[local]"))
			fmt.Println(output(page))
			return
		}

	}

	// Try to get page from cache
	page, err := checkCache(dbSource, platform, language, command)
	if err != nil {
		log.Fatal(err)
	}
	if len(page) > 0 {
		fmt.Println(GrayString("[cache]"))
		fmt.Println(output(page))
		return
	}

	// Try to find page in official repository
	page, err = checkRemote(platform, language, command, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	if len(page) > 0 {
		fmt.Println(GrayString("[global]"))
		fmt.Println(output(page))
		return
	}

	fmt.Printf("`%s` documentation is not available. Consider contributing Pull Request to https://github.com/tldr-pages/tldr\n", command)
}

func handleFlags() (source, platform, language, dbSource string, global, done bool, err error) {
	localPath, err := getLocalPath()
	if err != nil {
		return
	}
	dbPath, err := getDBPath()
	if err != nil {
		return
	}

	fVersion := flag.Bool("version", false, "show program's version number and exit")
	fUpdateCache := flag.Bool("update_cache", false, "Update the local cache of pages and exit")
	fPlatform := flag.String("platform", "linux", "Override the operating system [linux, osx, sunos, windows, common]")
	fList := flag.Bool("list", false, "List all available commands for operating system")
	fSource := flag.String("source", localPath, "Override the default page source")
	fLanguage := flag.String("language", "en", "Override the default language")
	fGlobal := flag.Bool("global", false, "Force to get info from official repository")
	fDBSource := flag.String("dbSource", dbPath, "Override the default cache db path")
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
		done = true
		return
	case *fList:
		fmt.Printf("%q\n", printList(*fSource, *fPlatform, *fLanguage))
		done = true
		return
	}

	source = *fSource
	platform = *fPlatform
	language = *fLanguage
	global = *fGlobal
	dbSource = *fDBSource
	if len(flag.Args()) != 1 {
		flag.Usage()
		done = true
		return
	}
	return
}
