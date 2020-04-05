package history_test

import (
	"path/filepath"
	"testing"

	"github.com/tchaudhry91/zsh-archaeologist/history"
)

func TestParseLine(t *testing.T) {
	type TestCase struct {
		Name      string
		Line      string
		Command   string
		Timestamp uint64
		Valid     bool
	}

	cases := []TestCase{
		TestCase{"Good", ": 1575991023:0;ls", "ls", 1575991023, true},
		TestCase{"BadTimestamp", ": i1575991023:0;ls", "ls", 1575991023, false},
		TestCase{"BadLine", ":", "ls", 1575991023, false},
		TestCase{"BadLine2", ";", "ls", 1575991023, false},
	}

	for _, c := range cases {
		entry, err := history.ParseLine(c.Line)
		if err != nil && c.Valid {
			t.Errorf("Unexpected error on case: %s, %v", c.Name, err)
			continue
		}
		if err == nil && !c.Valid {
			t.Errorf("Failed to error on case: %s", c.Name)
			continue
		}
		if !c.Valid {
			continue
		}
		if c.Timestamp != entry.Timestamp {
			t.Errorf("Incorrect Timestamp Parsed, Have:%d Want:%d", entry.Timestamp, c.Timestamp)
		}
		if c.Command != entry.Command {
			t.Errorf("Incorrect Command Parsed, Have:%s, Want:%s", entry.Command, c.Command)
		}
	}
}

func TestParseFile(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("samples", "*.test_history"))
	if err != nil {
		t.Errorf("Failed to gather sample files: %v", err)
	}
	for _, f := range files {
		ee, err := history.ParseFile(f)
		if err != nil {
			t.Errorf("Failed to parse sample file:%s because %v", f, err)
		}
		t.Logf("Found %d Entries in file:%s", len(ee), f)
	}
}
