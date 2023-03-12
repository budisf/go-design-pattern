package model

import "time"

type UserRoleResponseDTO struct {
	ID           *uint   `json:"user_id"`
	Name         *string `json:"name"`
	AuthServerId *uint   `json:"auth_server_id"`
	Nip          *string `json:"nip"`
	RoleId       *uint   `json:"role_id"`
	RoleName     *string `json:"role_name"`
	RoleLabel    *string `json:"role_label"`
	ParentId     *uint   `json:"parent_id"`
	NameRoleHead *string `json:"parent_role_name"`
	CreatedAt    *uint64 `json:"created_at"`
	UpdatedAt    *uint64 `json:"updated_at"`
}

type UserResponseEpochDTO struct {
	ID           *uint   `json:"user_id"`
	Name         *string `json:"name"`
	AuthServerId *uint   `json:"auth_server_id"`
	Nip          *string `json:"nip"`
	RoleId       *uint   `json:"role_id"`
	RoleName     *string `json:"role_name"`
	Email        *string `json:"email"`
	Label        *string `json:"label"`
	CreatedAt    *uint64 `json:"created_at"`
	UpdatedAt    *uint64 `json:"updated_at"`
}

type UserResponseDTO struct {
	ID           *uint      `json:"user_id"`
	Name         *string    `json:"name"`
	AuthServerId *uint      `json:"auth_server_id"`
	Nip          *string    `json:"nip"`
	RoleId       *uint      `json:"role_id"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
