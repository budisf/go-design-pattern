package models

type RegionsRequest struct {
	Name       string `json:"name" binding:"required"`
	DistrictID int    `json:"district_id" binding:"required"`
}
