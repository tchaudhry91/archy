package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/tchaudhry91/zsh-archaeologist/history"
)

// ImportLocalConfig is a struct to hold config values for local file imports.
type ImportLocalConfig struct {
	// RootCfg is the root configuration.
	ImportCfg *ImportConfig
}

// Exec performs the ImportLocal action
func (cfg *ImportLocalConfig) Exec(ctx context.Context, args []string) error {
	entries := []history.Entry{}
	for _, f := range args {
		ee, err := history.ParseFile(f, "localhost")
		if err != nil {
			return fmt.Errorf("Failed to import from file: %s because %v", f, err)
		}
		// TO-DO: Make efficient later
		entries = history.MergeHistory(entries, ee)
	}
	err := history.WriteHistoryFile(entries, cfg.ImportCfg.RootCfg.HistoryFile)
	if err != nil {
		return err
	}
	return nil
}

// NewImportLocalCommand returns a command that can be used to import a file locally
func NewImportLocalCommand(importCfg *ImportConfig) *ffcli.Command {
	cfg := &ImportLocalConfig{
		ImportCfg: importCfg,
	}
	fs := flag.NewFlagSet("arch import local", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "local",
		ShortUsage: "archy import local [flags] files...",
		ShortHelp:  "Import and merge history from other files into your local base.",
		FlagSet:    fs,
		Exec:       cfg.Exec,
	}
}

// ImportConfig contains values for import related options
type ImportConfig struct {
	// RootConfig is the root config reference
	RootCfg *RootConfig
	// MaxLength denotes the number of entries to write in the file.
	MaxLength int64
}

// Exec function for this command.
func (c *ImportConfig) Exec(context.Context, []string) error {
	return flag.ErrHelp
}

// NewImportCommand returns a command that can be used for import entries
func NewImportCommand(rootCfg *RootConfig) *ffcli.Command {
	cfg := &ImportConfig{
		RootCfg: rootCfg,
	}
	fs := flag.NewFlagSet("arch import", flag.ExitOnError)
	fs.Int64Var(&cfg.MaxLength, "l", 100000, "The maximum number of entries to write to the file. This may be overwritten by your ZSH config")
	importLocalCmd := NewImportLocalCommand(cfg)
	sub := []*ffcli.Command{importLocalCmd}
	return &ffcli.Command{
		Name:        "import",
		ShortUsage:  "archy import [flags] sub-cmd",
		ShortHelp:   "Import allows you to merge entries into your base file",
		FlagSet:     fs,
		Exec:        cfg.Exec,
		Subcommands: sub,
	}
}
