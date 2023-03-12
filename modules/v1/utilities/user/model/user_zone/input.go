package model

type UserZoneRequestDTO struct {
	SalesZoneId   *uint   `json:"sales_zone_id" binding:"required,gte=1"`
	SalesZoneType *string `json:"sales_zone_type" binding:"EnumVersionTwo=group_territories&areas&regions&districts"`
}

type UserZoneRequestQueryParamDTO struct {
	SalesZoneId   *string `binding:"required,gte=1" form:"sales_zone_id"`
	SalesZoneType *string `binding:"EnumVersionTwo=group_territories&areas&regions&districts" form:"sales_zone_type"`
}

type UserZoneRequestParamsDTO struct {
	SalesZoneId   *uint   `form:"sales_zone_id" binding:"required,gte=1"`
	SalesZoneType *string `form:"sales_zone_type" binding:"required,EnumVersionTwo=group_territories&areas&regions&districts"`
	RoleName      *string `form:"role_name" binding:"required,EnumVersionTwo=nsm&sm&asm&field-force"`
}

type GetZoneByUserID struct {
	SalesZoneType *string `form:"sales_zone_type" binding:"omitempty,EnumVersionTwo=group_territories&areas&regions&districts"`
	UserId        *uint   `json:"user_id" form:"user_id" binding:"omitempty,gte=1"`
}
