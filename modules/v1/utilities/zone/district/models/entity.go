package models

import (
	regions "ethical-be/modules/v1/utilities/zone/region/models"
	"time"

	"gorm.io/gorm"
)

type Districts struct {
	ID                   uint
	Name                 string            `json:"name" gorm:"type:varchar(256)"`
	CreatedAt            time.Time         `json:"created_at" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt            time.Time         `json:"updated_at"`
	DeletedAt            gorm.DeletedAt    `json:"deleted_at"`
	IsVacant             bool              `json:"is_vacant" gorm:"type:boolean;default:true"`
	RegionsUnderDistrict []regions.Regions `gorm:"foreignKey:DistrictID;References:ID"`
}
