package handler

import (
	model "ethical-be/modules/v1/utilities/zone/district/models"
	service "ethical-be/modules/v1/utilities/zone/district/services"
	respon "ethical-be/pkg/api-response"

	helper "ethical-be/pkg/helpers"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DistrictHandler struct {
	districtService service.RegionsService
}

func NewDistrictHandler(districtService service.RegionsService) *DistrictHandler {
	return &DistrictHandler{districtService}
}

func (handler *DistrictHandler) GetAll(ctx *gin.Context) {
	district, err := handler.districtService.FindAll()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, respon.DataNotFound())
			return
		}
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	var districtsResponse []model.DistrictResponse
	for _, d := range district {
		districtResponse := responseDistrict(d)
		districtsResponse = append(districtsResponse, districtResponse)
	}
	ctx.JSON(http.StatusOK, respon.Success(districtsResponse))
}

func (handler *DistrictHandler) GetByID(ctx *gin.Context) {
	idString := ctx.Param("id_district")
	id, _ := strconv.Atoi(idString)

	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID District"))
		return
	}

	district, err := handler.districtService.FindById(id)
	if district.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID District"))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	// districtResponse := responseDistrict(district)
	ctx.JSON(http.StatusOK, respon.Success(district))
}

func (handler *DistrictHandler) CreateDistrict(ctx *gin.Context) {
	var districtRequest model.DistrictRequest

	err := ctx.ShouldBindJSON(&districtRequest)
	if err != nil {
		errorMessage := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessage))
		return
	}
	district, err := handler.districtService.Create(districtRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, respon.Success(responseDistrict(district)))
}

func (handler *DistrictHandler) UpdateDistrict(ctx *gin.Context) {

	var districtRequest model.DistrictRequest

	err := ctx.ShouldBindJSON(&districtRequest)
	if err != nil {
		errorMessages := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessages))
		return
	}
	idString := ctx.Param("id_district")
	id, _ := strconv.Atoi(idString)
	r, errBy := handler.districtService.FindById(id)
	if r.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID"))
		return
	}
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	district, err := handler.districtService.Update(id, districtRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, respon.Success(responseDistrict(district)))
}

func (handler *DistrictHandler) DeleteDistrict(ctx *gin.Context) {
	idString := ctx.Param("id_district")
	id, _ := strconv.Atoi(idString)
	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID District"))
		return
	}
	district, err := handler.districtService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	if district.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID District"))
		return
	}

	errBy := handler.districtService.Delete(id)
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(errBy.Error()))
		return
	}
	ctx.JSON(http.StatusOK, respon.StatusOK("Delete Success"))
}

func responseDistrict(d model.Districts) model.DistrictResponse {
	districtResponse := model.DistrictResponse{
		ID:        d.ID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,
		IsVacant:  d.IsVacant,
	}
	return districtResponse
}
