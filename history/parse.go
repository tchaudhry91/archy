package history

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseFile returns Entry objects from a given zsh_history file
func ParseFile(location string, machine string) ([]Entry, error) {
	var errRet error
	ee := []Entry{}
	f, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(zshEntrySplitFunc)

	for scanner.Scan() {
		text := scanner.Text()
		entry, err := ParseLine(text)
		if err != nil {
			errRet = err
			continue
		}
		if !verifyEntrySanity(entry) {
			continue
		}
		entry.Machine = machine
		ee = append(ee, entry)
	}
	return ee, errRet
}

// SliceEntries reduces the slice of entries to only include entries AFTER the given timestamp
func SliceEntries(lastTS uint64, entries []Entry) []Entry {
	reduced := []Entry{}
	for _, e := range entries {
		if e.Timestamp > lastTS {
			reduced = append(reduced, e)
		}
	}
	return reduced
}

// ParseLine parses a single line of text to an Entry Object
func ParseLine(line string) (Entry, error) {
	split := strings.Split(line, ";")
	if len(split) < 2 {
		return Entry{}, fmt.Errorf("Invalid Line Found, not enough ';' found:%s", line)
	}
	timestampElapsed := split[0]
	command := split[1]

	// re-split timestamp_elapsed to remove elapsed
	split = strings.Split(timestampElapsed, ":")
	if len(split) < 2 {
		return Entry{}, fmt.Errorf("Invalid Line Found, not enough ':' found:%s", line)
	}
	timestampStr := strings.TrimSpace(split[1])
	timestamp, err := strconv.ParseUint(timestampStr, 10, 64)
	if err != nil {
		return Entry{}, fmt.Errorf("Invalid Line Found, unparsable timestamp found:%s", timestampStr)
	}
	return Entry{
		Timestamp: timestamp,
		Command:   command,
	}, nil
}

func verifyEntrySanity(e Entry) bool {
	if e.Command == "" || e.Timestamp == 0 {
		return false
	}
	return true
}

// zshEntrySplitFunc splits files on every zsh entry
func zshEntrySplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), "\n: "); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	return
}
