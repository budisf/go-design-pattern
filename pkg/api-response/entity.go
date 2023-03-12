package apiresponse

import helperDatabases "ethical-be/pkg/helpers/databases"

/*
   |--------------------------------------------------------------------------
   | Default Struct Response API
   |--------------------------------------------------------------------------
   |
   | You can chang this every momen you want    |
*/

type Meta struct {
	Message interface{} `json:"message"`
	Status  int         `json:"status"`
	Code    string      `json:"code"`
}

type Response struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
}

/*
   |--------------------------------------------------------------------------
   | Default Struct Pagination data
   |--------------------------------------------------------------------------
   |
*/

type PaginationQuery struct {
	Limit  *int `form:"limit,default=10" binding:"required,number"`
	Page   *int `form:"page,default=1" binding:"required,number"`
	Offset *int `form:"offset,default=0" binding:"omitempty,number"`
}

type Pagination struct {
	MetaData *helperDatabases.ResponseBackPaginationDTO `json:"metadata"`
	Records  interface{}                                `json:"records"`
}

type PaginationOld struct {
	CurrentPage int         `json:"current_page"`
	TotalPage   int         `json:"total_page"`
	TotalData   int         `json:"total_data"`
	DataPerPage interface{} `json:"data_per_page"`
}
