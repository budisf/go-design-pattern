package routes

import (
	"ethical-be/app/config"

	"github.com/gin-gonic/gin"

	district "ethical-be/modules/v1/utilities/zone/district/handler"
)

func District(router *gin.Engine, conf *config.Conf, districtHandler *district.DistrictHandler) {

	v1 := router.Group("/v1/district")

	v1.POST("", districtHandler.CreateDistrict)
	v1.GET("", districtHandler.GetAll)
	v1.GET("/:id_district", districtHandler.GetByID)
	v1.PUT("/:id_district", districtHandler.UpdateDistrict)
	v1.DELETE("/:id_district", districtHandler.DeleteDistrict)

}
