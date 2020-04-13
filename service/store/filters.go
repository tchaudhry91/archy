package store

import "go.mongodb.org/mongo-driver/bson"

// SelectAllFilter returns all documents
func SelectAllFilter() bson.D {
	return bson.D{}
}

// SelectForUserFilter returns all documents for a single user
func SelectForUserFilter(user string) bson.D {
	return bson.D{
		{Key: "user", Value: user},
	}
}

// SelectSinceTimestampFilter returns all documents since a particular timestamp
func SelectSinceTimestampFilter(ts uint64) bson.D {
	return bson.D{
		{
			Key: "entry.timestamp",
			Value: map[string]uint64{
				"$gt": ts,
			},
		},
	}
}

// AndMergeFilters combines multiple filters with an and operation
func AndMergeFilters(filters ...bson.D) bson.D {
	return bson.D{
		{
			Key:   "$and",
			Value: filters,
		},
	}
}
