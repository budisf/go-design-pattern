package models

import (
	"time"

	"gorm.io/gorm"
)

type DistrictResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	IsVacant  bool           `json:"is_vacant" gorm:"type:boolean;default:true"`
}

type DistrictResponseEpochTime struct {
	ID        *uint   `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *uint   `json:"created_at"`
	UpdatedAt *uint   `json:"updated_at"`
	IsVacant  *bool   `json:"is_vacant"`
}
