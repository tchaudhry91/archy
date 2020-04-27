package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tchaudhry91/archy/history"
)

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "supply filepaths you want to import entries from other local zsh_history files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filesToImport := args

		base, err := history.ParseFile(baseHistoryFile, hostname)
		if err != nil {
			fmt.Printf("Unable to parse base history file: %v", err)
		}

		for _, f := range filesToImport {
			new, err := history.ParseFile(f, hostname)
			if err != nil {
				fmt.Printf("Unable to parse file to import:%v", err)
			}
			base = history.MergeHistory(base, new)
		}
		index := int64(0)
		if int64(len(base)) > maxEntries {
			index = int64(len(base)) - maxEntries
		}
		base = base[index:]
		history.WriteHistoryFile(base, baseHistoryFile)
	},
}

func init() {
	importCmd.AddCommand(localCmd)
}
