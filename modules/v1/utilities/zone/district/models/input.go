package models

type DistrictRequest struct {
	Name string `json:"name" binding:"required"`
}
