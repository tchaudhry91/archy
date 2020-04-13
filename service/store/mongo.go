package store

import (
	"fmt"
	"time"

	"github.com/tchaudhry91/zsh-archaeologist/history"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
)

const collection = "entries"
const database = "history"

// MongoStore is a struct to represent MongoDB for history entries
type MongoStore struct {
	client *mongo.Client
}

// NewMongoStore returns a new MongoDB backed store
func NewMongoStore(uri string) (ms *MongoStore, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return ms, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return ms, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return ms, err
	}

	return &MongoStore{client}, nil
}

// GetEntries is used to obtain entries based on filters
func (s *MongoStore) GetEntries(ctx context.Context, user string, filter bson.D) ([]history.Entry, error) {
	entries := []history.Entry{}
	coll := s.client.Database(database).Collection(collection)
	cur, err := coll.Find(ctx, AndMergeFilters(filter, SelectForUserFilter(user)))
	defer cur.Close(ctx)
	if err != nil {
		return entries, err
	}
	for cur.Next(ctx) {
		var result EntryDocument
		err := cur.Decode(&result)
		if err != nil {
			return entries, err
		}
		entries = append(entries, result.Entry)
	}
	if cur.Err() != nil {
		return entries, err
	}

	return entries, nil
}

// StoreEntries Stores the entries to the mongo store for the given user
func (s *MongoStore) StoreEntries(ctx context.Context, user string, entries []history.Entry) error {
	coll := s.client.Database(database).Collection(collection)

	models := []mongo.WriteModel{}
	for _, e := range entries {
		m := mongo.NewUpdateOneModel().SetFilter(EntryDocument{User: user, Entry: e}).SetUpdate(EntryDocument{e, user}).SetUpsert(true)
		models = append(models, m)
	}

	opts := options.BulkWrite().SetOrdered(false)

	res, err := coll.BulkWrite(ctx, models, opts)
	if err != nil {
		return err
	}
	fmt.Println(res.InsertedCount)
	fmt.Println(res.UpsertedCount)
	fmt.Println(res.ModifiedCount)
	return nil
}
