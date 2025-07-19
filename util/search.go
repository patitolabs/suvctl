package util

import (
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
	OutputStudents(*students)
}

func (c *Client) SearchProfessor(name, lastname string) {
	professors, err := c.SuvClient.SearchProfessor(name, lastname)
	cobra.CheckErr(err)
	OutputProfessors(*professors)
}
