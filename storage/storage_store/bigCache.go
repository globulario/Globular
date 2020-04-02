package storage_store

import (
	"time"

	"encoding/json"

	"github.com/allegro/bigcache"
)

// Implement the storage service with big store.
type BigCache_store struct {
	cache *bigcache.BigCache // The actual cache.

	// Sychronization.
	actions chan map[string]interface{}
}

func (self *BigCache_store) run() {
	for {
		select {
		case action := <-self.actions:

			if action["name"].(string) == "Open" {
				action["result"].(chan error) <- self.open(action["options"].(string))
			} else if action["name"].(string) == "SetItem" {
				action["result"].(chan error) <- self.cache.Set(action["key"].(string), action["val"].([]byte))
			} else if action["name"].(string) == "GetItem" {
				val, err := self.cache.Get(action["key"].(string))
				action["results"].(chan map[string]interface{}) <- map[string]interface{}{"val": val, "err": err}
			} else if action["name"].(string) == "RemoveItem" {
				action["result"].(chan error) <- self.cache.Delete(action["key"].(string))
			} else if action["name"].(string) == "Clear" {
				action["result"].(chan error) <- self.cache.Reset()
			} else if action["name"].(string) == "Drop" {
				action["result"].(chan error) <- self.cache.Reset()
			} else if action["name"].(string) == "Close" {
				action["result"].(chan error) <- self.cache.Close()
				break // exit here.
			}

		}
	}
}

// Use it to use the store.
func NewBigCache_store() *BigCache_store {
	s := new(BigCache_store)
	s.actions = make(chan map[string]interface{}, 0)

	go func() {
		s.run()
	}()

	return s
}

func (self *BigCache_store) Open(options string) error {
	action := map[string]interface{}{"name": "Open", "result": make(chan error), "options": options}
	self.actions <- action
	return <-action["result"].(chan error)
}

func (self *BigCache_store) open(optionsStr string) error {
	var config bigcache.Config
	var err error
	if len(optionsStr) == 0 {
		config = bigcache.Config{
			// number of shards (must be a power of 2)
			Shards: 1024,
			// time after which entry can be evicted
			LifeWindow: 2 * time.Minute,
			// rps * lifeWindow, used only in initial memory allocation
			MaxEntriesInWindow: 1000 * 10 * 60,
			// max entry size in bytes, used only in initial memory allocation
			MaxEntrySize: 500,
			// prints information about additional memory allocation
			Verbose: true,
			// cache will not allocate more memory than this limit, value in MB
			// if value is reached then the oldest entries can be overridden for the new ones
			// 0 value means no size limit
			HardMaxCacheSize: 4000,
			// callback fired when the oldest entry is removed because of its
			// expiration time or no space left for the new entry. Default value is nil which
			// means no callback and it prevents from unwrapping the oldest entry.
			OnRemove: func(key string, data []byte) {
				/** Nothing here **/
			},
		}
	} else {
		// init the config from string.
		json.Unmarshal([]byte(optionsStr), &config)
	}

	// init the underlying cache.
	self.cache, err = bigcache.NewBigCache(config)

	return err
}

// Close the store.
func (self *BigCache_store) Close() error {
	action := map[string]interface{}{"name": "Close", "result": make(chan error)}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Set item
func (self *BigCache_store) SetItem(key string, val []byte) error {
	action := map[string]interface{}{"name": "SetItem", "result": make(chan error), "key": key, "val": val}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Get item with a given key.
func (self *BigCache_store) GetItem(key string) ([]byte, error) {
	action := map[string]interface{}{"name": "GetItem", "results": make(chan map[string]interface{}, 0), "key": key}
	self.actions <- action
	results := <-action["results"].(chan map[string]interface{})
	if results["err"] != nil {
		return nil, results["err"].(error)
	}

	return results["val"].([]byte), nil
}

// Remove an item
func (self *BigCache_store) RemoveItem(key string) error {
	action := map[string]interface{}{"name": "RemoveItem", "result": make(chan error), "key": key}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Clear the data store.
func (self *BigCache_store) Clear() error {
	action := map[string]interface{}{"name": "Clear", "result": make(chan error)}
	self.actions <- action
	return <-action["result"].(chan error)
}

// Drop the data store.
func (self *BigCache_store) Drop() error {
	action := map[string]interface{}{"name": "Drop", "result": make(chan error)}
	self.actions <- action
	return <-action["result"].(chan error)
}
