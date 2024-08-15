package util

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (c *Client) Login(usercode, password string) {
	session, err := c.SuvClient.Login(usercode, password)
	cobra.CheckErr(err)

	c.SetPhpSession(*session)

	if viper.GetBool("detailed") {
		fmt.Println("Login successful. Session ID:", *session)
	} else {
		fmt.Println("Login successful")
	}
}

func (c *Client) Logout() {
	err := c.SuvClient.Logout()

	if viper.GetBool("force") {
		viper.Set("session", "")
		viper.Set("detailed", false)
		viper.WriteConfig()
	}

	cobra.CheckErr(err)

	c.SetPhpSession("")
	fmt.Println("Logout successful")
}
