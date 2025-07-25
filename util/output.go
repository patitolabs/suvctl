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
	OutputText  OutputFormat = "text"
	OutputTable OutputFormat = "table"
	OutputJSON  OutputFormat = "json"
	OutputRaw   OutputFormat = "raw"
)

// GetOutputFormat returns the current output format from viper config
func GetOutputFormat() OutputFormat {
	format := viper.GetString("output")
	switch format {
	case "text":
		return OutputText
	case "json":
		return OutputJSON
	case "raw":
		return OutputRaw
	default:
		return OutputTable
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
	case OutputRaw:
		outputGradesRaw(grades)
	case OutputTable:
		outputGradesTable(grades)
	default:
		outputGradesText(grades)
	}
}

// OutputStudents outputs students in the specified format
func OutputStudents(students []gosuv2.StudentBasicResponse) {
	format := GetOutputFormat()
	switch format {
	case OutputJSON:
		outputStudentsJSON(students)
	case OutputRaw:
		outputStudentsRaw(students)
	case OutputTable:
		outputStudentsTable(students)
	default:
		outputStudentsText(students)
	}
}

// OutputProfessors outputs professors in the specified format
func OutputProfessors(professors []gosuv2.ProfessorBasicResponse) {
	format := GetOutputFormat()
	switch format {
	case OutputJSON:
		outputProfessorsJSON(professors)
	case OutputRaw:
		outputProfessorsRaw(professors)
	case OutputTable:
		outputProfessorsTable(professors)
	default:
		outputProfessorsText(professors)
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

	// Analyze which columns have data
	columns := analyzeGradeColumns(grades)

	// Build the table
	printGradeTableHeader(columns)
	printGradeTableSeparator(columns)

	for _, grade := range grades {
		printGradeTableRow(grade, columns)
	}

	printGradeTableFooter(columns)
}

// Column represents a table column with its properties
type Column struct {
	Name    string
	Width   int
	Align   string // "left", "right", "center"
	HasData bool
}

func analyzeGradeColumns(grades []gosuv2.SuvCurrentCourseGrades) []Column {
	// Initialize all possible columns
	columns := []Column{
		{"Course", 9, "right", true},
		{"Course Name", 36, "left", true},
		{"Attempt", 9, "right", true},
		{"Unit 1", 9, "right", false},
		{"Unit 2", 9, "right", false},
		{"Unit 3", 9, "right", false},
		{"Unit 4", 9, "right", false},
		{"Unit 5", 9, "right", false},
		{"Unit 6", 9, "right", false},
		{"Subst", 9, "right", false},
		{"Failed", 9, "right", false},
		{"Average", 9, "right", false},
		{"Final Avg", 9, "right", false},
		{"Status", 10, "center", true},
	}

	// Check which columns have data
	for _, grade := range grades {
		if grade.Average1 != 0 {
			columns[3].HasData = true
		}
		if grade.Average2 != 0 {
			columns[4].HasData = true
		}
		if grade.Average3 != 0 {
			columns[5].HasData = true
		}
		if grade.Average4 != 0 {
			columns[6].HasData = true
		}
		if grade.Average5 != 0 {
			columns[7].HasData = true
		}
		if grade.Average6 != 0 {
			columns[8].HasData = true
		}
		if grade.Substitute != 0 {
			columns[9].HasData = true
		}
		if grade.Postponed != 0 {
			columns[10].HasData = true
		}
		if grade.Average != 0 {
			columns[11].HasData = true
		}
		if grade.FinalAverage != 0 {
			columns[12].HasData = true
		}
	}

	// Filter to only include columns with data
	var activeColumns []Column
	for _, col := range columns {
		if col.HasData {
			activeColumns = append(activeColumns, col)
		}
	}

	return activeColumns
}

func printGradeTableHeader(columns []Column) {
	// Top border
	fmt.Print("╭")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", col.Width))
		if i < len(columns)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("╮")

	// Header row
	fmt.Print("│")
	for _, col := range columns {
		padding := col.Width - len(col.Name)
		switch col.Align {
		case "center":
			leftPad := padding / 2
			rightPad := padding - leftPad
			fmt.Printf("%s%s%s│", strings.Repeat(" ", leftPad), col.Name, strings.Repeat(" ", rightPad))
		case "right":
			fmt.Printf("%s%s │", strings.Repeat(" ", padding-1), col.Name)
		default: // left
			fmt.Printf(" %-*s│", col.Width-1, col.Name)
		}
	}
	fmt.Println()
}

func printGradeTableSeparator(columns []Column) {
	fmt.Print("├")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", col.Width))
		if i < len(columns)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")
}

func printGradeTableFooter(columns []Column) {
	fmt.Print("╰")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", col.Width))
		if i < len(columns)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Println("╯")
}

func printGradeTableRow(grade gosuv2.SuvCurrentCourseGrades, columns []Column) {
	finalStatus := determineFinalStatus(grade)

	// Truncate course name if too long
	courseName := grade.CourseName
	maxNameLen := 34 // Account for padding
	if len(courseName) > maxNameLen {
		courseName = courseName[:maxNameLen-3] + "..."
	}

	fmt.Print("│")

	for _, col := range columns {
		var content string
		var useColor bool
		var colorCode string

		switch col.Name {
		case "Course":
			content = fmt.Sprintf("%d", grade.CourseID)
		case "Course Name":
			content = courseName
		case "Attempt":
			content = fmt.Sprintf("%d", grade.Attempt)
		case "Unit 1":
			content, useColor, colorCode = formatGradeValue(grade.Average1)
		case "Unit 2":
			content, useColor, colorCode = formatGradeValue(grade.Average2)
		case "Unit 3":
			content, useColor, colorCode = formatGradeValue(grade.Average3)
		case "Unit 4":
			content, useColor, colorCode = formatGradeValue(grade.Average4)
		case "Unit 5":
			content, useColor, colorCode = formatGradeValue(grade.Average5)
		case "Unit 6":
			content, useColor, colorCode = formatGradeValue(grade.Average6)
		case "Subst":
			content, useColor, colorCode = formatGradeValue(grade.Substitute)
		case "Failed":
			content, useColor, colorCode = formatGradeValue(grade.Postponed)
		case "Average":
			content, useColor, colorCode = formatGradeValue(grade.Average)
		case "Final Avg":
			content, useColor, colorCode = formatGradeValue(grade.FinalAverage)
		case "Status":
			content = finalStatus
			useColor = true
			colorCode = getStatusColor(finalStatus)
		}

		// Apply formatting based on column alignment
		padding := col.Width - len(content)
		if useColor {
			switch col.Align {
			case "center":
				leftPad := padding / 2
				rightPad := padding - leftPad
				fmt.Printf("%s%s%s%s\033[0m│", strings.Repeat(" ", leftPad), colorCode, content, strings.Repeat(" ", rightPad))
			case "right":
				fmt.Printf("%s%s%s\033[0m │", strings.Repeat(" ", padding-1), colorCode, content)
			default: // left
				fmt.Printf(" %s%s\033[0m%s│", colorCode, content, strings.Repeat(" ", padding-1))
			}
		} else {
			switch col.Align {
			case "center":
				leftPad := padding / 2
				rightPad := padding - leftPad
				fmt.Printf("%s%s%s│", strings.Repeat(" ", leftPad), content, strings.Repeat(" ", rightPad))
			case "right":
				fmt.Printf("%s%s │", strings.Repeat(" ", padding-1), content)
			default: // left
				fmt.Printf(" %-*s│", col.Width-1, content)
			}
		}
	}
	fmt.Println()

	// Show warning for disqualified students
	if grade.Disabled {
		fmt.Print("│")
		for i, col := range columns {
			if i == 1 { // Course Name column
				warning := "\033[31mWARNING: Student disqualified\033[0m"
				fmt.Printf(" %-*s│", col.Width-1, warning)
			} else {
				fmt.Printf("%s│", strings.Repeat(" ", col.Width))
			}
		}
		fmt.Println()
	}
}

func formatGradeValue(grade float32) (string, bool, string) {
	if grade == 0 {
		return "-", false, ""
	}

	content := fmt.Sprintf("%.2f", grade)
	useColor := true
	colorCode := getGradeColor(grade)

	return content, useColor, colorCode
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

	// Calculate column widths based on content
	columns := []Column{
		{"Student ID", 13, "left", true},
		{"Student Name", 36, "left", true},
		{"DNI", 13, "left", true},
	}

	// Adjust column widths based on actual content
	maxIDLen := len("Student ID")
	maxNameLen := len("Student Name")
	maxDNILen := len("DNI")

	for _, student := range students {
		if len(student.StudentID) > maxIDLen {
			maxIDLen = len(student.StudentID)
		}
		if len(student.StudentName) > maxNameLen {
			maxNameLen = len(student.StudentName)
		}
		if len(student.DNI) > maxDNILen {
			maxDNILen = len(student.DNI)
		}
	}

	// Apply minimum widths and padding
	columns[0].Width = max(maxIDLen+2, 13)
	columns[1].Width = max(maxNameLen+2, 36)
	columns[2].Width = max(maxDNILen+2, 13)

	fmt.Println("Students found:")
	printTableHeader(columns)
	printTableSeparator(columns)

	for _, student := range students {
		printStudentTableRow(student, columns)
	}

	printTableFooter(columns)
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

	// Calculate column widths based on content
	columns := []Column{
		{"Code", 13, "left", true},
		{"Professor Name", 36, "left", true},
		{"DNI", 13, "left", true},
		{"Worker ID", 13, "left", true},
	}

	// Adjust column widths based on actual content
	maxCodeLen := len("Code")
	maxNameLen := len("Professor Name")
	maxDNILen := len("DNI")
	maxWorkerLen := len("Worker ID")

	for _, professor := range professors {
		if len(professor.Code) > maxCodeLen {
			maxCodeLen = len(professor.Code)
		}
		if len(professor.ProfessorName) > maxNameLen {
			maxNameLen = len(professor.ProfessorName)
		}
		if len(professor.DNI) > maxDNILen {
			maxDNILen = len(professor.DNI)
		}
		if len(professor.WorkerID) > maxWorkerLen {
			maxWorkerLen = len(professor.WorkerID)
		}
	}

	// Apply minimum widths and padding
	columns[0].Width = max(maxCodeLen+2, 13)
	columns[1].Width = max(maxNameLen+2, 36)
	columns[2].Width = max(maxDNILen+2, 13)
	columns[3].Width = max(maxWorkerLen+2, 13)

	fmt.Println("Professors found:")
	printTableHeader(columns)
	printTableSeparator(columns)

	for _, professor := range professors {
		printProfessorTableRow(professor, columns)
	}

	printTableFooter(columns)
}

// Text output functions (existing behavior)
func outputGradesText(grades []gosuv2.SuvCurrentCourseGrades) {
	for _, grade := range grades {
		prettyPrintGradeCourse(grade)
		fmt.Println()
	}
}

func outputStudentsText(students []gosuv2.StudentBasicResponse) {
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

func outputProfessorsText(professors []gosuv2.ProfessorBasicResponse) {
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

// Raw JSON output functions (no indentation for piping)
func outputGradesRaw(grades []gosuv2.SuvCurrentCourseGrades) {
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

	output, err := json.Marshal(gradeData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func outputStudentsRaw(students []gosuv2.StudentBasicResponse) {
	var studentData []StudentData
	for _, student := range students {
		data := StudentData{
			StudentID:   student.StudentID,
			StudentName: student.StudentName,
			DNI:         student.DNI,
		}
		studentData = append(studentData, data)
	}

	output, err := json.Marshal(studentData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

func outputProfessorsRaw(professors []gosuv2.ProfessorBasicResponse) {
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

	output, err := json.Marshal(professorData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}

// Generic table helper functions
func printTableHeader(columns []Column) {
	// Top border
	fmt.Print("╭")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", col.Width))
		if i < len(columns)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("╮")

	// Header row
	fmt.Print("│")
	for _, col := range columns {
		padding := col.Width - len(col.Name)
		switch col.Align {
		case "center":
			leftPad := padding / 2
			rightPad := padding - leftPad
			fmt.Printf("%s%s%s│", strings.Repeat(" ", leftPad), col.Name, strings.Repeat(" ", rightPad))
		case "right":
			fmt.Printf("%s%s │", strings.Repeat(" ", padding-1), col.Name)
		default: // left
			fmt.Printf(" %-*s│", col.Width-1, col.Name)
		}
	}
	fmt.Println()
}

func printTableSeparator(columns []Column) {
	fmt.Print("├")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", col.Width))
		if i < len(columns)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")
}

func printTableFooter(columns []Column) {
	fmt.Print("╰")
	for i, col := range columns {
		fmt.Print(strings.Repeat("─", col.Width))
		if i < len(columns)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Println("╯")
}

func printStudentTableRow(student gosuv2.StudentBasicResponse, columns []Column) {
	fmt.Print("│")

	for _, col := range columns {
		var content string

		switch col.Name {
		case "Student ID":
			content = student.StudentID
		case "Student Name":
			content = student.StudentName
		case "DNI":
			content = student.DNI
		}

		// Apply formatting based on column alignment
		padding := col.Width - len(content)
		switch col.Align {
		case "center":
			leftPad := padding / 2
			rightPad := padding - leftPad
			fmt.Printf("%s%s%s│", strings.Repeat(" ", leftPad), content, strings.Repeat(" ", rightPad))
		case "right":
			fmt.Printf("%s%s │", strings.Repeat(" ", padding-1), content)
		default: // left
			fmt.Printf(" %-*s│", col.Width-1, content)
		}
	}
	fmt.Println()
}

func printProfessorTableRow(professor gosuv2.ProfessorBasicResponse, columns []Column) {
	fmt.Print("│")

	for _, col := range columns {
		var content string

		switch col.Name {
		case "Code":
			content = professor.Code
		case "Professor Name":
			content = professor.ProfessorName
		case "DNI":
			content = professor.DNI
		case "Worker ID":
			content = professor.WorkerID
		}

		// Apply formatting based on column alignment
		padding := col.Width - len(content)
		switch col.Align {
		case "center":
			leftPad := padding / 2
			rightPad := padding - leftPad
			fmt.Printf("%s%s%s│", strings.Repeat(" ", leftPad), content, strings.Repeat(" ", rightPad))
		case "right":
			fmt.Printf("%s%s │", strings.Repeat(" ", padding-1), content)
		default: // left
			fmt.Printf(" %-*s│", col.Width-1, content)
		}
	}
	fmt.Println()
}
