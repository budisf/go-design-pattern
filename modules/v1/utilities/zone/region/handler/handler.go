package handler

import (
	model "ethical-be/modules/v1/utilities/zone/region/models"
	service "ethical-be/modules/v1/utilities/zone/region/services"
	respon "ethical-be/pkg/api-response"

	helper "ethical-be/pkg/helpers"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegionsHandler struct {
	regionsService service.RegionsService
}

func NewRegionsHandler(regionsService service.RegionsService) *RegionsHandler {
	return &RegionsHandler{regionsService}
}

func (handler *RegionsHandler) GetAll(ctx *gin.Context) {
	regions, err := handler.regionsService.FindAll()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, respon.DataNotFound())
			return
		}
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	var regionsResponse []model.RegionsResponse
	for _, r := range regions {
		regionResponse := responseRegion(r)
		regionsResponse = append(regionsResponse, regionResponse)
	}
	ctx.JSON(http.StatusOK, respon.Success(regionsResponse))
}

func (handler *RegionsHandler) GetByID(ctx *gin.Context) {
	idString := ctx.Param("id_region")
	id, _ := strconv.Atoi(idString)

	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Regions"))
		return
	}

	regions, err := handler.regionsService.FindById(id)
	if regions.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Regions"))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	// regionsResponse := responseRegion(regions.Regions)
	ctx.JSON(http.StatusOK, respon.Success(regions))
}

func (handler *RegionsHandler) CreateRegions(ctx *gin.Context) {
	var regionsRequest model.RegionsRequest

	err := ctx.ShouldBindJSON(&regionsRequest)
	if err != nil {
		errorMessage := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessage))
		return
	}
	regions, err := handler.regionsService.Create(regionsRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, respon.Success(responseRegion(regions)))
}

func (handler *RegionsHandler) UpdateRegions(ctx *gin.Context) {

	var regionRequest model.RegionsRequest

	err := ctx.ShouldBindJSON(&regionRequest)
	if err != nil {
		errorMessages := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessages))
		return
	}
	idString := ctx.Param("id_region")
	id, _ := strconv.Atoi(idString)
	r, errBy := handler.regionsService.FindById(id)
	if r.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID"))
		return
	}
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	regions, err := handler.regionsService.Update(id, regionRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, respon.Success(regions))
}

func (handler *RegionsHandler) DeleteRegions(ctx *gin.Context) {
	idString := ctx.Param("id_region")
	id, _ := strconv.Atoi(idString)
	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Regions"))
		return
	}
	regions, err := handler.regionsService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	if regions.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Regions"))
		return
	}

	errBy := handler.regionsService.Delete(id)
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(errBy.Error()))
		return
	}
	ctx.JSON(http.StatusOK, respon.StatusOK("Delete Success"))
}

func responseRegion(r model.Regions) model.RegionsResponse {
	regionResponse := model.RegionsResponse{
		ID:         r.ID,
		Name:       r.Name,
		DistrictID: r.DistrictID,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		DeletedAt:  r.DeletedAt,
		IsVacant:   r.IsVacant,
	}
	return regionResponse
}
