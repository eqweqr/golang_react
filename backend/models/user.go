package models

import "fmt"

type User struct {
	ID    int
	Name  string
	Email string
	Phone string
	Role  string
	Pass  string
}

func (rec User) String() string {
	return fmt.Sprintf("employee %s:\n", rec.Name)
}
