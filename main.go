package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Platform *string
	Source   *string
	DBSource *string
	Language *string
	Global   *bool
}

func main() {
	log.SetFlags(log.Lshortfile)
	// Process flags or print information if required
	cfg, done, err := handleFlags()
	if err != nil {
		log.Fatal(err)
	}
	if done {
		return
	}

	// Get command name
	command := flag.Args()[0]

	// If flag not set, first search result in local files
	if !*cfg.Global {
		// Get page from local folder
		var page []string
		page, err = checkLocal(cfg, command)
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
	page, err := checkCache(cfg, command)
	if err != nil {
		switch err {
		case ErrorInvalidKey:
		// Do nothing
		case ErrorCacheNotExists:
			fmt.Println("Cache is empty, you can download it: " + os.Args[0] + " -update_cache")
		default:
			log.Fatal(err)
		}
	}
	if len(page) > 0 {
		fmt.Println(GrayString("[cache]"))
		fmt.Println(output(page))
		return
	}

	// Try to find page in official repository
	page, err = checkRemote(cfg, command)
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

func handleFlags() (cfg *Config, done bool, err error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return
	}

	cfg = &Config{
		Platform: flag.String("platform", "linux", "Override the operating system [linux, osx, sunos, windows, common]"),
		Source:   flag.String("source", getLocalPath(homeDir), "Override the default page source"),
		DBSource: flag.String("dbSource", getDBPath(homeDir), "Override the default cache db path"),
		Language: flag.String("language", "en", "Override the default language"),
		Global:   flag.Bool("global", false, "Force to get info from official repository"),
	}
	fVersion := flag.Bool("version", false, "show program's version number and exit")
	fUpdateCache := flag.Bool("update_cache", false, "Update the local cache of pages and exit")
	fList := flag.Bool("list", false, "List all available commands for operating system")

	flag.Usage = func() {
		fmt.Println(getVersion())
		fmt.Println()
		fmt.Printf("usage: %s [options] command\n\n", os.Args[0])
		fmt.Printf("cache info: %s\n\n", getCacheInfo(cfg))
		fmt.Println("Go command line client for tldr")
		fmt.Println()
		fmt.Println("optional arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()

	switch {
	case *fVersion:
		fmt.Println(getVersion())
		done = true
		return
	case *fList:
		if *cfg.Global {
			commands, err2 := printGlobalList(cfg)
			if err2 != nil {
				err = err2
				return
			}
			fmt.Printf("%q\n", commands)
		} else {
			fmt.Printf("%q\n", printLocalList(cfg))
		}
		done = true
		return
	case *fUpdateCache:
		if err = updateCache(cfg); err != nil {
			return
		}
		fmt.Println("Cache successfully updated")
		done = true
		return
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		done = true
		return
	}
	return
}
