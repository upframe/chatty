package boltdb

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

// The list of bucket names to require in the database.
var bucketNames = []string{
	"teams",
}

// Setup opens the database at file and sets it up.
func Setup(file string) error {
	var err error
	db, err = bolt.Open(file, 0600, nil)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range bucketNames {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket %s: %v", bucket, err)
			}
		}
		return nil
	})
	return err
}

// saveToDB saves val by key into bucket by gob-encoding it.
func saveToDB(bucket, key string, val interface{}) error {
	enc, err := jsonEncode(val)
	if err != nil {
		return fmt.Errorf("error encoding for database: %v", err)
	}
	return saveToDBRaw(bucket, []byte(key), enc)
}

// saveToDBRaw saves the value with key in bucket.
func saveToDBRaw(bucket string, key, val []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(key, val)
	})
}

// loadFromBucket loads key from bucket into into, decoded.
func loadFromBucket(into interface{}, bucket *bolt.Bucket, key []byte) error {
	v := bucket.Get([]byte(key))
	if v != nil {
		return jsonDecode(v, into)
	}
	return nil
}

// loadFromDB loads key from bucket into into, decoded.
func loadFromDB(into interface{}, bucket, key string) error {
	return db.View(func(tx *bolt.Tx) error {
		return loadFromBucket(into, tx.Bucket([]byte(bucket)), []byte(key))
	})
}

// loadFromDBRaw loads the value from bucket at key, no decoding.
func loadFromDBRaw(bucket, key string) ([]byte, error) {
	var v []byte
	return v, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get([]byte(key))
		return nil
	})
}

// isUnique returns true if key is not found in bucket.
func isUnique(bucket, key string) bool {
	var unique bool
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		unique = b.Get([]byte(key)) == nil
		return nil
	})
	return unique
}

// jsonEncode json encodes value.
func jsonEncode(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(value)
	return buf.Bytes(), err
}

// jsonDecode gob decodes buf into into.
func jsonDecode(buf []byte, into interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(buf))
	return dec.Decode(into)
}
