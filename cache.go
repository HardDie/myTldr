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

func getDBPath() string {
	user, err := user.Current()
	if err != nil {
		// Application can't continue
		os.Exit(1)
	}
	return user.HomeDir + "/" + DBDefaultPath
}

func buildBucketName(platform, language string) string {
	name := platform
	if language != "en" {
		name += "." + language
	}
	return name
}

func openBoldDB(source string) (db *bolt.DB, clean func(), err error) {
	// Create folder if not exists
	err = os.MkdirAll(source, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
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

func checkCache(source, platform, language, name string) (page []string, err error) {
	db, clean, err := openBoldDB(source)
	if err != nil {
		return
	}
	defer clean()

	var data []byte
	bucketName := buildBucketName(platform, language)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errors.New("Data not exists")
		}
		data = b.Get([]byte(name))
		if data == nil {
			return errors.New("Data not exists")
		}
		return nil
	})
	if err != nil {
		return
	}

	if len(data) > 0 {
		page = strings.Split(string(data), "\n")
		return
	}

	err = errors.New("file not found")
	return
}

func putCache(source, platform, language, name string, data []byte) {
	db, clean, err := openBoldDB(source)
	if err != nil {
		return
	}
	defer clean()

	bucketName := buildBucketName(platform, language)
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(name), data)
	})
	if err != nil {
		log.Fatal(err)
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
