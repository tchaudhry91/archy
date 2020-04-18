package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var baseHistoryFile string
var token string
var hostname string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "archy",
	Short: "A zsh history manipulator",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Attempt to get hostname
	host, err := os.Hostname()
	if err != nil {
		host = ""
	}

	// Attempt to get homedir
	home, err := homedir.Dir()
	if err != nil {
		home = ""
	}

	// Persistent Flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.archy.yaml)")
	rootCmd.PersistentFlags().StringVar(&baseHistoryFile, "baseHistoryFile", filepath.Join(home, ".zsh_history"), "base history file to operate on")
	rootCmd.PersistentFlags().StringVar(&hostname, "hostname", host, "Override the hostname value for entries from this machine")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "the token to communicate with the remote service")
	rootCmd.PersistentFlags().StringVar(&remoteAddr, "remote", "https://archy.tux-sudo.com", "Address of the remote service to contact")

	viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))
	viper.BindPFlag("baseHistoryFile", rootCmd.Flags().Lookup("baseHistoryFile"))
	viper.BindPFlag("hostname", rootCmd.Flags().Lookup("hostname"))
	viper.BindPFlag("remote", rootCmd.Flags().Lookup("remote"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".archy" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".archy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
