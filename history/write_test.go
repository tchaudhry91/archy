package history_test

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/tchaudhry91/zsh-archaeologist/history"
)

func getSampleEntries() ([]history.Entry, error) {
	return history.ParseFile("samples/sample1.test_history", "test")
}

func getSampleEntries2() ([]history.Entry, error) {
	return history.ParseFile("samples/sample2.test_history", "test2")
}

func getFileHash(fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}

func TestSampleWrite(t *testing.T) {
	ee, err := getSampleEntries()
	if err != nil {
		t.Errorf("Failed to get sample entries:%v", err)
	}

	// Write Entries to a file
	fname := filepath.Join(os.TempDir(), ".zsh_history")
	err = history.WriteHistoryFile(ee, fname)
	if err != nil {
		t.Errorf("Failed to write entries to files:%v", err)
	}

	// Compare files
	hashOriginal, err := getFileHash(filepath.Join("samples", "sample1.test_history"))
	if err != nil {
		t.Errorf("Failed to Compute Hash:%v", err)
	}
	hashNew, err := getFileHash(fname)
	if err != nil {
		t.Errorf("Failed to Compute Hash:%v", err)
	}
	if hashOriginal != hashNew {
		t.Errorf("Hashes do not match. Want %s, have %s", hashOriginal, hashNew)
	}
	os.Remove(fname)

}
