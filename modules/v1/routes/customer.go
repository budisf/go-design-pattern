package routes

import (
	"ethical-be/app/config"
	// "ethical-be/app/middlewares"

	"github.com/gin-gonic/gin"

	customer "ethical-be/modules/v1/utilities/customer/handler"
)

func Customer(router *gin.Engine, conf *config.Conf, customerHandler customer.CustomerHandler) {

	/**
	All customer Route
	*/

	v1 := router.Group("/v1/customer")

	/*
	 jwt middleware
	*/
	// auth := middlewares.AuthJwt()

	v1.GET("", customerHandler.Index)
	v1.GET("/:id", customerHandler.GetById)
	v1.POST("", customerHandler.Create)
	v1.PUT("/:id", customerHandler.Edit)
	v1.DELETE("/:id", customerHandler.Delete)

}
