package relations

import (
	models "ethical-be/modules/v1/utilities/role/model"
)

type RolesRelation struct {
	models.Roles
	ParentRole models.ParentRoles `gorm:"foreignKey:ParentId" json:"parent_role"`
}

func (RolesRelation) TableName() string {
	return "roles"
}
