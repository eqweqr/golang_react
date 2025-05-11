package models

import "fmt"

type Record struct {
	ID   int
	Name string
	Age  int
}

func (rec Record) String() string {
	return fmt.Sprintf("employee %s: %d\n", rec.Name, rec.Age)
}
