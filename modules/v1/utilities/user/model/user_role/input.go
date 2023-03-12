package model

type UserRoleRequestDTO struct {
	RoleId *uint `json:"role_id" binding:"required,gte=1"`
}
