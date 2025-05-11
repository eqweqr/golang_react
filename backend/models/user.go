package models

import "fmt"

type User struct {
	ID    int
	Name  string
	Age   int
	Email string
	Phone string
	Role  string
	Pass  string
}

func (rec User) String() string {
	return fmt.Sprintf("employee %s: %d\n", rec.Name, rec.Age)
}
