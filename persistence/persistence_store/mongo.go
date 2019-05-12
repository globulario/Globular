package persistence_store

import (
	"context"
	//	"fmt"
	//	"log"

	//	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"

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

//////////////////////////////////////////////////////////////////////////////////
// Insert
//////////////////////////////////////////////////////////////////////////////////
/**
 * Insert one value in the store.
 */
func (self *MongoStore) InsertOne(ctx context.Context, database string, collection string, entity interface{}) (interface{}, error) {
	// Get the collection object.
	collection_ := self.client.Database(database).Collection(collection)

	result, err := collection_.InsertOne(ctx, entity)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

/**
 * Insert many results at time.
 */
func (self *MongoStore) InsertMany(ctx context.Context, database string, collection string, entities []interface{}) ([]interface{}, error) {
	// Get the collection object.
	collection_ := self.client.Database(database).Collection(collection)

	// return self.client.Ping(ctx, nil)
	insertManyResult, err := collection_.InsertMany(ctx, entities)
	if err != nil {
		return nil, err
	}

	return insertManyResult.InsertedIDs, nil
}
