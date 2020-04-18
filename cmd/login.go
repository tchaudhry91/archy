package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
			fmt.Printf("Could not login: %v", err)
			return
		}
		fmt.Printf("Succesfully Logged In\n Token:%s", token)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	registerCmd.PersistentFlags().StringVar(&loginU, "user", "", "Username")
	registerCmd.PersistentFlags().StringVar(&loginP, "password", "", "Password")
}
