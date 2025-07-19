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

	OutputGrades(suvGradesResponse.Courses)
}

func (c *Client) ListGradesByCourseId(courseId []string) {
	suvGradesResponse, err := c.SuvClient.GetSuvGradesResponse()
	cobra.CheckErr(err)

	courseIdMap := make(map[int]struct{})
	for _, id := range courseId {
		courseIdInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Invalid course ID: %s\n", id)
			continue
		}
		courseIdMap[courseIdInt] = struct{}{}
	}

	found := false
	var foundGrades []gosuv2.SuvCurrentCourseGrades
	for _, grade := range suvGradesResponse.Courses {
		if _, exists := courseIdMap[grade.CourseID]; exists {
			foundGrades = append(foundGrades, grade)
			found = true
		}
	}

	if !found {
		fmt.Println("No courses found.")
		os.Exit(1)
	}

	OutputGrades(foundGrades)
}

func (c *Client) ListGradesByCourseName(courseName []string) {
	suvGradesResponse, err := c.SuvClient.GetSuvGradesResponse()
	cobra.CheckErr(err)

	found := false
	var foundGrades []gosuv2.SuvCurrentCourseGrades
	for _, grade := range suvGradesResponse.Courses {
		for _, name := range courseName {
			if strings.Contains(grade.CourseName, name) {
				foundGrades = append(foundGrades, grade)
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("No courses found.")
		os.Exit(1)
	}

	OutputGrades(foundGrades)
}

func prettyPrintGradeCourse(grade gosuv2.SuvCurrentCourseGrades) {
	fmt.Println("Course ID:", grade.CourseID)
	fmt.Println("Course:", grade.CourseName)
	fmt.Println("Time:", grade.Attempt)
	printAverage(grade.Average1, "Average of Unit 1:")
	printAverage(grade.Average2, "Average of Unit 2:")
	printAverage(grade.Average3, "Average of Unit 3:")
	printAverage(grade.Average4, "Average of Unit 4:")
	printAverage(grade.Average5, "Average of Unit 5:")
	printAverage(grade.Average6, "Average of Unit 6:")
	printAverage(grade.Substitute, "Substitute exam:")
	printAverage(grade.Average, "Course Average:")
	printAverage(grade.Postponed, "Failed:")
	printAverage(grade.FinalAverage, "Course Final Average:")

	if grade.Disabled {
		fmt.Println("\033[31mWarning: the student was disqualified in this course\033[0m")
	}

	printFinalStatus(grade)
}

func printAverage(grade float32, message string) {
	if grade != 0 {
		printGrade(grade, message)
	}
}

func printGrade(grade float32, message string) {
	// If grade < 13.5 print the message in the default color, and the number in red
	// Else, print the message in the default color, and the number in light blue
	if grade < 13.5 {
		fmt.Printf("%s \033[31m%.2f\033[0m\n", message, grade)
	} else {
		fmt.Printf("%s \033[94m%.2f\033[0m\n", message, grade)
	}
}

func printFinalStatus(grade gosuv2.SuvCurrentCourseGrades) {
	if grade.FinalStatus == gosuv2.PassedStatus {
		// Print the final status in green
		fmt.Println("Final status: \033[32mPASSED\033[0m")
	} else {
		if grade.Average1 != 0 && grade.Average2 != 0 && grade.Average3 != 0 {
			if grade.Average >= 14 || grade.FinalAverage >= 14 {
				// Print the final status in green, student passed
				fmt.Println("Final status: \033[32mPASSED\033[0m")
			} else {
				// Print the final status in red, student failed
				fmt.Println("Final status: \033[31mFAILED\033[0m")
			}
		} else {
			// Print the final status in yellow, semester isn't over yet
			fmt.Println("Final status: \033[33mPENDING\033[0m")
		}
	}
}
