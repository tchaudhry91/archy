package cmd

import (
	"fmt"

	"github.com/tchaudhry91/zsh-archaeologist/service/client"

	"github.com/spf13/cobra"
	"github.com/tchaudhry91/zsh-archaeologist/history"
)

var remoteAddr string

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "exports your command history to a remote server",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := history.ParseFile(baseHistoryFile, hostname)
		if err != nil {
			panic(err)
		}

		c, err := client.NewHistoryClient(remoteAddr, "willbeatoken", 30)
		if err != nil {
			panic(err)
		}

		req := client.PutEntriesRequest{
			Entries: entries,
		}

		updated, err := c.PutEntries(req)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Succesfully Updated %d Entries\nm", updated)
	},
}

func init() {
	exportCmd.AddCommand(remoteCmd)

	exportCmd.PersistentFlags().StringVar(&remoteAddr, "remote", "https://archy.tux-sudo.com", "Address of the remote service to contact")
}
