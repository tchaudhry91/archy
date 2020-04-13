package store_test

import (
	"context"
	"testing"

	"github.com/tchaudhry91/zsh-archaeologist/history"
	"github.com/tchaudhry91/zsh-archaeologist/service/store"
)

const uri = "mongodb://mongoadmin:password@localhost:61686"

func TestMongoOperations(t *testing.T) {
	entries, err := history.ParseFile("../../history/samples/sample1.test_history", "localhost")
	if err != nil {
		t.Errorf("Cannot get entries:%v", err)
		t.FailNow()
	}

	s, err := store.NewMongoStore(uri)
	if err != nil {
		t.Errorf("Cannot create connection:%v", err)
		t.FailNow()
	}

	err = s.StoreEntries(context.Background(), "tchaudhryTest", entries)
	if err != nil {
		t.Errorf("First Pass Creating entries failed: %v", err)
	}
	ee, err := s.GetEntries(context.Background(), "tchaudhryTest", store.SelectSinceTimestampFilter(1575997968))
	if err != nil {
		t.Errorf("First Pass Read Failed:%v", err)
	}
	if len(ee) != len(entries) {
		t.Errorf("Incorrect Number of Entries found, want:%d, have:%d", len(entries), len(ee))
	}
}
