package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tchaudhry91/archy/service/client"
)

var machine string
var command string

var tsEnd int64
var tsBegin int64
var limit int64

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search across your command history on the remote service",
	Run: func(cmd *cobra.Command, args []string) {
		token = detectToken()
		c, err := client.NewHistoryClient(remoteAddr, token, 100)
		if err != nil {
			panic(err)
		}

		req := client.GetEntriesRequest{
			Machine: machine,
			Command: command,
			Start:   uint64(tsBegin),
			End:     uint64(tsEnd),
			Limit:   limit,
		}

		entries, err := c.GetEntries(req)
		if err != nil {
			fmt.Printf("Failed to get entries from remote: %v\n", err)
		}
		for _, e := range entries {
			fmt.Printf("%d\t%s\t%s\n", e.Timestamp, e.Machine, e.Command)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().Int64Var(&tsEnd, "tsEnd", time.Now().Unix(), "Unix timestamp till when to search. Default: Now")
	searchCmd.PersistentFlags().Int64Var(&tsBegin, "timestampBegin", time.Now().AddDate(-1, 0, 0).Unix(), "Unix timestamp from when to search. Default: 1 year ago")
	searchCmd.PersistentFlags().Int64Var(&limit, "limit", 100, "Limit the number of results")
	searchCmd.PersistentFlags().StringVar(&machine, "machine", "", "Filter for a specific machine")
	searchCmd.PersistentFlags().StringVar(&command, "command", "", "Filter for a specific command regex")
}
