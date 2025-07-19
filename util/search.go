package util

import (
	"fmt"

	"github.com/patitolabs/gosuv2"
	"github.com/spf13/cobra"
)

func (c *Client) SearchStudent(code, name, lastname, dni string) {
	var (
		err      error
		students *[]gosuv2.StudentBasicResponse
	)

	if code != "" {
		students, err = c.SuvClient.SearchStudentByCode(code)
	}
	if name != "" && lastname != "" {
		students, err = c.SuvClient.SearchStudentByName(name, lastname)
	}
	if dni != "" {
		students, err = c.SuvClient.SearchStudentByDni(dni)
	}

	cobra.CheckErr(err)
	c.printSearchStudentResponse(*students)
}

func (c *Client) SearchProfessor(name, lastname string) {
	professors, err := c.SuvClient.SearchProfessor(name, lastname)
	cobra.CheckErr(err)
	c.parseSearchProfessorResponse(*professors)
}

func (c *Client) printSearchStudentResponse(students []gosuv2.StudentBasicResponse) {
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

func (c *Client) parseSearchProfessorResponse(professors []gosuv2.ProfessorBasicResponse) {
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
