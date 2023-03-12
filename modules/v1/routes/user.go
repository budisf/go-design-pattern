package routes

import (
	"ethical-be/app/config"
	"ethical-be/app/middlewares"
	"ethical-be/modules/v1/utilities/user/handler"

	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine, conf *config.Conf, userHandler *handler.UserHandler, userRoleHandler *handler.UserRoleHandler, userZoneHandler *handler.UserZoneHandler) {

	/*
	 auth middleware
	*/
	auth := middlewares.AuthUser()
	/*
		API USER MASTER
	*/
	v1 := router.Group("/v1/user")
	/*
		@routes: {name_service}/v1/user
		@desc: to create new user
	*/
	v1.POST("", userHandler.Save)
	/*
		@routes: {name_service}/v1/user/:user_id
		@desc: to get detail data user
	*/
	v1.GET("/:user_id", userHandler.GetByID)
	/*
		@routes: {name_service}/v1/user/:user_id
		@desc: to update detail data user
	*/
	v1.PATCH("/:user_id", userHandler.UpdateById)
	/*
		@routes: {name_service}/v1/user/
		@desc: to get all data user with pagination or not
	*/
	v1.GET("", userHandler.GetAllPaginate)
	/*
		@routes: {name_service}/v1/user/:user_id
		@desc: to soft delete data user by id
	*/
	v1.DELETE("/:user_id", userHandler.DeleteById)
	/*
		API USER ROLE
	*/
	/*
		@routes: {name_service}/v1/user/:user_id/role
		@desc: to assign role user
	*/
	v1.PATCH("/:user_id/role", userRoleHandler.UpdateById)
	/*
		@routes: {name_service}/v1/user/:user_id/role
		@desc: to assign role user
	*/
	v1.GET("/:user_id/role", userRoleHandler.GetByID)
	/*
		@routes: {name_service}/v1/user/:user_id/zone
		@desc: to assign user zone based on gt, or area, or region
	*/
	v1.POST("/:user_id/zone", userZoneHandler.AssignUserZone)
	/*
		@routes: {name_service}/v1/user/:user_id/zone
		@desc: to get detail user by zone assign and not vacant
	*/
	v1.GET("/zone", auth, userZoneHandler.GetZoneByUserID)
	/*
		@routes: {name_service}/v1/user/:user_id/zone
		@desc: to finished user by zone assign and result is vacant
	*/
	v1.PATCH("/:user_id/zone/finished-assignment", userZoneHandler.UpdateFinishedAssigment)
	/*
		@routes: {name_service}/v1/user/:zone_id/zone-status
		@desc: to get detail user by zone assign and not vacant
	*/
	v1.GET("/zone-status", userZoneHandler.GetListUserByZoneId)
	/*
		@routes: {name_service}/v1/user/child
		@desc: to get child user by user id and role name
	*/
	v1.GET("/child", auth, userZoneHandler.GetChildByUser)
	/*
		@routes: {name_service}/v1/user/:user_id/zone
		@desc: to get detail user by zone assign and not vacant
	*/
	v1.GET("/:user_id/zone-detail", userZoneHandler.GetUserZoneByUserIDZoneIDZoneType)
	/*
		@routes: {name_service}/v1/user/child/vacant
		@desc: to get child vacant by user id
	*/
	v1.GET("/child/vacant", auth, userZoneHandler.GetChildVacantByUserId)
	/*
		@routes: {name_service}/v1/user/child/nonvacant
		@desc: to get child non vacant by user id
	*/
	v1.GET("/child/nonvacant", auth, userZoneHandler.GetChildNonVacantByUserId)
	/*
		@routes: {name_service}/v1/user/:user_id/subordinate-employees?sales_zone_id=1&sales_zone_type=
		@desc: to get child user by user id and role name
	*/
	v1.GET("/:user_id/subordinate-employees", userZoneHandler.GetSubordinateEmployeesByUserIDZoneIDZoneType)
	/*
		@routes: {name_service}/v1/user/profile
		@desc: to update detail data user
	*/
	v1.GET("/profile", auth, userHandler.GetByAuthID)

	v1.GET("/zone/vacant-child", auth, userZoneHandler.GetZoneChildVacantByUserId)

	v1.GET("/zone/impersonate", auth, userZoneHandler.ImpersonateAccessControlSales)
}
