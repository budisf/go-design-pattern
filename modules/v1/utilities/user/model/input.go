package model

type UserRequestDTO struct {
	Name         *string `json:"name" binding:"required,min=4"`
	AuthServerId *uint   `json:"auth_server_id" binding:"required,gte=1"`
	Nip          *string `json:"nip" binding:"required,min=4"`
	Email        *string `json:"email" binding:"required,email"`
	//RoleId       *uint   `json:"role_id" binding:"omitempty,gte=1"`
}

type UserRequestUpdateDTO struct {
	Name  *string `json:"name" binding:"required,min=4"`
	Email *string `json:"email" binding:"required,email"`
}
