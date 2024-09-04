package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

func (c *Client) ListGradesByCourseId(courseId []string) {
	suvGradesResponse, err := c.SuvClient.GetSuvGradesResponse()
	cobra.CheckErr(err)

	courseIdMap := make(map[string]struct{})
	for _, id := range courseId {
		courseIdMap[id] = struct{}{}
	}

	found := false
	for _, grade := range suvGradesResponse.Courses {
		if _, exists := courseIdMap[grade.IdCurso]; exists {
			prettyPrintGradeCourse(grade)
			fmt.Println()
			found = true
		}
	}

	if !found {
		fmt.Println("No courses found.")
		os.Exit(1)
	}
}

func (c *Client) ListGradesByCourseName(courseName []string) {
	suvGradesResponse, err := c.SuvClient.GetSuvGradesResponse()
	cobra.CheckErr(err)

	found := false
	for _, grade := range suvGradesResponse.Courses {
		for _, name := range courseName {
			if strings.Contains(grade.Curso, name) {
				prettyPrintGradeCourse(grade)
				fmt.Println()
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("No courses found.")
		os.Exit(1)
	}
}

func prettyPrintGradeCourse(grade gosuv2.SuvCurrentCourseGrades) {
	fmt.Println("Course ID:", grade.IdCurso)
	fmt.Println("Course:", grade.Curso)
	fmt.Println("Time:", grade.Vez)
	printAverage(grade.Promedio1, "Average of Unit 1:")
	printAverage(grade.Promedio2, "Average of Unit 2:")
	printAverage(grade.Promedio3, "Average of Unit 3:")
	printAverage(grade.Promedio4, "Average of Unit 4:")
	printAverage(grade.Promedio5, "Average of Unit 5:")
	printAverage(grade.Promedio6, "Average of Unit 6:")
	printAverage(grade.Sustitutorio, "Substitute exam:")
	printAverage(grade.Promedio, "Course Average:")
	printAverage(grade.Aplazado, "Failed:")
	printAverage(grade.PromedioFinal, "Course Final Average:")

	if grade.Inhabilitado != "0" {
		fmt.Println("\033[31mWarning: the student was disqualified in this course\033[0m")
	}

	printFinalStatus(grade)
}

func printAverage(gradeStr string, message string) {
	if gradeStr != "" {
		parseAndPrintGrade(gradeStr, message)
	}
}

func printFinalStatus(grade gosuv2.SuvCurrentCourseGrades) {
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
