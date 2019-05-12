package persistence_store

import "context"

/**
 * Represent a data store interface.
 */
type Store interface {

	/**
	 * Connect to the data store.
	 */
	Connect(host string, port int32, user string, password string, database string, timeout int32, options_str string) error

	/**
	 * return error if the connection is unreachable.
	 */
	Ping(ctx context.Context) error

	/**
	 * Insert one result.
	 */
	InsertOne(ctx context.Context, database string, collection string, entity interface{}) (interface{}, error)

	/**
	 * Insert many result at once.
	 */
	InsertMany(ctx context.Context, database string, collection string, entities []interface{}) ([]interface{}, error)
}
