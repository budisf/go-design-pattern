package models

type CustomerResponse struct {
	ID           uint
	CustomerName string  `json:"customer_name" gorm:"type:varchar(256)" binding:"required"`
	CustomerCode string  `json:"customer_code" gorm:"type:varchar(256)"`
	Specialist   *string `json:"specialis" gorm:"type:varchar(256)"`
	UpdatedAt    uint    `json:"updated_at"`
	CreatedAt    uint    `json:"created_at"`
}
