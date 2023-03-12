package models

import (
	gt "ethical-be/modules/v1/utilities/zone/gt/models"
	"time"

	"gorm.io/gorm"
)

type Areas struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(256)"`
	RegionID    int
	CreatedAt   time.Time             `json:"created_at" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt   time.Time             `json:"updated_at"`
	DeletedAt   gorm.DeletedAt        `json:"deleted_at"`
	IsVacant    bool                  `json:"is_vacant" gorm:"type:boolean;default:true"`
	GtUnderArea []gt.GroupTerritories `gorm:"foreignKey:AreaID;References:ID"`
}
