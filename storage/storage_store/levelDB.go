package storage_store

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB_store struct {
	path    string
	db      *leveldb.DB
	options string
	isOpen  bool
}

func NewLevelDB_store() *LevelDB_store {
	s := new(LevelDB_store)

	return s
}

// In that case the parameter contain the path.
func (self *LevelDB_store) Open(optionsStr string) error {
	if self.isOpen == true {
		return nil // the connection is already open.
	}

	self.options = optionsStr
	log.Println("--> try to open ", self.path, "db is open")
	var err error
	if len(self.path) == 0 {
		if len(optionsStr) == 0 {
			return errors.New("store path and store name must be given as options!")
		}

		options := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(optionsStr), &options)

		if options["path"] == nil {
			return errors.New("no store path was given in connection option!")
		}

		if options["name"] == nil {
			return errors.New("no store name was given in connection option!")
		}

		self.path = options["path"].(string) + string(os.PathSeparator) + options["name"].(string)

	}
	// Open the store.
	self.db, err = leveldb.OpenFile(self.path, nil)
	if err != nil {
		return err
	}
	log.Println("--> ", self.path, "db is open")
	self.isOpen = true
	return nil
}

// Close the store.
func (self *LevelDB_store) Close() error {
	if self.isOpen == false {
		return nil
	}
	self.isOpen = false
	return self.db.Close()
}

// Set item
func (self *LevelDB_store) SetItem(key string, val []byte) error {
	return self.db.Put([]byte(key), val, nil)
}

// Get item with a given key.
func (self *LevelDB_store) GetItem(key string) ([]byte, error) {
	return self.db.Get([]byte(key), nil)
}

// Remove an item
func (self *LevelDB_store) RemoveItem(key string) error {
	return self.db.Delete([]byte(key), nil)
}

// Clear the data store.
func (self *LevelDB_store) Clear() error {
	err := self.Drop()
	if err != nil {
		return err
	}

	// Recreate the db files and connection.
	return self.Open(self.path)
}

// Drop the data store.
func (self *LevelDB_store) Drop() error {
	// Close the db
	err := self.Close()
	if err != nil {
		return err
	}

	return os.RemoveAll(self.path)
}
