package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/patitolabs/gosuv2"
	"github.com/spf13/viper"
)

// OutputFormat represents the different output formats available
type OutputFormat string

const (
	OutputDefault OutputFormat = "default"
	OutputTable   OutputFormat = "table"
	OutputJSON    OutputFormat = "json"
)

// GetOutputFormat returns the current output format from viper config
func GetOutputFormat() OutputFormat {
	format := viper.GetString("output")
	switch format {
	case "table":
		return OutputTable
	case "json":
		return OutputJSON
	default:
		return OutputDefault
	}
}

// GradeData represents structured grade data for formatting
type GradeData struct {
	CourseID     int     `json:"course_id"`
	CourseName   string  `json:"course_name"`
	Attempt      int     `json:"attempt"`
	Average1     float32 `json:"average_1,omitempty"`
	Average2     float32 `json:"average_2,omitempty"`
	Average3     float32 `json:"average_3,omitempty"`
	Average4     float32 `json:"average_4,omitempty"`
	Average5     float32 `json:"average_5,omitempty"`
	Average6     float32 `json:"average_6,omitempty"`
	Substitute   float32 `json:"substitute,omitempty"`
	Average      float32 `json:"average,omitempty"`
	Postponed    float32 `json:"postponed,omitempty"`
	FinalAverage float32 `json:"final_average,omitempty"`
	Disabled     bool    `json:"disabled"`
	FinalStatus  string  `json:"final_status"`
}

// StudentData represents structured student data for formatting
type StudentData struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	DNI         string `json:"dni"`
}

// ProfessorData represents structured professor data for formatting
type ProfessorData struct {
	Code          string `json:"code"`
	ProfessorName string `json:"professor_name"`
	DNI           string `json:"dni"`
	WorkerID      string `json:"worker_id"`
}

// OutputGrades outputs grades in the specified format
func OutputGrades(grades []gosuv2.SuvCurrentCourseGrades) {
	format := GetOutputFormat()
	switch format {
	case OutputJSON:
		outputGradesJSON(grades)
	case OutputTable:
		outputGradesTable(grades)
	default:
		outputGradesDefault(grades)
	}
}

// OutputStudents outputs students in the specified format
func OutputStudents(students []gosuv2.StudentBasicResponse) {
	format := GetOutputFormat()
	switch format {
	case OutputJSON:
		outputStudentsJSON(students)
	case OutputTable:
		outputStudentsTable(students)
	default:
		outputStudentsDefault(students)
	}
}

// OutputProfessors outputs professors in the specified format
func OutputProfessors(professors []gosuv2.ProfessorBasicResponse) {
	format := GetOutputFormat()
	switch format {
	case OutputJSON:
		outputProfessorsJSON(professors)
	case OutputTable:
		outputProfessorsTable(professors)
	default:
		outputProfessorsDefault(professors)
	}
}

func outputGradesJSON(grades []gosuv2.SuvCurrentCourseGrades) {
	var gradeData []GradeData
	for _, grade := range grades {
		finalStatus := determineFinalStatus(grade)
		data := GradeData{
			CourseID:     grade.CourseID,
			CourseName:   grade.CourseName,
			Attempt:      grade.Attempt,
			Average1:     grade.Average1,
			Average2:     grade.Average2,
			Average3:     grade.Average3,
			Average4:     grade.Average4,
			Average5:     grade.Average5,
			Average6:     grade.Average6,
			Substitute:   grade.Substitute,
			Average:      grade.Average,
			Postponed:    grade.Postponed,
			FinalAverage: grade.FinalAverage,
			Disabled:     grade.Disabled,
			FinalStatus:  finalStatus,
		}
		gradeData = append(gradeData, data)
	}

	output, err := json.MarshalIndent(gradeData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func outputGradesTable(grades []gosuv2.SuvCurrentCourseGrades) {
	if len(grades) == 0 {
		fmt.Println("No courses found.")
		return
	}

	// ASCII table header
	fmt.Println("╭─────────┬────────────────────────────────────┬─────────┬─────────┬─────────┬─────────┬─────────┬─────────┬──────────┬─────────┬─────────┬──────────┬─────────────┬──────────╮")
	fmt.Println("│ Course  │ Course Name                        │ Attempt │  Unit 1 │  Unit 2 │  Unit 3 │  Unit 4 │  Unit 5 │   Unit 6 │   Subst │  Failed │  Average │ Final Avg   │  Status  │")
	fmt.Println("├─────────┼────────────────────────────────────┼─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼──────────┼─────────┼─────────┼──────────┼─────────────┼──────────┤")

	for _, grade := range grades {
		finalStatus := determineFinalStatus(grade)

		// Truncate course name if too long
		courseName := grade.CourseName
		if len(courseName) > 34 {
			courseName = courseName[:31] + "..."
		}

		fmt.Printf("│ %7d │ %-34s │ %7d │", grade.CourseID, courseName, grade.Attempt)

		// Print averages with proper spacing
		printTableGrade(grade.Average1)
		printTableGrade(grade.Average2)
		printTableGrade(grade.Average3)
		printTableGrade(grade.Average4)
		printTableGrade(grade.Average5)
		printTableGrade(grade.Average6)
		printTableGrade(grade.Substitute)
		printTableGrade(grade.Postponed)
		printTableGrade(grade.Average)
		printTableGrade(grade.FinalAverage)

		// Status with color
		statusColor := getStatusColor(finalStatus)
		fmt.Printf(" %s%-8s\033[0m │\n", statusColor, finalStatus)

		if grade.Disabled {
			fmt.Println("│         │ \033[31mWARNING: Student disqualified\033[0m       │         │         │         │         │         │         │          │         │         │          │             │          │")
		}
	}

	fmt.Println("╰─────────┴────────────────────────────────────┴─────────┴─────────┴─────────┴─────────┴─────────┴─────────┴──────────┴─────────┴─────────┴──────────┴─────────────┴──────────╯")
}

func outputStudentsJSON(students []gosuv2.StudentBasicResponse) {
	var studentData []StudentData
	for _, student := range students {
		data := StudentData{
			StudentID:   student.StudentID,
			StudentName: student.StudentName,
			DNI:         student.DNI,
		}
		studentData = append(studentData, data)
	}

	output, err := json.MarshalIndent(studentData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func outputStudentsTable(students []gosuv2.StudentBasicResponse) {
	if len(students) == 0 {
		fmt.Println("No students found")
		return
	}

	fmt.Println("Students found:")
	fmt.Println("╭─────────────┬────────────────────────────────────┬─────────────╮")
	fmt.Println("│ Student ID  │ Student Name                       │     DNI     │")
	fmt.Println("├─────────────┼────────────────────────────────────┼─────────────┤")

	for _, student := range students {
		studentName := student.StudentName
		if len(studentName) > 34 {
			studentName = studentName[:31] + "..."
		}
		fmt.Printf("│ %-11s │ %-34s │ %-11s │\n", student.StudentID, studentName, student.DNI)
	}

	fmt.Println("╰─────────────┴────────────────────────────────────┴─────────────╯")
}

func outputProfessorsJSON(professors []gosuv2.ProfessorBasicResponse) {
	var professorData []ProfessorData
	for _, professor := range professors {
		data := ProfessorData{
			Code:          professor.Code,
			ProfessorName: professor.ProfessorName,
			DNI:           professor.DNI,
			WorkerID:      professor.WorkerID,
		}
		professorData = append(professorData, data)
	}

	output, err := json.MarshalIndent(professorData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func outputProfessorsTable(professors []gosuv2.ProfessorBasicResponse) {
	if len(professors) == 0 {
		fmt.Println("No professors found")
		return
	}

	fmt.Println("Professors found:")
	fmt.Println("╭─────────────┬────────────────────────────────────┬─────────────┬─────────────╮")
	fmt.Println("│    Code     │ Professor Name                     │     DNI     │  Worker ID  │")
	fmt.Println("├─────────────┼────────────────────────────────────┼─────────────┼─────────────┤")

	for _, professor := range professors {
		professorName := professor.ProfessorName
		if len(professorName) > 34 {
			professorName = professorName[:31] + "..."
		}
		fmt.Printf("│ %-11s │ %-34s │ %-11s │ %-11s │\n", professor.Code, professorName, professor.DNI, professor.WorkerID)
	}

	fmt.Println("╰─────────────┴────────────────────────────────────┴─────────────┴─────────────╯")
}

// Default output functions (existing behavior)
func outputGradesDefault(grades []gosuv2.SuvCurrentCourseGrades) {
	for _, grade := range grades {
		prettyPrintGradeCourse(grade)
		fmt.Println()
	}
}

func outputStudentsDefault(students []gosuv2.StudentBasicResponse) {
	if len(students) == 0 {
		fmt.Println("No students found")
	} else {
		fmt.Println("Students found:")
		for _, student := range students {
			fmt.Println()
			fmt.Println("Code:", student.StudentID)
			fmt.Println("Name:", student.StudentName)
			fmt.Println("DNI:", student.DNI)
		}
	}
}

func outputProfessorsDefault(professors []gosuv2.ProfessorBasicResponse) {
	if len(professors) == 0 {
		fmt.Println("No professors found")
	} else {
		fmt.Println("Professors found:")
		for _, professor := range professors {
			fmt.Println()
			fmt.Println("Code:", professor.Code)
			fmt.Println("Name:", professor.ProfessorName)
			fmt.Println("DNI:", professor.DNI)
			fmt.Println("Worker ID:", professor.WorkerID)
		}
	}
}

// Helper functions
func printTableGrade(grade float32) {
	if grade != 0 {
		color := getGradeColor(grade)
		fmt.Printf(" %s%7.2f\033[0m │", color, grade)
	} else {
		fmt.Printf("    -    │")
	}
}

func getGradeColor(grade float32) string {
	if grade < 13.5 {
		return "\033[31m" // Red
	}
	return "\033[94m" // Light blue
}

func getStatusColor(status string) string {
	switch strings.ToUpper(status) {
	case "PASSED":
		return "\033[32m" // Green
	case "FAILED":
		return "\033[31m" // Red
	case "PENDING":
		return "\033[33m" // Yellow
	default:
		return "\033[0m" // Default
	}
}

func determineFinalStatus(grade gosuv2.SuvCurrentCourseGrades) string {
	if grade.FinalStatus == gosuv2.PassedStatus {
		return "PASSED"
	} else {
		if grade.Average1 != 0 && grade.Average2 != 0 && grade.Average3 != 0 {
			if grade.Average >= 14 || grade.FinalAverage >= 14 {
				return "PASSED"
			} else {
				return "FAILED"
			}
		} else {
			return "PENDING"
		}
	}
}
