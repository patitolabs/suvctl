package cmd

import "github.com/spf13/cobra"

var gradesCmd = &cobra.Command{
	Use:   "grades",
	Short: "List the grades of the current period",
	Run:   grades,
}

func init() {
	rootCmd.AddCommand(gradesCmd)
}

func grades(cmd *cobra.Command, args []string) {
	c.ListGrades()
}
