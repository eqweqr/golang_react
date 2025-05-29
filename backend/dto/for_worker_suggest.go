package dto

type FullSuggestionDTO struct {
	OrderId     int
	ModelName   string
	Warranty    bool
	Comment     string
	OrderStatus string
	ClientName  string
	ClientPhone string
	ClientEmail string
	WorkType    string
	Proposed    bool
	WorkerName  interface{}
	WorkerPhone interface{}
	WorkerEmail interface{}
	Summary     interface{}
	Term        interface{}
	ConfTime    interface{}
	CreatedAt   interface{}
}
