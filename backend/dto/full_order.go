package dto

type FullOrderDTO struct {
	OrderId     int
	ModelName   string
	Warranty    bool
	Comment     string
	OrderStatus string
	ClientName  string
	ClientPhone string
	ClientEmail string
	WorkType    string
	WorkerName  interface{}
	WorkerPhone interface{}
	WorkerEmail interface{}
	Summary     interface{}
	ConfTime    interface{}
	CreatedAt   interface{}
}
