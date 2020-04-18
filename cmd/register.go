package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tchaudhry91/zsh-archaeologist/service/client"
)

var user, password string

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register to the remote service",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewHistoryClient(remoteAddr, "", 100)
		if err != nil {
			panic(err)
		}

		req := client.RegisterRequest{
			User:     user,
			Password: password,
		}

		err = c.Register(req)
		if err != nil {
			fmt.Printf("Could not register: %v", err)
			return
		}
		fmt.Printf("Succesfully Registered! Please login to get the token")
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.PersistentFlags().StringVar(&user, "user", "", "Username to register")
	registerCmd.PersistentFlags().StringVar(&password, "password", "", "Password to register")
}
