package models

import "time"

type Customer struct {
	ID           uint
	CustomerCode string     `json:"customer_code" gorm:"type:varchar(256)"`
	CustomerName string     `json:"customer_name" gorm:"type:varchar(256)"`
	Specialist   *string    `json:"specialist" gorm:"type:varchar(256)"`
	Position     *string    `json:"position" gorm:"type:varchar(256)"`
	UpdatedAt    *time.Time `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at"`
	IsDeleted    bool       `json:"is_deleted"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type CustomerOutlet struct {
	ID         uint
	CustomerId uint `json:"customer_id"`
	OutletID   uint `json:"outlet_id"`
	Customer   Customer
	UpdatedAt  *time.Time `json:"updated_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
