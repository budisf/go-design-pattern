package routes

import (
	"ethical-be/app/config"

	"github.com/gin-gonic/gin"

	gt "ethical-be/modules/v1/utilities/zone/gt/handler"
)

func GroupTerritories(router *gin.Engine, conf *config.Conf, gtHandler *gt.GtHandler) {

	v1 := router.Group("/v1/group_territories")

	v1.POST("", gtHandler.CreateGt)
	v1.GET("", gtHandler.GetAll)
	v1.GET("/:id_group_territories", gtHandler.GetByID)
	v1.PUT("/:id_group_territories", gtHandler.UpdateGt)
	v1.DELETE("/:id_group_territories", gtHandler.Delete)
	v1.PATCH("/:id_group_territories/change_area", gtHandler.UpdateAreas)

}
