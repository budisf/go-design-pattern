package model

import "time"

type RoleResponseDTO struct {
	RoleId     *uint        `json:"role_id"`
	Name       *string      `json:"role_name"`
	Label      *string      `json:"role_label"`
	ParentId   *uint        `json:"parent_id"`
	CreatedAt  *time.Time   `json:"created_at"`
	UpdatedAt  *time.Time   `json:"updated_at"`
	ParentRole *ParentRoles `gorm:"foreignKey:ParentId" json:"parent_role"`
}

type RoleResponse struct {
	RoleId     *uint                `json:"role_id"`
	Name       *string              `json:"role_name"`
	Label      *string              `json:"role_label"`
	CreatedAt  *uint64              `json:"created_at"`
	UpdatedAt  *uint64              `json:"updated_at"`
	ParentRole *ParentRolesResponse `gorm:"foreignKey:ParentId" json:"parent_role"`
}

type ParentRolesResponse struct {
	RoleID    *uint   `json:"role_id"`
	Name      *string `json:"role_name"`
	Label     *string `json:"role_label"`
	CreatedAt *uint64 `json:"created_at"`
	UpdatedAt *uint64 `json:"updated_at"`
}
