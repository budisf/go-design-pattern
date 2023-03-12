package models

type AllRequest struct {
	Name      string `json:"name" binding:"required"`
	RegionID int    `json:"region_id" binding:"required"`
}

type AreasRequest struct {
	Name string `json:"name" binding:"required"`
}

type RegionRequest struct {
	RegionID int `json:"region_id" binding:"required"`
}
