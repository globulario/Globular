package storage_store

import (
	"time"

	"encoding/json"

	"github.com/allegro/bigcache"
)

// Implement the storage service with big store.
type BigCache_store struct {
	cache *bigcache.BigCache // The actual cache.
}

func NewBigCache_store() *BigCache_store {
	s := new(BigCache_store)
	return s
}

func (self *BigCache_store) Open(optionsStr string) error {
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
	return self.cache.Close()
}

// Set item
func (self *BigCache_store) SetItem(key string, val []byte) error {
	return self.cache.Set(key, val)
}

// Get item with a given key.
func (self *BigCache_store) GetItem(key string) ([]byte, error) {
	return self.cache.Get(key)
}

// Remove an item
func (self *BigCache_store) RemoveItem(key string) error {
	return self.cache.Delete(key)
}

// Clear the data store.
func (self *BigCache_store) Clear() error {
	return self.cache.Reset()
}

// Drop the data store.
func (self *BigCache_store) Drop() error {
	return self.cache.Reset()
}
