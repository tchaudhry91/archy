package history

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// WriteHistoryFile write the entries back to the history file (default ~/.zsh_history) and backs up the last one
func WriteHistoryFile(entries []Entry, historyFile string) error {
	if historyFile == "" {
		historyFile = "~/.zsh_history"
	}
	// Backup file if exists
	err := BackupHistoryFileIfExists(historyFile)
	if err != nil {
		return err
	}

	// Write File
	f, err := os.Create(historyFile)
	if err != nil {
		return err
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	for i, e := range entries {
		line := ConvertEntryToLine(e)
		// Trim the newline if this is the last line
		if i == (len(entries) - 1) {
			line = strings.TrimRight(line, "\n")
		}
		_, err = writer.WriteString(line)
		if err != nil {
			// Move the old history back
			RestoreOldHistory(historyFile)
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

// ConvertEntryToLine is a reverse parse operation that converts an Entry to a line for the zsh_history file
func ConvertEntryToLine(entry Entry) string {
	return fmt.Sprintf(": %d:0;%s\n", entry.Timestamp, entry.Command)
}

// BackupHistoryFileIfExists creates a backup of the given history file in place with a .back suffix
func BackupHistoryFileIfExists(historyFile string) error {
	// Backup history file if it already exists
	existing := true
	data, err := ioutil.ReadFile(historyFile)
	if os.IsNotExist(err) {
		// Old File Doesn't Exist
		existing = false
	}
	if err != nil && existing {
		return fmt.Errorf("Failed to open old history file:%v", err)
	}
	if existing {
		err := ioutil.WriteFile(historyFile+".back", data, 0600)
		if err != nil {
			return fmt.Errorf("Failed to backup old history file:%v", err)
		}
	}
	return nil
}

// RestoreOldHistory moves the backup back to the main historyFile
func RestoreOldHistory(historyFile string) error {
	return os.Rename(historyFile+".back", historyFile)
}
