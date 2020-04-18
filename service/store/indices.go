package store

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// EnsureHistoryIndices makes sure the required indices exist on the database
func (s *MongoStore) EnsureHistoryIndices(ctx context.Context, recreate bool) error {
	err := s.client.Ping(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "unable to communicate with database")
	}

	col := s.client.Database(database).Collection(collection)
	iv := col.Indexes()

	// Unique Name-ID for Machines
	found, err := s.checkIndex(ctx, &iv, "unique-timestamp-command")
	if err != nil {
		return errors.Wrap(err, "unable to check for existence of index")
	}
	if !found || recreate {
		if recreate && found {
			_, err := iv.DropOne(ctx, "unique-timestamp-command", nil)
			if err != nil {
				return err
			}
		}
		// Create Unique Timestamp Command
		unique := *options.Index().SetUnique(true).SetName("unique-timestamp-command")
		model := mongo.IndexModel{
			Keys: bson.M{
				"user":      1,
				"timestamp": -1,
				"command":   -2,
			},
			Options: &unique,
		}

		_, err = iv.CreateOne(ctx, model, nil)
		if err != nil {
			return errors.Wrap(err, "failed to create unique-timestamp-command index")
		}
	}
	return nil
}

// EnsureUserIndices makes sure usernames are unique
func (s *MongoStore) EnsureUserIndices(ctx context.Context, recreate bool) error {
	col := s.client.Database(database).Collection(collectionUsers)
	iv := col.Indexes()

	// Unique user
	found, err := s.checkIndex(ctx, &iv, "unique-username")
	if err != nil {
		return errors.Wrap(err, "unable to check for existence of index")
	}
	if !found || recreate {
		if recreate && found {
			_, err := iv.DropOne(ctx, "unique-username", nil)
			if err != nil {
				return err
			}
		}
		// Create Unique Timestamp Command
		unique := *options.Index().SetUnique(true).SetName("unique-username")
		model := mongo.IndexModel{
			Keys: bson.M{
				"user": 1,
			},
			Options: &unique,
		}
		_, err = iv.CreateOne(ctx, model, nil)
		if err != nil {
			return errors.Wrap(err, "unique-username")
		}
	}
	return nil
}

// checkIndex checks if a particular index exists in an IndexView
func (s *MongoStore) checkIndex(ctx context.Context, iv *mongo.IndexView, name string) (bool, error) {
	var found bool
	cur, err := iv.List(ctx)
	if err != nil {
		return found, errors.Wrap(err, "unable to list indices")
	}
	for cur.Next(ctx) {
		d := &bson.Raw{}
		err := cur.Decode(d)
		if err != nil {
			return found, errors.Wrap(err, "unable to decode bson index doc")
		}
		if d.Lookup("name").StringValue() == name {
			return true, nil
		}
	}
	err = cur.Close(ctx)
	if err != nil {
		return found, errors.Wrap(err, "unable to close index")
	}
	return found, nil
}
