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
	gradesCmd.Flags().StringArrayP("course", "n", []string{}, "Filter by course name")
}

func grades(cmd *cobra.Command, args []string) {
	courseIds, err := cmd.Flags().GetStringArray("courseid")
	cobra.CheckErr(err)

	courseNames, err := cmd.Flags().GetStringArray("course")
	cobra.CheckErr(err)

	if len(courseIds) > 0 {
		c.ListGradesByCourseId(courseIds)
	}

	if len(courseNames) > 0 {
		c.ListGradesByCourseName(courseNames)
	}

	if len(courseIds) == 0 && len(courseNames) == 0 {
		c.ListGrades()
	}
}
