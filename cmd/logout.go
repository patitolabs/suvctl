package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Destroy a session for the user and remove it from the config file",
	Run:   logout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	logoutCmd.Flags().BoolP("force", "F", false, "force the logout without asking for confirmation")

	viper.BindPFlag("force", logoutCmd.Flags().Lookup("force"))
}

func logout(cmd *cobra.Command, args []string) {
	if session == "" {
		cmd.Println("No session to logout")
		fmt.Println()
		cmd.Usage()
		os.Exit(1)
		return
	}

	c.Logout()
}
