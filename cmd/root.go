package cmd

import (
	"fmt"
	"path"

	"github.com/adrg/xdg"
	"github.com/patitolabs/gosuv2"
	"github.com/patitolabs/suvctl/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "suvctl",
		Short: "A command-line tool for SUV2 at National University of Trujillo",
	}

	c       *util.Client
	config  *gosuv2.SuvConfig
	session string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.config/suvctl/suvctl.yaml)")
	rootCmd.PersistentFlags().StringP("host", "U", "", "SUV host FQDN (default is suv2.unitru.edu.pe)")
	rootCmd.PersistentFlags().StringP("session", "S", "", "session for SUV operations")
	rootCmd.PersistentFlags().BoolP("detailed", "d", false, "show detailed information")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "show version information")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("session", rootCmd.PersistentFlags().Lookup("session"))
	viper.BindPFlag("detailed", rootCmd.PersistentFlags().Lookup("detailed"))
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		configPath := path.Join(xdg.ConfigHome, "suvctl")

		viper.AddConfigPath(configPath)

		viper.SetConfigName("suvctl")
		viper.SetConfigType("yaml")
		viper.SetConfigType("yml")

		viper.SetDefault("host", "suv2.unitru.edu.pe")
	}

	err := viper.ReadInConfig()

	if viper.GetBool("detailed") {
		fmt.Println("suvctl is running in detailed mode")
		fmt.Println()
	}

	if err == nil {
		if viper.GetBool("detailed") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
			fmt.Println()
		}
	}

	if viper.GetBool("version") {
		versionCmd.Run(versionCmd, []string{})
		fmt.Println()
	}

	if viper.GetString("session") != "" {
		session = viper.GetString("session")
	} else {
		session = rootCmd.Flags().Lookup("session").Value.String()
	}

	config = util.ReadConfig()
	c = util.NewClient(config)

	if session != "" {
		c.SetPhpSession(session)

		if viper.GetBool("detailed") {
			fmt.Println("Using session:", session)
			fmt.Println()
		}
	}
}
