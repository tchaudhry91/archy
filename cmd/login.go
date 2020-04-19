package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tchaudhry91/zsh-archaeologist/service/client"
)

var loginU string
var loginP string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to the remote service",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewHistoryClient(remoteAddr, "", 100)
		if err != nil {
			panic(err)
		}

		req := client.LoginRequest{
			User:     loginU,
			Password: loginP,
		}

		token, err = c.Login(req)
		if err != nil {
			fmt.Printf("Could not login: %v\n", err)
			return
		}
		fmt.Printf("Succesfully Logged In\n Token:%s\n", token)
		if err = viper.WriteConfig(); err != nil {
			fmt.Println("Could not write config back to file:", err)
		}
		fmt.Printf("Your token has been updated in your config.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().StringVar(&loginU, "user", "", "Username")
	loginCmd.PersistentFlags().StringVar(&loginP, "password", "", "Password")
}
