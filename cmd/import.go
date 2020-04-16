package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var maxEntries int64

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import allows you to add external entries to your base zsh history",
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Flags
	importCmd.PersistentFlags().Int64Var(&maxEntries, "maxEntries", 10000, "The maximum number of newest entries to write. This may be overridden by your zsh config")
	viper.BindPFlag("maxEntries", importCmd.Flags().Lookup("maxEntries"))
}
