package handler

import (
	model "ethical-be/modules/v1/utilities/zone/area/models"
	service "ethical-be/modules/v1/utilities/zone/area/services"
	reg "ethical-be/modules/v1/utilities/zone/region/services"
	respon "ethical-be/pkg/api-response"

	helper "ethical-be/pkg/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AreasHandler struct {
	areasService   service.AreasService
	regionsService reg.RegionsService
}

func NewAreasHandler(areasService service.AreasService, regionsService reg.RegionsService) *AreasHandler {
	return &AreasHandler{areasService, regionsService}
}

func (handler *AreasHandler) GetAll(ctx *gin.Context) {
	areas, err := handler.areasService.FindAll()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, respon.DataNotFound())
			return
		}
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	var areasResponse []model.AreasResponse
	for _, a := range areas {
		areaResponse := responseAreas(a)
		areasResponse = append(areasResponse, areaResponse)
	}
	ctx.JSON(http.StatusOK, respon.Success(areasResponse))
}

func (handler *AreasHandler) GetByID(ctx *gin.Context) {

	idString := ctx.Param("id_area")
	id, _ := strconv.Atoi(idString)

	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Area"))
		return
	}

	areas, err := handler.areasService.FindById(id)
	if areas.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Area"))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	// areaResponse := responseAreas(areas)
	ctx.JSON(http.StatusOK, respon.Success(areas))
}

func (handler *AreasHandler) CreateAreas(ctx *gin.Context) {
	var areasRequest model.AllRequest

	err := ctx.ShouldBindJSON(&areasRequest)
	if err != nil {
		errorMassages := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMassages))
		return
	}

	cekReg, _ := handler.regionsService.FindById(areasRequest.RegionID)
	if cekReg.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("Regions ID"))
		return
	}

	areas, _ := handler.areasService.Create(areasRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(err))
		return
	}
	ctx.JSON(http.StatusOK, respon.Success(responseAreas(areas)))
}

func (handler *AreasHandler) UpdateArea(ctx *gin.Context) {
	var areaRequest model.AreasRequest

	err := ctx.ShouldBindJSON(&areaRequest)
	if err != nil {
		errorMessages := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessages))
		return
	}
	idString := ctx.Param("id_area")
	id, _ := strconv.Atoi(idString)
	a, errBy := handler.areasService.FindById(id)
	if a.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID"))
		return
	}
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	areas, err := handler.areasService.Update(id, areaRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, respon.Success(responseAreas(areas)))
}

func (handler *AreasHandler) UpdateRegion(ctx *gin.Context) {
	var regRequest model.RegionRequest

	err := ctx.ShouldBindJSON(&regRequest)
	if err != nil {
		errorMessages := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessages))
		return
	}
	idString := ctx.Param("id_area")
	id, _ := strconv.Atoi(idString)
	a, errBy := handler.areasService.FindById(id)
	if a.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID"))
		return
	}
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	reg, err := handler.areasService.UpdateRegion(id, regRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, respon.Success(responseAreas(reg)))
}

func (handler *AreasHandler) Delete(ctx *gin.Context) {
	idString := ctx.Param("id_area")
	id, _ := strconv.Atoi(idString)

	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Area"))
		return
	}
	Areas, err := handler.areasService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	if Areas.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Area"))
		return
	}
	err = handler.areasService.Delete(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, respon.StatusOK("Delete Success"))
}

func responseAreas(a model.Areas) model.AreasResponse {
	areaResponse := model.AreasResponse{
		ID:        a.ID,
		Name:      a.Name,
		RegionID:  a.RegionID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		DeletedAt: a.DeletedAt,
		IsVacant:  a.IsVacant,
	}
	return areaResponse
}
