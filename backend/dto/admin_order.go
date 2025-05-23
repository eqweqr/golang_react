package dto

type AdminOrderDTO struct {
	Id       int
	Username string
	Status   string
}

func NewAdminOrderDTO(id int, username string, status string) *AdminOrderDTO {
	return &AdminOrderDTO{
		Id:       id,
		Username: username,
		Status:   status,
	}
}
