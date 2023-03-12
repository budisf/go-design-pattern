package routes

import (
	"ethical-be/app/config"

	"github.com/gin-gonic/gin"

	region "ethical-be/modules/v1/utilities/zone/region/handler"
)

func Regions(router *gin.Engine, conf *config.Conf, regionsHandler *region.RegionsHandler) {

	v1 := router.Group("/v1/region")

	v1.POST("", regionsHandler.CreateRegions)
	v1.GET("", regionsHandler.GetAll)
	v1.GET("/:id_region", regionsHandler.GetByID)
	v1.PUT("/:id_region", regionsHandler.UpdateRegions)
	v1.DELETE("/:id_region", regionsHandler.DeleteRegions)

}
