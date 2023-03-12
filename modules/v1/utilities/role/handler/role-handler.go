package handler

import (
	"errors"
	model "ethical-be/modules/v1/utilities/role/model"
	"ethical-be/modules/v1/utilities/role/model/relations"
	"ethical-be/modules/v1/utilities/role/service"
	res "ethical-be/pkg/api-response"
	helper "ethical-be/pkg/helpers"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IRoleHandler interface {
	Save(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetAllPaginate(ctx *gin.Context)
	UpdateById(ctx *gin.Context)
	DeleteById(ctx *gin.Context)
	UpdateParentRole(ctx *gin.Context)
	GetChildPositionByUserId(c *gin.Context)
}

type RoleHandler struct {
	roleService service.IRoleService
}

func InitRoleHandler(roleService service.IRoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func responseRole(r *relations.RolesRelation) *model.RoleResponse {
	createdAtEpoch := uint64(helper.ConvertDateToUnix(r.CreatedAt))
	var parentRoleResponse *model.ParentRolesResponse
	var updatedAtEpoch *uint64
	if r.UpdatedAt.IsZero() {
		updatedAtEpoch = nil
	} else {
		updatedAtEpochTemp := uint64(helper.ConvertDateToUnix(r.UpdatedAt))
		updatedAtEpoch = &updatedAtEpochTemp
	}

	if r.ParentRole.ID == nil {
		parentRoleResponse = nil
	} else {

		createdRoleParentEpoch := uint64(helper.ConvertDateToUnix(*r.ParentRole.CreatedAt))
		var updatedRoleParentEpoch *uint64

		if r.ParentRole.UpdatedAt == nil {
			updatedRoleParentEpoch = nil
		} else {
			updatedRoleParentEpochTemp := uint64(helper.ConvertDateToUnix(*r.ParentRole.UpdatedAt))
			updatedRoleParentEpoch = &updatedRoleParentEpochTemp
		}

		parentRoleResponse = &model.ParentRolesResponse{
			RoleID:    r.ParentRole.ID,
			Label:     r.ParentRole.Label,
			CreatedAt: &createdRoleParentEpoch,
			UpdatedAt: updatedRoleParentEpoch,
		}
	}

	roleResponse := model.RoleResponse{
		RoleId:     r.ID,
		Name:       r.Name,
		Label:      r.Label,
		CreatedAt:  &createdAtEpoch,
		UpdatedAt:  updatedAtEpoch,
		ParentRole: parentRoleResponse,
	}
	return &roleResponse
}

func (handler *RoleHandler) Save(ctx *gin.Context) {
	var body model.RoleRequestDTO
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(errorMessages[0]))
			return
		}
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}
	httpResponse, errorMessage := handler.roleService.Save(&body)
	fmt.Println(httpResponse)
	if errorMessage != nil || httpResponse == 500 {
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpResponse == 404 {
		ctx.SecureJSON(http.StatusNotFound, res.NotFound("Data parent_id not found on table. Please check again"))
		return
	}
	ctx.SecureJSON(http.StatusOK, res.StatusOK("Successfully create the data role"))
	return
}

func (handler *RoleHandler) GetById(ctx *gin.Context) {
	idString := ctx.Param("role_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)

	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	httpStatusResponse, errorMessage, dataRole := handler.roleService.GetById(&idString)

	if errorMessage != nil || httpStatusResponse == 500 {
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	}

	if httpStatusResponse == 404 {
		ctx.SecureJSON(http.StatusNotFound, res.DataNotFound())
		return
	}
	role := responseRole(dataRole)
	ctx.SecureJSON(http.StatusOK, res.Success(role))
	return
}

func (handler *RoleHandler) GetAllPaginate(ctx *gin.Context) {
	var httpStatus int
	userIDAccessTokenInt := helper.GetUserIdFromMiddleware(ctx)
	userIDAccessToken := strconv.FormatUint(uint64(userIDAccessTokenInt), 10)
	// validate query params
	resultQueryParamPaginate, errResultQueryParamPaginate := helper.QueryParamPaginateTransform(ctx)
	if errResultQueryParamPaginate != nil {
		helper.HandleError(errResultQueryParamPaginate)
		httpStatus = http.StatusBadRequest
		ctx.SecureJSON(httpStatus, res.BadRequest(errResultQueryParamPaginate.Error()))
	}
	fmt.Printf("resultQueryParamPaginate : %+v\n", *resultQueryParamPaginate.OrderBy)
	httpStatusResponse, errorMessage, responseJson := handler.roleService.GetAllByPaginate(resultQueryParamPaginate, &userIDAccessToken)
	if errorMessage != nil || httpStatusResponse == 500 {
		helper.HandleError(errorMessage)
		httpStatus = http.StatusInternalServerError
		ctx.SecureJSON(httpStatus, res.ServerError(errorMessage.Error()))
		return
	}
	ctx.SecureJSON(httpStatus, res.Success(responseJson))
	return
}

func (handler *RoleHandler) UpdateById(ctx *gin.Context) {
	var body model.RoleRequestUpdateDTO
	idString := ctx.Param("role_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)

	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(errorMessages[0]))
			return
		}
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}

	httpStatusResponse, errorMessage, responseJson := handler.roleService.UpdateById(&idString, &body)
	if httpStatusResponse == 500 || errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	}
	if httpStatusResponse == 404 {
		ctx.SecureJSON(http.StatusNotFound, res.DataNotFound())
		return
	}
	ctx.SecureJSON(http.StatusOK, res.Success(responseJson))
	return
}

func (handler *RoleHandler) DeleteById(ctx *gin.Context) {
	idString := ctx.Param("role_id")
	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)

	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	httpStatusResponse, errorMessage := handler.roleService.DeleteById(&idString)
	if httpStatusResponse == 500 || errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	}
	if httpStatusResponse == 404 {
		ctx.SecureJSON(http.StatusNotFound, res.DataNotFound())
		return
	}
	ctx.SecureJSON(http.StatusOK, res.Success("Successfully deleted the data role"))
	return
}

func (handler *RoleHandler) UpdateParentRole(ctx *gin.Context) {
	var bodyPayload model.RoleRequestUpdateHeadRoleDTO
	var idHeadRoleString *string
	idRoleString := ctx.Param("role_id")
	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idRoleString)

	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	err := ctx.ShouldBindJSON(&bodyPayload)
	if err != nil {
		fmt.Println(err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			ctx.JSON(http.StatusBadRequest, res.BadRequest(errorMessages[0]))
			return
		}
		ctx.JSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}

	if bodyPayload.IdHeadRole != nil {
		idHeadRoleUint64 := uint64(*bodyPayload.IdHeadRole)
		idHeadRoleUint64Temp := strconv.FormatUint(idHeadRoleUint64, 10)
		idHeadRoleString = &idHeadRoleUint64Temp
	} else {
		idHeadRoleString = nil
	}

	httpStatusResponse, errorMessage, responseJson := handler.roleService.UpdateChangeParentRole(&idRoleString, idHeadRoleString)
	if httpStatusResponse == 404 {
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	if httpStatusResponse == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	}
	role := responseRole(responseJson)
	ctx.JSON(http.StatusOK, res.Success(role))
	return
}

func (handler *RoleHandler) GetChildPositionByUserId(c *gin.Context) {

	userId := helper.GetUserIdFromMiddleware(c)
	userIdString := fmt.Sprintf("%v", userId)

	result, err := handler.roleService.GetChildPositionByUserId(userIdString)

	if err != nil {
		if err.Error() == "404" {
			c.JSON(http.StatusNotFound, res.NotFound("Role_id User"))
			return
		}
		c.JSON(http.StatusInternalServerError, res.ServerError(err.Error()))
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, res.DataNotFound())
		return
	}

	c.JSON(http.StatusOK, res.Success(result))
}
