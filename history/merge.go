package history

import (
	"sort"
)

// MergeHistory combines to Entry slices and sorts them by timestamps
func MergeHistory(entrySlices ...[]Entry) []Entry {
	combined := []Entry{}
	for _, ee := range entrySlices {
		combined = append(combined, ee...)
	}
	sort.Slice(combined, func(i, j int) bool {
		if combined[i].Timestamp == combined[j].Timestamp {
			return combined[i].Command < combined[j].Command
		}
		return combined[i].Timestamp < combined[j].Timestamp
	})
	return deduplicate(combined)
}

// deduplicate removes duplicate entries
func deduplicate(sortedEntries []Entry) []Entry {
	// Entries are considered duplicate when the timestamp and command are the same
	dedup := []Entry{}
	for i := range sortedEntries {
		if len(dedup) > 0 && (sortedEntries[i].Timestamp == dedup[len(dedup)-1].Timestamp) && (sortedEntries[i].Command == dedup[len(dedup)-1].Command) {
			continue
		}
		dedup = append(dedup, sortedEntries[i])
	}
	return dedup
}
