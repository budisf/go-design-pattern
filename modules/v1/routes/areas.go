package routes

import (
	"ethical-be/app/config"

	"github.com/gin-gonic/gin"

	area "ethical-be/modules/v1/utilities/zone/area/handler"
)

func Areas(router *gin.Engine, conf *config.Conf, areaHandler *area.AreasHandler) {

	v1 := router.Group("/v1/area")

	v1.POST("/", areaHandler.CreateAreas)
	v1.GET("/", areaHandler.GetAll)
	v1.GET("/:id_area", areaHandler.GetByID)
	v1.PUT("/:id_area", areaHandler.UpdateArea)
	v1.DELETE("/:id_area", areaHandler.Delete)
	v1.PUT("/:id_area/change_region", areaHandler.UpdateRegion)

}
