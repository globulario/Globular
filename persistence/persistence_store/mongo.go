package persistence_store

import (
	"context"
	//	"fmt"
	// "log"

	"strconv"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"encoding/json"

	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * Implementation of the the Store interface with mongo db.
 */
type MongoStore struct {
	client *mongo.Client
}

/**
 * Connect to the remote/local mongo server
 * TODO add more connection options via the option_str and options package.
 */
func (self *MongoStore) Connect(host string, port int32, user string, password string, database string, timeout int32, options_str string) (err error) {

	// basic connection string to begin with.
	connectionStr := "mongodb://" + host + ":" + strconv.Itoa(int(port))

	self.client, err = mongo.NewClient(options.Client().ApplyURI(connectionStr))

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	self.client.Connect(ctx)

	if len(database) > 0 {
		// In that case if the database dosent exist I will return an error.
		if self.client.Database(database) == nil {
			return errors.New("No database with name " + database + " exist on this store.")
		}
	}

	return err
}

/**
 * Return the nil on success.
 */
func (self *MongoStore) Ping(ctx context.Context) error {
	return self.client.Ping(ctx, nil)
}

/**
 * return the number of entry in a table.
 */
func (self *MongoStore) Count(ctx context.Context, database string, collection string, query string, optionsStr string) (int64, error) {
	var opts *options.CountOptions
	if len(optionsStr) > 0 {
		opts = new(options.CountOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return int64(0), err
		}
	}

	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return int64(0), err
	}

	count, err := self.client.Database(database).Collection(collection).CountDocuments(ctx, q, opts)
	return count, err
}

func (self *MongoStore) CreateDatabase(ctx context.Context, name string) error {
	return errors.New("MongoDb will create your database at first insert.")
}

/**
 * Delete a database
 */
func (self *MongoStore) DeleteDatabase(ctx context.Context, name string) error {
	return self.client.Database(name).Drop(ctx)
}

/**
 * Create a Collection
 */
func (self *MongoStore) CreateCollection(ctx context.Context, database string, name string) error {
	return errors.New("MongoDb will create your collection at first insert.")
}

/**
 * Delete collection
 */
func (self *MongoStore) DeleteCollection(ctx context.Context, database string, name string) error {
	err := self.client.Database(name).Collection(name).Drop(ctx)
	return err
}

//////////////////////////////////////////////////////////////////////////////////
// Insert
//////////////////////////////////////////////////////////////////////////////////
/**
 * Insert one value in the store.
 */
func (self *MongoStore) InsertOne(ctx context.Context, database string, collection string, entity interface{}, optionsStr string) (interface{}, error) {

	var opts *options.InsertOneOptions
	if len(optionsStr) > 0 {
		opts = new(options.InsertOneOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return int64(0), err
		}
	}

	// Get the collection object.
	collection_ := self.client.Database(database).Collection(collection)

	result, err := collection_.InsertOne(ctx, entity, opts)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

/**
 * Insert many results at time.
 */
func (self *MongoStore) InsertMany(ctx context.Context, database string, collection string, entities []interface{}, optionsStr string) ([]interface{}, error) {

	var opts *options.InsertManyOptions
	if len(optionsStr) > 0 {
		opts = new(options.InsertManyOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return nil, err
		}
	}

	// Get the collection object.
	collection_ := self.client.Database(database).Collection(collection)

	// return self.client.Ping(ctx, nil)
	insertManyResult, err := collection_.InsertMany(ctx, entities, opts)
	if err != nil {
		return nil, err
	}

	return insertManyResult.InsertedIDs, nil
}

//////////////////////////////////////////////////////////////////////////////////
// Read
//////////////////////////////////////////////////////////////////////////////////

/**
 * Find many values from a query
 */
func (self *MongoStore) Find(ctx context.Context, database string, collection string, query string, fields []string, optionsStr string) ([]interface{}, error) {
	if self.client.Database(database) == nil {
		return nil, errors.New("No database found with name " + database)
	}

	if self.client.Database(database).Collection(collection) == nil {
		return nil, errors.New("No collection found with name " + collection)
	}

	collection_ := self.client.Database(database).Collection(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return nil, err
	}

	var opts *options.FindOptions
	if len(optionsStr) > 0 {
		opts = new(options.FindOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return nil, err
		}
	}

	cur, err := collection_.Find(ctx, q, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)

	for cur.Next(ctx) {
		entity := make(map[string]interface{})
		err := cur.Decode(&entity)
		if err != nil {
			return nil, err
		}
		// In that case I will return the whole entity
		if len(fields) == 0 {
			results = append(results, entity)
		} else {
			values := make([]interface{}, len(fields))
			for i := 0; i < len(fields); i++ {
				values[i] = entity[fields[i]]
			}
			results = append(results, values)
		}
	}

	// In case of error
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

/**
 * Find one result at time.
 */
func (self *MongoStore) FindOne(ctx context.Context, database string, collection string, query string, fields []string, optionsStr string) (interface{}, error) {

	if self.client.Database(database) == nil {
		return nil, errors.New("No database found with name " + database)
	}

	if self.client.Database(database).Collection(collection) == nil {
		return nil, errors.New("No collection found with name " + collection)
	}

	collection_ := self.client.Database(database).Collection(collection)
	entity := make(map[string]interface{})

	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return nil, err
	}

	var opts *options.FindOneOptions
	if len(optionsStr) > 0 {
		opts = new(options.FindOneOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return nil, err
		}
	}

	err = collection_.FindOne(ctx, q, opts).Decode(&entity)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if len(fields) == 0 {
		result = entity
	} else {
		values := make([]interface{}, len(fields))

		for i := 0; i < len(fields); i++ {
			values[i] = entity[fields[i]]
		}
		result = values
	}

	return result, nil
}

//////////////////////////////////////////////////////////////////////////////////
// Update
//////////////////////////////////////////////////////////////////////////////////

/**
 * Update one or more value that match the query.
 */
func (self *MongoStore) Update(ctx context.Context, database string, collection string, query string, value string, optionsStr string) error {
	if self.client.Database(database) == nil {
		return errors.New("No database found with name " + database)
	}

	if self.client.Database(database).Collection(collection) == nil {
		return errors.New("No collection found with name " + collection)
	}

	collection_ := self.client.Database(database).Collection(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return err
	}

	v := make(map[string]interface{})
	err = json.Unmarshal([]byte(value), &v)
	if err != nil {
		return err
	}

	var opts *options.UpdateOptions
	if len(optionsStr) > 0 {
		opts = new(options.UpdateOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return err
		}
	}

	_, err = collection_.UpdateMany(ctx, q, v, opts)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Update one document at time
 */
func (self *MongoStore) UpdateOne(ctx context.Context, database string, collection string, query string, value string, optionsStr string) error {
	if self.client.Database(database) == nil {
		return errors.New("No database found with name " + database)
	}

	if self.client.Database(database).Collection(collection) == nil {
		return errors.New("No collection found with name " + collection)
	}

	collection_ := self.client.Database(database).Collection(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return err
	}

	v := make(map[string]interface{})
	err = json.Unmarshal([]byte(value), &v)
	if err != nil {
		return err
	}

	var opts *options.UpdateOptions
	if len(optionsStr) > 0 {
		opts = new(options.UpdateOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return err
		}
	}

	_, err = collection_.UpdateOne(ctx, q, v, opts)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Replace a document by another.
 */
func (self *MongoStore) ReplaceOne(ctx context.Context, database string, collection string, query string, value string, optionsStr string) error {
	if self.client.Database(database) == nil {
		return errors.New("No database found with name " + database)
	}

	if self.client.Database(database).Collection(collection) == nil {
		return errors.New("No collection found with name " + collection)
	}

	collection_ := self.client.Database(database).Collection(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return err
	}

	v := make(map[string]interface{})
	err = json.Unmarshal([]byte(value), &v)
	if err != nil {
		return err
	}

	var opts *options.ReplaceOptions
	if len(optionsStr) > 0 {
		opts = new(options.ReplaceOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return err
		}
	}

	_, err = collection_.ReplaceOne(ctx, q, v, opts)
	if err != nil {
		return err
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////////
// Delete
//////////////////////////////////////////////////////////////////////////////////

/**
 * Remove one or more value depending of the query results.
 */
func (self *MongoStore) Delete(ctx context.Context, database string, collection string, query string, optionsStr string) error {
	if self.client.Database(database) == nil {
		return errors.New("No database found with name " + database)
	}

	if self.client.Database(database).Collection(collection) == nil {
		return errors.New("No collection found with name " + collection)
	}

	collection_ := self.client.Database(database).Collection(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return err
	}

	var opts *options.DeleteOptions
	if len(optionsStr) > 0 {
		opts = new(options.DeleteOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return err
		}
	}
	_, err = collection_.DeleteMany(ctx, q, opts)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Remove one document at time
 */
func (self *MongoStore) DeleteOne(ctx context.Context, database string, collection string, query string, optionsStr string) error {
	if self.client.Database(database) == nil {
		return errors.New("No database found with name " + database)
	}

	if self.client.Database(database).Collection(collection) == nil {
		return errors.New("No collection found with name " + collection)
	}

	collection_ := self.client.Database(database).Collection(collection)
	q := make(map[string]interface{})
	err := json.Unmarshal([]byte(query), &q)
	if err != nil {
		return err
	}

	var opts *options.DeleteOptions
	if len(optionsStr) > 0 {
		opts = new(options.DeleteOptions)
		err := json.Unmarshal([]byte(optionsStr), opts)
		if err != nil {
			return err
		}
	}

	_, err = collection_.DeleteOne(ctx, q, opts)
	if err != nil {
		return err
	}

	return nil
}
