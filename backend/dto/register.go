package dto

type RegisterDTO struct {
	Token string
}

func NewRegisterDTO(token string) *RegisterDTO {
	return &RegisterDTO{
		Token: token,
	}
}
