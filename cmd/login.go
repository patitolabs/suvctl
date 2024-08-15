package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Create a session for the user and store it for further use",
	Run:   login,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("usercode", "u", "", "user for SUV operations")
	loginCmd.Flags().StringP("password", "p", "", "password for SUV operations")

	loginCmd.MarkFlagsRequiredTogether("usercode", "password")

	viper.BindPFlag("usercode", loginCmd.Flags().Lookup("usercode"))
	viper.BindPFlag("password", loginCmd.Flags().Lookup("password"))
}

func login(cmd *cobra.Command, args []string) {
	usercode := viper.GetString("usercode")
	password := viper.GetString("password")

	if usercode == "" || password == "" {
		usercode = cmd.Flag("usercode").Value.String()
		password = cmd.Flag("password").Value.String()
	}

	if usercode == "" || password == "" {
		cmd.Println("You must provide a user code and password")
		cmd.Println()
		cmd.Usage()
		os.Exit(1)
		return
	}

	c.Login(usercode, password)
}
