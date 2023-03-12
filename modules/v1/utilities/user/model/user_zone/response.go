package model

import "ethical-be/modules/v1/utilities/user/model"

type UserZoneDetailResponseDTO struct {
	ID            *uint         `json:"user_zone_id"`
	UserId        *uint         `json:"user_id"`
	Username      *string       `json:"username"`
	UserNip       *string       `json:"user_nip"`
	NameSalesZone *string       `json:"name_sales_zone"`
	SalesZoneId   *uint         `json:"sales_zone_id"`
	SalesZoneType SalesZoneType `json:"sales_zone_type"`
	AssignedDate  *uint64       `json:"assigned_date"`
}

type UserZoneDetailJoinUserResponseDTO struct {
	ID                         *uint         `json:"user_zone_id"`
	UserId                     *uint         `json:"user_id"`
	SalesZoneId                *uint         `json:"sales_zone_id"`
	SalesZoneType              SalesZoneType `json:"sales_zone_type"`
	AssignedDate               *uint64       `json:"assigned_date"`
	model.UserResponseEpochDTO `json:"users"`
}

type GetListUserByZoneResponseDTO struct {
	ID            *uint         `json:"user_zone_id"`
	UserId        *uint         `json:"user_id"`
	Username      *string       `json:"username"`
	UserNip       *string       `json:"user_nip"`
	SalesZoneId   *uint         `json:"sales_zone_id"`
	NameSalesZone *string       `json:"name_sales_zone"`
	SalesZoneType SalesZoneType `json:"sales_zone_type"`
	AssignedDate  *uint64       `json:"assigned_date"`
	FinishedDate  *uint64       `json:"finished_date"`
	CreatedAt     *uint64       `json:"created_at"`
	UpdatedAt     *uint64       `json:"updated_at"`
	UserRoleId    *uint         `json:"user_role_id"`
}

type ResponseListUserByZoneStatusResponseDTO struct {
	Vacant       *bool       `json:"vacant"`
	TotalRecords *uint       `json:"total_records"`
	Records      interface{} `json:"records"`
}

type UserResponse struct {
	UserId        uint    `json:"user_id"`
	UserName      *string `json:"user_name"`
	Nip           *string `json:"nip"`
	RoleID        *uint   `json:"role_id"`
	RoleName      *string `json:"role_name"`
	SalesZoneType string  `json:"zone_type"`
	SalesZoneId   *uint   `json:"zone_id"`
}

type GetZoneByUserIDResponse struct {
	SalesZoneID   uint   `json:"sales_zone_id"`
	SalesZoneType string `json:"sales_zone_type"`
	SalesZoneName string `json:"sales_zone_name"`
}

type ZoneType struct {
	ZoneType  string `json:"zone_type"`
	ZoneID    int    `json:"zone_id"`
	Name      string `json:"name"`
	IsVacant  bool   `json:"is_vacant"`
	RoleID    int    `json:"role_id"`
	RoleName  string `json:"role_name"`
	RoleLabel string `json:"role_label"`
}

type ZoneTypeRole struct {
	RoleID    int    `json:"role_id"`
	RoleName  string `json:"role_name"`
	RoleLabel string `json:"role_label"`
}
