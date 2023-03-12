package model

type RoleRequestDTO struct {
	Name     *string `json:"role_name" binding:"required,min=2"`
	Label    *string `json:"role_label" binding:"required,min=6"`
	ParentId *uint   `json:"parent_id" binding:"omitempty,gte=1"`
}

type RoleRequestUpdateDTO struct {
	Name  *string `json:"role_name" binding:"required,min=2"`
	Label *string `json:"role_label" binding:"required,min=6"`
}

type RoleRequestUpdateHeadRoleDTO struct {
	IdHeadRole *uint `json:"id_head_role" binding:"omitempty,gte=1"`
}
