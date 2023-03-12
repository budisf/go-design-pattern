package models

import "time"

type CustomerRequest struct {
	ID           uint
	CustomerName string     `json:"customer_name" gorm:"type:varchar(256)" binding:"required"`
	CustomerCode string     `json:"customer_code" gorm:"type:varchar(256)" binding:"required"`
	Specialist   *string    `json:"specialist" gorm:"type:varchar(256)"`
	UpdatedAt    *time.Time `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CustomerOutletRequest struct {
	ID         uint
	CustomerId uint `json:"customer_id" binding:"required"`
	OutletID   uint `json:"outlet_id" binding:"required"`
	Customer   Customer
	UpdatedAt  *time.Time `json:"updated_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
