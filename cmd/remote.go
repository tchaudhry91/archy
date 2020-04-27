package cmd

import (
	"fmt"

	"github.com/tchaudhry91/archy/service/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tchaudhry91/archy/history"
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "exports your command history to a remote server",
	Run: func(cmd *cobra.Command, args []string) {
		token = detectToken()
		lastTS := detectLast()
		entries, err := history.ParseFile(baseHistoryFile, hostname)
		if err != nil {
			panic(err)
		}

		// Reduce entry slice for new entries only
		if lastTS > 0 {
			entries = history.SliceEntries(lastTS, entries)
		}
		if len(entries) == 0 {
			fmt.Println("No new entries to send")
			return
		}

		c, err := client.NewHistoryClient(remoteAddr, token, 180)
		if err != nil {
			panic(err)
		}

		req := client.PutEntriesRequest{
			Entries: entries,
		}

		updated, err := c.PutEntries(req)
		if err != nil {
			fmt.Printf("Could not update entries: %v", err)
			return
		}
		fmt.Printf("Succesfully Updated %d Entries\n", updated)
		lastTS = entries[len(entries)-1].Timestamp
		viper.Set("last", lastTS)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Failed to write back config:%v", err)
		}
	},
}

func detectToken() string {
	if token == "" {
		if viper.InConfig("token") {
			token = viper.GetString("token")
		}
	}
	return token
}

func detectLast() uint64 {
	var last uint64
	if viper.InConfig("last") {
		last = viper.GetUint64("last")
	}
	return last
}

func init() {
	exportCmd.AddCommand(remoteCmd)
}
