package models

type AllRequest struct {
	Name    string `json:"name" binding:"required"`
	AreaID int    `json:"area_id" binding:"required"`
}

type GtRequest struct {
	Name string `json:"name" binding:"required"`
}

type AreaRequest struct {
	AreaID int `json:"area_id" binding:"required"`
}
