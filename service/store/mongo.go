package store

import (
	"errors"
	"time"

	"github.com/tchaudhry91/archy/history"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
)

const collection = "entries"
const collectionUsers = "users"
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return ms, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return ms, err
	}

	st := &MongoStore{client}
	err = st.EnsureUserIndices(ctx, true)
	if err != nil {
		return st, err
	}

	return st, nil
}

// GetEntries is used to obtain entries based on filters
func (s *MongoStore) GetEntries(ctx context.Context, user string, filter bson.D, limit int64) ([]history.Entry, error) {
	entries := []history.Entry{}
	coll := s.client.Database(database).Collection(collection)
	if coll == nil {
		return nil, errors.New("Could not get collection")
	}

	cur, err := coll.Find(ctx, AndMergeFilters(filter, SelectForUserFilter(user)), options.Find().SetLimit(limit).SetSort(bson.D{{Key: "entry.timestamp", Value: -1}}))
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
func (s *MongoStore) StoreEntries(ctx context.Context, user string, entries []history.Entry) (changed int64, err error) {
	coll := s.client.Database(database).Collection(collection)
	if coll == nil {
		return 0, errors.New("Could not get collection")
	}

	models := []mongo.WriteModel{}
	for _, e := range entries {
		m := mongo.NewUpdateOneModel().SetFilter(EntryDocument{User: user, Entry: e}).SetUpdate(EntryDocument{e, user}).SetUpsert(true)
		models = append(models, m)
	}

	opts := options.BulkWrite().SetOrdered(false)

	res, err := coll.BulkWrite(ctx, models, opts)
	if err != nil {
		return 0, errors.New("Failed to Bulk write entries")
	}
	return (res.InsertedCount + res.UpsertedCount + res.ModifiedCount), nil
}

// PutUser stores the user in the databse
func (s *MongoStore) PutUser(ctx context.Context, u *User) error {
	coll := s.client.Database(database).Collection(collectionUsers)
	if coll == nil {
		return errors.New("Could not get collection")
	}
	_, err := coll.InsertOne(ctx, *u)
	return err
}

// GetUser retrieves a user from the database
func (s *MongoStore) GetUser(ctx context.Context, user string) (*User, error) {
	coll := s.client.Database(database).Collection(collectionUsers)
	if coll == nil {
		return nil, errors.New("Could not get collection")
	}
	res := coll.FindOne(ctx, SelectForUserFilter(user))
	if res.Err() != nil {
		return nil, res.Err()
	}
	u := User{}
	err := res.Decode(&u)
	return &u, err
}
