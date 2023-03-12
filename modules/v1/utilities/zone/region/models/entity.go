package models

import (
	area "ethical-be/modules/v1/utilities/zone/area/models"
	"time"

	"gorm.io/gorm"
)

type Regions struct {
	ID               uint           `gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"type:varchar(256)"`
	DistrictID       int            `json:"district_id"`
	CreatedAt        time.Time      `json:"created_at" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at"`
	IsVacant         bool           `json:"is_vacant" gorm:"type:boolean;default:true"`
	AreasUnderRegion []area.Areas   `gorm:"foreignKey:RegionID;References:ID"`
}
