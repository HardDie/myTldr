package main

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

var (
	ErrorCannotFindBucket = errors.New("can't find the bucket")
	ErrorPutData          = errors.New("can't put data to bucket")
	ErrorInvalidKey       = errors.New("such key not found")
	ErrorDBClosed         = errors.New("DB is closed")
)

type BoltWrapper struct {
	db *bolt.DB
}

func NewBoltWrapper(dbpath string) (bw *BoltWrapper, err error) {
	db, err := bolt.Open(dbpath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return
	}
	bw = &BoltWrapper{
		db: db,
	}
	return
}

func (bw *BoltWrapper) Close() (err error) {
	if bw.db == nil {
		err = ErrorDBClosed
		return
	}

	err = bw.db.Close()
	bw.db = nil
	return
}

func (bw *BoltWrapper) CreateBucketIfNotExists(name string) (err error) {
	if bw.db == nil {
		err = ErrorDBClosed
		return
	}

	return bw.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
}

func (bw *BoltWrapper) PutDataToBucket(bucketName, key string, value []byte) (err error) {
	if bw.db == nil {
		err = ErrorDBClosed
		return
	}

	return bw.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return ErrorCannotFindBucket
		}
		err = b.Put([]byte(key), value)
		if err != nil {
			return ErrorPutData
		}
		return nil
	})
}

func (bw *BoltWrapper) GetDataFromBucket(bucketName, key string) (value string, err error) {
	if bw.db == nil {
		err = ErrorDBClosed
		return
	}

	err = bw.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return ErrorCannotFindBucket
		}
		data := b.Get([]byte(key))
		if data == nil {
			return ErrorInvalidKey
		}
		value = string(data)
		return nil
	})
	return
}

func (bw *BoltWrapper) GetKeysFromBucket(bucketName string) (keys []string, err error) {
	err = bw.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return ErrorCannotFindBucket
		}
		return b.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})
	if err != nil {
		return
	}
	return
}

func (bw *BoltWrapper) PrintAllData() (err error) {
	if bw.db == nil {
		err = ErrorDBClosed
		return
	}

	return bw.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			log.Println("bucket:", string(name))
			return b.ForEach(func(k, v []byte) error {
				log.Println("   key:", string(k), "value:", strings.Split(string(v), "\n")[0])
				return nil
			})
		})
	})
}
