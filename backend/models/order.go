package models

import "fmt"

type Order struct {
	ID          int
	ModelName   string
	Warranty    bool
	Comment     string
	ClientId    int
	WorkerId    int
	WorkType    string
	OrderStatus string
	ConfTime    string
	Summary     string
	CreatedAt   string
}

func (ord Order) String() string {
	return fmt.Sprintf("Order %s:\n", ord.ModelName)
}
