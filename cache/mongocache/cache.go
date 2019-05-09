/*
Package mongocache implements the presto.Cache interface using a mongodb database for persistent storage.
*/
package mongocache

import (
	"context"

	"github.com/benpate/derp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// MongoCache represents a single database connection
type MongoCache struct {
	client     *mongo.Client
	database   string
	collection string
}

// Record represents the key/value pair that is written into the mongodb cache.
type Record struct {
	ID    primitive.ObjectID `bson:"_id"`
	Key   string             `bson:"key"`
	Value string             `bson:"value"`
}

// New returns a fully initialized mongodb cache
func New(client *mongo.Client, database string, collection string) *MongoCache {

	return &MongoCache{
		client:     client,
		database:   database,
		collection: collection,
	}
}

// Get retrieves a single value from the cache.  If the value does not exist
// in the cache, then "" is returned
func (cache *MongoCache) Get(key string) string {

	record := Record{}
	ctx := context.TODO()
	collection := cache.client.Database(cache.database).Collection(cache.collection)
	criteria := bson.M{"key": key}

	if err := collection.FindOne(ctx, criteria).Decode(&record); err != nil {
		return ""
	}

	return record.Value
}

// Set adds/updates a value in the cache.
func (cache *MongoCache) Set(key string, value string) *derp.Error {

	ctx := context.TODO()
	collection := cache.client.Database(cache.database).Collection(cache.collection)
	criteria := bson.M{"key": key}
	update := bson.M{"value": value}

	T := true
	options := options.UpdateOptions{Upsert: &T}

	if _, err := collection.UpdateOne(ctx, criteria, update, &options); err != nil {
		return derp.New(derp.CodeInternalError, "mongocache.Set", "Error setting cache value", err.Error(), key, value)
	}

	return nil
}
