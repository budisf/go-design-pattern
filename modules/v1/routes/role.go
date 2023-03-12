package routes

import (
	"ethical-be/app/config"
	"ethical-be/app/middlewares"
	"ethical-be/modules/v1/utilities/role/handler"

	"github.com/gin-gonic/gin"
)

func Role(router *gin.Engine, conf *config.Conf, roleHandler *handler.RoleHandler) {

	auth := middlewares.AuthUser()
	v1 := router.Group("/v1/role")
	/*
		@routes: {name_service}/v1/role
		@desc: get all the data pagination
	*/
	v1.GET("", auth, roleHandler.GetAllPaginate)
	/*
		@routes: {name_service}/v1/role
		@desc: to create new role
	*/
	v1.POST("", roleHandler.Save)
	/*
		@routes: {name_service}/v1/role/:role_id
		@desc: to get detail role
	*/
	v1.GET("/:role_id", roleHandler.GetById)
	/*
		@routes: {name_service}/v1/role/:role_id
		@desc: to update role
	*/
	v1.PATCH("/:role_id", roleHandler.UpdateById)
	/*
		@routes: {name_service}/v1/role/:role_id
		@desc: to delete role
	*/
	v1.DELETE("/:role_id", roleHandler.DeleteById)
	/*
		@routes: {name_service}/v1/role/:role_id
		@desc: to update role with change head role
	*/
	v1.PATCH("/:role_id/change-role-head", roleHandler.UpdateParentRole)
	/*
		@routes: {name_service}/v1/child/:user_id
		@desc: to get all child position by user id
	*/
	v1.GET("/child", auth, roleHandler.GetChildPositionByUserId)
}
