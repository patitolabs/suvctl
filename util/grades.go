package util

import (
	"fmt"
	"strconv"

	"github.com/patitolabs/gosuv2"
	"github.com/spf13/cobra"
)

func (c *Client) ListGrades() {
	suvGradesResponse, err := c.SuvClient.GetSuvGradesResponse()
	cobra.CheckErr(err)

	for _, grade := range suvGradesResponse.Courses {
		prettyPrintGradeCourse(grade)
		fmt.Println()
	}
}

func prettyPrintGradeCourse(grade gosuv2.SuvCurrentCourseGrades) {
	fmt.Println("Course ID:", grade.IdCurso)
	fmt.Println("Course:", grade.Curso)
	fmt.Println("Time:", grade.Vez)
	if grade.Promedio1 != "" {
		parseAndPrintGrade(grade.Promedio1, "Average of Unit 1:")
	}
	if grade.Promedio2 != "" {
		parseAndPrintGrade(grade.Promedio2, "Average of Unit 2:")
	}
	if grade.Promedio3 != "" {
		parseAndPrintGrade(grade.Promedio3, "Average of Unit 3:")
	}
	if grade.Promedio4 != "" {
		parseAndPrintGrade(grade.Promedio4, "Average of Unit 4:")
	}
	if grade.Promedio5 != "" {
		parseAndPrintGrade(grade.Promedio5, "Average of Unit 5:")
	}
	if grade.Promedio6 != "" {
		parseAndPrintGrade(grade.Promedio6, "Average of Unit 6:")
	}
	if grade.Sustitutorio != "" {
		parseAndPrintGrade(grade.Sustitutorio, "Substitute exam:")
	}
	if grade.Promedio != "" {
		parseAndPrintGrade(grade.Promedio, "Course Average:")
	}
	if grade.Aplazado != "" {
		parseAndPrintGrade(grade.Aplazado, "Failed:")
	}
	if grade.PromedioFinal != "" {
		parseAndPrintGrade(grade.PromedioFinal, "Course Final Average:")
	}
	if grade.Inhabilitado != "0" {
		fmt.Println("\033[31mWarning: the student was disqualified in this course\033[0m")
	}
	if grade.EstadoFinal == "1" {
		// Print the final status in green
		fmt.Println("Final status: \033[32mPASSED\033[0m")
	} else {
		if grade.Promedio1 != "" && grade.Promedio2 != "" && grade.Promedio3 != "" {
			// Print the final status in red
			fmt.Println("Final status: \033[31mFAILED\033[0m")
		} else {
			// Print the final status in yellow, semester isn't over yet
			fmt.Println("Final status: \033[33mPENDING\033[0m")
		}
	}
}

func parseAndPrintGrade(gradeStr string, message string) {
	grade, err := strconv.ParseFloat(gradeStr, 32)
	if err != nil {
		fmt.Println("Error parsing grade:", err)
		return
	}

	// If grade < 13.5 print the message in the default color, and the number in red
	// Else, print the message in the default color, and the number in light blue
	if grade < 13.5 {
		fmt.Printf("%s \033[31m%.2f\033[0m\n", message, grade)
	} else {
		fmt.Printf("%s \033[94m%.2f\033[0m\n", message, grade)
	}
}
