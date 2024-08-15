package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version   = "devbuild"
	commit    = "none"
	buildDate = "unknown"
	platform  = "unknown"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version information of suvctl",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(`
                      _   _ 
                     | | | |
  ___ _   ___   _____| |_| |
 / __| | | \ \ / / __| __| |
 \__ \ |_| |\ V / (__| |_| |
 |___/\__,_| \_/ \___|\__|_|
			`)
			cmd.Println("suvctl", version)
			cmd.Println("commit", commit)
			cmd.Println("built on", buildDate)
			cmd.Println("running on", platform)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
