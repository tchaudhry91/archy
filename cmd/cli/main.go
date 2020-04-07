package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

// RootConfig is a struct to hold config values for the entire CLI
type RootConfig struct {
	// HistoryFile is the location of the base zsh history to operate on.
	HistoryFile string
}

func main() {
	cfg := &RootConfig{}
	rootFlagSet := flag.NewFlagSet("archy", flag.ExitOnError)
	rootFlagSet.StringVar(&cfg.HistoryFile, "f", "~/.zsh_history", "Base History File")

	sub := []*ffcli.Command{
		NewImportCommand(cfg),
	}

	rootCmd := &ffcli.Command{
		Name:        "archy",
		ShortUsage:  "archy [flags] sub-cmd",
		Subcommands: sub,
		FlagSet:     rootFlagSet,
		Exec:        cfg.Exec,
	}

	if err := rootCmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

// Exec function for this command.
func (c *RootConfig) Exec(context.Context, []string) error {
	// The root command has no meaning, so if it gets executed,
	// display the usage text to the user instead.
	return flag.ErrHelp
}
