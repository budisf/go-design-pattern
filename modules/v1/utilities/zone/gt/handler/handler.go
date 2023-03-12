package handler

import (
	area "ethical-be/modules/v1/utilities/zone/area/services"
	model "ethical-be/modules/v1/utilities/zone/gt/models"
	service "ethical-be/modules/v1/utilities/zone/gt/services"
	respon "ethical-be/pkg/api-response"

	helper "ethical-be/pkg/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GtHandler struct {
	gtService   service.GtServices
	areaService area.AreasService
}

func NewGtHandler(gtService service.GtServices, areaService area.AreasService) *GtHandler {
	return &GtHandler{gtService, areaService}
}

func (handler *GtHandler) GetAll(ctx *gin.Context) {
	gts, err := handler.gtService.FindAll()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, respon.DataNotFound())
			return
		}
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	var gtsResponse []model.GtResponse
	for _, gt := range gts {
		gtResponse := responseGt(gt)
		gtsResponse = append(gtsResponse, gtResponse)
	}
	ctx.JSON(http.StatusOK, respon.Success(gtsResponse))
}

func (handler *GtHandler) GetByID(ctx *gin.Context) {

	idString := ctx.Param("id_group_territories")
	id, _ := strconv.Atoi(idString)

	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Group Territories"))
		return
	}

	gt, err := handler.gtService.FindById(id)
	if gt.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Group Territories"))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	// areaResponse := responseGt(gt)
	ctx.JSON(http.StatusOK, respon.Success(gt))
}

func (handler *GtHandler) CreateGt(ctx *gin.Context) {
	var gtRequest model.AllRequest

	err := ctx.ShouldBindJSON(&gtRequest)
	if err != nil {
		errorMassages := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMassages))
		return
	}

	cekAreas, _ := handler.areaService.FindById(gtRequest.AreaID)
	if cekAreas.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("Area ID"))
		return
	}

	gts, _ := handler.gtService.Create(gtRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(err))
		return
	}
	ctx.JSON(http.StatusOK, respon.Success(responseGt(gts)))
}

func (handler *GtHandler) UpdateGt(ctx *gin.Context) {
	var gtRequest model.GtRequest

	err := ctx.ShouldBindJSON(&gtRequest)
	if err != nil {
		errorMessages := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessages))
		return
	}
	idString := ctx.Param("id_group_territories")
	id, _ := strconv.Atoi(idString)
	g, errBy := handler.gtService.FindById(id)
	if g.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID group territories"))
		return
	}
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	gt, err := handler.gtService.Update(id, gtRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, respon.Success(responseGt(gt)))
}

func (handler *GtHandler) UpdateAreas(ctx *gin.Context) {
	var areaRequest model.AreaRequest

	err := ctx.ShouldBindJSON(&areaRequest)
	if err != nil {
		errorMessage := helper.ErrorMessage(err)
		ctx.JSON(http.StatusBadRequest, respon.BadRequest(errorMessage))
		return
	}
	idString := ctx.Param("id_group_territories")
	id, _ := strconv.Atoi(idString)
	g, errBy := handler.gtService.FindById(id)
	if g.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID group territories"))
		return
	}
	if errBy != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	gt, err := handler.gtService.UpdateAreas(id, areaRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, respon.Success(responseGt(gt)))
}

func (handler *GtHandler) Delete(ctx *gin.Context) {
	idString := ctx.Param("id_group_territories")
	id, _ := strconv.Atoi(idString)

	if id == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Group Territories"))
		return
	}
	Gt, err := handler.gtService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}
	if Gt.ID == 0 {
		ctx.JSON(http.StatusNotFound, respon.NotFound("ID Group Territories"))
		return
	}
	err = handler.gtService.Delete(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, respon.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, respon.StatusOK("Delete Success"))
}

func responseGt(g model.GroupTerritories) model.GtResponse {
	gtResponse := model.GtResponse{
		ID:        g.ID,
		Name:      g.Name,
		AreaID:    g.AreaID,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		DeletedAt: g.DeletedAt,
		IsVacant:  g.IsVacant,
	}
	return gtResponse
}
