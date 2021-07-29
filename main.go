package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Platform *string
	Source   *string
	DBSource *string
	Language *string
}

func main() {
	homeDir, err := getHomeDir()
	if err != nil {
		return
	}

	// Default flags
	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().BoolP("help", "h", false, "show this help message and exit")
	rootCmd.Flags().BoolP("version", "v", false, "show program's version number and exit")
	rootCmd.Flags().BoolP("update_cache", "u", false, "Update the local cache of pages and exit")
	rootCmd.Flags().StringP("platform", "p", "linux", "Override the operating system [linux, osx, sunos, windows, common]")
	rootCmd.Flags().BoolP("list", "l", false, "List all available commands for operating system")
	rootCmd.Flags().StringP("source", "s", getLocalPath(homeDir), "Override the default page source")
	rootCmd.Flags().StringP("language", "L", "en", "Override the default language")

	// Custom flags
	rootCmd.Flags().BoolP("global", "g", false, "Force to get info from official repository")
	rootCmd.Flags().StringP("db_source", "D", getDBPath(homeDir), "Override the default cache db path")

	// Setup
	cobra.CheckErr(rootCmd.Execute())
}

var rootCmd = &cobra.Command{
	Use:   "tldr",
	Short: "Golang command line client for tldr",
	Long:  "Golang command line client for tldr",
	Run:   main_call,
}

func main_call(cmd *cobra.Command, args []string) {
	// Print version
	v, err := cmd.Flags().GetBool("version")
	if err != nil {
		log.Fatal(err.Error())
	}
	if v {
		fmt.Println(getVersion())
		return
	}

	// Build config
	platform, err := cmd.Flags().GetString("platform")
	if err != nil {
		log.Fatal(err.Error())
	}
	source, err := cmd.Flags().GetString("source")
	if err != nil {
		log.Fatal(err.Error())
	}
	dbSource, err := cmd.Flags().GetString("db_source")
	if err != nil {
		log.Fatal(err.Error())
	}
	language, err := cmd.Flags().GetString("language")
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg := &Config{
		Platform: &platform,
		Source:   &source,
		DBSource: &dbSource,
		Language: &language,
	}

	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Print list of available commands
	list, err := cmd.Flags().GetBool("list")
	if err != nil {
		log.Fatal(err.Error())
	}
	if list {
		if global {
			commands, err2 := printGlobalList(cfg)
			if err2 != nil {
				err = err2
				return
			}
			fmt.Printf("%q\n", commands)
		} else {
			fmt.Printf("%q\n", printLocalList(cfg))
		}
		return
	}

	// Update cache
	update, err := cmd.Flags().GetBool("update_cache")
	if err != nil {
		log.Fatal(err.Error())
	}
	if update {
		if err = updateCache(cfg); err != nil {
			return
		}
		fmt.Println("Cache successfully updated")
		return
	}

	if len(args) != 1 {
		cmd.Help()
		return
	}

	// Get command name
	command := args[0]

	// If flag not set, first search result in local files
	if !global {
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
		switch err {
		case ErrorPageNotExists:
			// Do nothing
		default:
			log.Fatal(err)
		}
	}
	if len(page) > 0 {
		fmt.Println(GrayString("[global]"))
		fmt.Println(output(page))
		return
	}

	fmt.Printf("`%s` documentation is not available. Consider contributing Pull Request to https://github.com/tldr-pages/tldr\n", command)
}
