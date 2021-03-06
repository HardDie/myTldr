package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	DBDefaultPath = ".cache/tldr"
)

func getDBPath() (path string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	path = usr.HomeDir + "/" + DBDefaultPath
	return
}

func buildBucketName(cfg *Config) string {
	if cfg.Language == "en" {
		return cfg.Platform
	} else {
		return cfg.Platform + "." + cfg.Language
	}
}

func openBoldDB(source string) (db *bolt.DB, clean func(), err error) {
	// Create folder if not exists
	if !isFileExists(source) {
		err = os.MkdirAll(source, 0777)
		if err != nil {
			return
		}
	}

	db, err = bolt.Open(source+"/"+"pages.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return
	}

	clean = func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func checkCache(cfg *Config, name string) (page []string, err error) {
	// If DB not exist, return
	if !isFileExists(cfg.Source) {
		return
	}

	// Open DB
	db, clean, err := openBoldDB(cfg.Source)
	if err != nil {
		return
	}
	defer clean()

	var data []byte
	// Build bucket name from input values
	bucketName := buildBucketName(cfg)
	err = db.View(func(tx *bolt.Tx) error {
		// Check if bucket exists
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errors.New("Data not exists")
		}
		// Check if data exists
		data = b.Get([]byte(name))
		if data == nil {
			return errors.New("Data not exists")
		}
		return nil
	})
	if err != nil {
		err = nil
		return
	}

	// Split data to lines
	page = strings.Split(string(data), "\n")
	return
}

func putCache(cfg *Config, name string, data []byte) (err error) {
	db, clean, err := openBoldDB(cfg.Source)
	if err != nil {
		return
	}
	defer clean()

	bucketName := buildBucketName(cfg)
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(name), data)
	})
	if err != nil {
		return
	}
	return
}

func printAllCache(source string) {
	db, clean, err := openBoldDB(source)
	if err != nil {
		return
	}
	defer clean()

	err = db.View(func(tx *bolt.Tx) error {
		_ = tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println("bucket:", string(name))
			return b.ForEach(func(k, v []byte) error {
				fmt.Println("   key:", string(k), "value:", strings.Split(string(v), "\n")[0])
				return nil
			})
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return
}
