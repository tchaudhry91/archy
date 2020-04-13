package store

import "github.com/tchaudhry91/zsh-archaeologist/history"

// EntryDocument is the storage translation for an entry
type EntryDocument struct {
	history.Entry
	User string
}
