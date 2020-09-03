package storage_store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB_store struct {
	path    string
	db      *leveldb.DB
	options string
	isOpen  bool

	// Sychronized action channel.
	actions chan map[string]interface{}
}

// Manage the concurent access of the db.
func (self *LevelDB_store) run() {
	for {
		select {
		case action := <-self.actions:
			if action["name"].(string) == "Open" {
				action["result"].(chan error) <- self.open(action["path"].(string))
			} else if action["name"].(string) == "SetItem" {
				if action["val"] != nil {
					action["result"].(chan error) <- self.setItem(action["key"].(string), action["val"].([]byte))
				} else {
					action["result"].(chan error) <- self.setItem(action["key"].(string), nil)
				}
			} else if action["name"].(string) == "GetItem" {
				val, err := self.getItem(action["key"].(string))
				action["results"].(chan map[string]interface{}) <- map[string]interface{}{"val": val, "err": err}
			} else if action["name"].(string) == "RemoveItem" {
				action["result"].(chan error) <- self.removeItem(action["key"].(string))
			} else if action["name"].(string) == "Clear" {
				action["result"].(chan error) <- self.clear()
			} else if action["name"].(string) == "Drop" {
				action["result"].(chan error) <- self.drop()
			} else if action["name"].(string) == "Close" {
				action["result"].(chan error) <- self.close()
				break // exit here.
			}

		}
	}
}

func NewLevelDB_store() *LevelDB_store {
	s := new(LevelDB_store)
	s.actions = make(chan map[string]interface{}, 0)
	go func() {
		s.run()
	}()
	return s
}

// In that case the parameter contain the path.
func (self *LevelDB_store) open(optionsStr string) error {
	if self.isOpen == true {
		return nil // the connection is already open.
	}
	fmt.Println("open store at path ", optionsStr)
	self.options = optionsStr

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

	self.isOpen = true
	return nil
}

// Close the store.
func (self *LevelDB_store) close() error {
	if self.isOpen == false {
		return nil
	}
	self.isOpen = false
	return self.db.Close()
}

// Set item
func (self *LevelDB_store) setItem(key string, val []byte) error {
	return self.db.Put([]byte(key), val, nil)
}

// Get item with a given key.
func (self *LevelDB_store) getItem(key string) ([]byte, error) {
	// Here I will make a little trick to give more flexibility to the tool...
	if strings.HasSuffix(key, "*") {
		// I will made use of iterator to ket the values
		values := "["
		iter := self.db.NewIterator(util.BytesPrefix([]byte(key[:len(key)-2])), nil)

		for ok := iter.Last(); ok; ok = iter.Prev() {
			if values != "[" {
				values += ","
			}
			values += string(iter.Value())
		}

		values += "]"

		iter.Release()
		return []byte(values), nil // I will return the stringnify value

	}

	return self.db.Get([]byte(key), nil)
}

// Remove an item or a range of items whit same path
func (self *LevelDB_store) removeItem(key string) error {
	if strings.HasSuffix(key, "*") {
		// I will made use of iterator to ket the values
		iter := self.db.NewIterator(util.BytesPrefix([]byte(key[:len(key)-1])), nil)
		for ok := iter.Last(); ok; ok = iter.Prev() {
			self.db.Delete([]byte(iter.Key()), nil)
		}
		iter.Release()

	}
	return self.db.Delete([]byte(key), nil)
}

// Clear the data store.
func (self *LevelDB_store) clear() error {
	err := self.Drop()
	if err != nil {
		return err
	}

	// Recreate the db files and connection.
	return self.Open(self.path)
}

// Drop the data store.
func (self *LevelDB_store) drop() error {
	// Close the db
	err := self.Close()
	if err != nil {
		return err
	}
	return os.RemoveAll(self.path)
}

//////////////////////// Synchronized LevelDB access ///////////////////////////

// Open the store with a give file path.
func (self *LevelDB_store) Open(path string) error {
	path = strings.ReplaceAll(path, "\\", "/")
	action := map[string]interface{}{"name": "Open", "result": make(chan error), "path": path}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Close the store.
func (self *LevelDB_store) Close() error {
	action := map[string]interface{}{"name": "Close", "result": make(chan error)}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Set item
func (self *LevelDB_store) SetItem(key string, val []byte) error {
	action := map[string]interface{}{"name": "SetItem", "result": make(chan error), "key": key, "val": val}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Get item with a given key.
func (self *LevelDB_store) GetItem(key string) ([]byte, error) {
	action := map[string]interface{}{"name": "GetItem", "results": make(chan map[string]interface{}, 0), "key": key}
	self.actions <- action
	results := <-action["results"].(chan map[string]interface{})
	if results["err"] != nil {
		return nil, results["err"].(error)
	}

	return results["val"].([]byte), nil
}

// Remove an item
func (self *LevelDB_store) RemoveItem(key string) error {
	action := map[string]interface{}{"name": "RemoveItem", "result": make(chan error), "key": key}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Clear the data store.
func (self *LevelDB_store) Clear() error {
	action := map[string]interface{}{"name": "Clear", "result": make(chan error)}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Drop the data store.
func (self *LevelDB_store) Drop() error {
	action := map[string]interface{}{"name": "Drop", "result": make(chan error)}
	self.actions <- action
	return <-action["result"].(chan error)
}
