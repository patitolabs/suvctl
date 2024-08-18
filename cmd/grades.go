package cmd

import "github.com/spf13/cobra"

var gradesCmd = &cobra.Command{
	Use:   "grades",
	Short: "List the grades of the current period",
	Run:   grades,
}

func init() {
	rootCmd.AddCommand(gradesCmd)

	gradesCmd.Flags().StringArrayP("courseid", "i", []string{}, "Filter by course ID")
}

func grades(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		c.ListGrades()
	} else {
		if courseId, err := cmd.Flags().GetStringArray("courseid"); err == nil {
			c.ListGradesByCourseId(courseId)
		}
	}
}
