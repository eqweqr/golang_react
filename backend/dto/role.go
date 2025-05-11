package dto

type RoleDTO struct {
	Role string
}

func NewRoleDTO(role string) *RoleDTO {
	return &RoleDTO{
		Role: role,
	}
}
