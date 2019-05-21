package persistence_store

import (
	"context"
)

/**
 * Represent a data store interface.
 */
type Store interface {

	/**
	 * Create a database
	 */
	CreateDatabase(ctx context.Context, name string) error

	/**
	 * Delete a database
	 */
	DeleteDatabase(ctx context.Context, name string) error

	/**
	 * Create a Collection
	 */
	CreateCollection(ctx context.Context, database string, name string) error

	/**
	 * Delete collection
	 */
	DeleteCollection(ctx context.Context, database string, name string) error

	/**
	 * Connect to the data store.
	 */
	Connect(host string, port int32, user string, password string, database string, timeout int32, options_str string) error

	/**
	 * return error if the connection is unreachable.
	 */
	Ping(ctx context.Context) error

	/**
	 * return the number of entry in a table.
	 */
	Count(ctx context.Context, database string, collection string, query string, options string) (int64, error)

	/**
	 * Insert one result.
	 */
	InsertOne(ctx context.Context, database string, collection string, entity interface{}, options string) (interface{}, error)

	/**
	 * Insert many result at once.
	 */
	InsertMany(ctx context.Context, database string, collection string, entities []interface{}, options string) ([]interface{}, error)

	/**
	 * Find many values from a query
	 */
	Find(ctx context.Context, database string, collection string, query string, fields []string, options string) ([]interface{}, error)

	/**
	 * Find one result at time.
	 */
	FindOne(ctx context.Context, database string, collection string, query string, fields []string, options string) (interface{}, error)

	/**
	 * Update document that match a given condition whit a given value.
	 */
	Update(ctx context.Context, database string, collection string, query string, value string, options string) error

	/**
	 * Update document that match a given condition whit a given value.
	 */
	UpdateOne(ctx context.Context, database string, collection string, query string, value string, options string) error

	/**
	 * Replace one document that match a given condition whit a given value.
	 */
	ReplaceOne(ctx context.Context, database string, collection string, query string, value string, options string) error

	/**
	 * Remove one or more value depending of the query results.
	 */
	Delete(ctx context.Context, database string, collection string, query string, options string) error

	/**
	 * Remove one value depending of the query results.
	 */
	DeleteOne(ctx context.Context, database string, collection string, query string, options string) error
}
