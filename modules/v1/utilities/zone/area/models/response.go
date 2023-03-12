package models

import (
	"time"

	"gorm.io/gorm"
)

type AreasResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	RegionID  int            `json:"region_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	IsVacant  bool           `json:"is_vacant" gorm:"type:boolean;default:true"`
}

type AreasResponseEpochTime struct {
	ID        *uint   `json:"id"`
	Name      *string `json:"name"`
	RegionID  *int    `json:"region_id"`
	CreatedAt *uint   `json:"created_at"`
	UpdatedAt *uint   `json:"updated_at"`
	IsVacant  *bool   `json:"is_vacant"`
}
