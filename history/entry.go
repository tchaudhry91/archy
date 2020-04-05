package history

// Entry is struct that holds values for a particular history line
type Entry struct {
	Timestamp uint64
	Machine   string
	Command   string
}

// NewEntry consturcts a new Entry
func NewEntry(timestamp uint64, machine string, command string) Entry {
	return Entry{
		Timestamp: timestamp,
		Machine:   machine,
		Command:   command,
	}
}
