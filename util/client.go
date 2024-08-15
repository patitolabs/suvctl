package util

import (
	"github.com/patitolabs/gosuv2"
	"github.com/spf13/viper"
)

type Client struct {
	SuvConfig *gosuv2.SuvConfig
	SuvClient *gosuv2.SuvClient
}

func ReadConfig() *gosuv2.SuvConfig {
	return &gosuv2.SuvConfig{
		Host:       viper.GetString("host"),
		PhpSession: viper.GetString("session"),
		Detailed:   viper.GetBool("detailed"),
	}
}

func NewClient(config *gosuv2.SuvConfig) *Client {
	return &Client{
		SuvConfig: config,
		SuvClient: gosuv2.NewSuvClient(*config),
	}
}

func (c *Client) SetPhpSession(session string) {
	c.SuvConfig.PhpSession = session
	c.SuvClient.LoadPhpSession()
	viper.Set("session", &session)
	viper.Set("detailed", false)
	viper.WriteConfig()
}
