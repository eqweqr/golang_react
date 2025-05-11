package models

import "fmt"

type Order struct {
	ID    int
	Name  string
	Age   int
	Email string
	Phone string
}

func (ord Order) String() string {
	return fmt.Sprintf("employee %s: %d\n", ord.Name, ord.Age)
}
