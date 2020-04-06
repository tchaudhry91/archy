package history_test

import (
	"testing"

	"github.com/tchaudhry91/zsh-archaeologist/history"
)

func TestMergeEntries(t *testing.T) {
	entries1, err := getSampleEntries()
	if err != nil {
		t.Errorf("Failed to get sample entries:%v", err)
	}
	entries2, err := getSampleEntries2()
	if err != nil {
		t.Errorf("Failed to get sample entries:%v", err)
	}

	combined := history.MergeHistory(entries1, entries2, entries2, entries1)

	if len(combined) != (len(entries1) + len(entries2)) {
		t.Errorf("Incorrect Merge Length, Want:%d, Have:%d", (len(entries1) + len(entries2)), len(combined))
	}
	t.Logf("Merged into total %d entries", len(combined))
}
