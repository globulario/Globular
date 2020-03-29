package storage_store

/**
 * A key value data store.
 */
type Store interface {
	// Open the store
	Open(optionsStr string) error

	// Close the store
	Close() error

	// Set item
	SetItem(key string, val []byte) error

	// Get item with a given key.
	GetItem(key string) ([]byte, error)

	// Remove an item
	RemoveItem(key string) error

	// Clear the data store.
	Clear() error

	// Drop the data store.
	Drop() error
}
