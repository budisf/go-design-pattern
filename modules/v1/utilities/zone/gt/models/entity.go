package models

import (
	"time"

	"gorm.io/gorm"
)

type GroupTerritories struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"type:varchar(256)"`
	AreaID    int
	CreatedAt time.Time      `json:"created_at" gorm:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	IsVacant  bool           `json:"is_vacant" gorm:"type:boolean;default:true"`
}
