package cmd

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search a user by code and show its information",
	Run:   search,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().BoolP("professors", "t", false, "search professors (default is students)")
	searchCmd.Flags().StringP("code", "c", "", "code of the user to search")
	searchCmd.Flags().StringP("name", "n", "", "name of the user to search")
	searchCmd.Flags().StringP("lastname", "l", "", "lastname of the user to search")
	searchCmd.Flags().StringP("dni", "i", "", "DNI of the user to search")

	searchCmd.MarkFlagsRequiredTogether("name", "lastname")
	searchCmd.MarkFlagsMutuallyExclusive("code", "name", "dni")
	searchCmd.MarkFlagsMutuallyExclusive("professors", "dni")
}

func search(cmd *cobra.Command, args []string) {
	professors := cmd.Flag("professors").Value.String() == "true"
	code := cmd.Flag("code").Value.String()
	name := cmd.Flag("name").Value.String()
	lastname := cmd.Flag("lastname").Value.String()
	dni := cmd.Flag("dni").Value.String()

	if code == "" && name == "" && dni == "" {
		cmd.Println("You must provide a code, name or dni")
		cmd.Println()
		cmd.Usage()
		return
	}

	if professors {
		c.SearchProfessor(name, lastname)
	} else {
		c.SearchStudent(code, name, lastname, dni)
	}
}
