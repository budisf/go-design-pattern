package handler

import (
	"errors"
	model2 "ethical-be/modules/v1/utilities/user/model"
	model3 "ethical-be/modules/v1/utilities/user/model/user_role"
	"ethical-be/modules/v1/utilities/user/service"
	res "ethical-be/pkg/api-response"
	helper "ethical-be/pkg/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IUserRoleHandler interface {
	GetByID(ctx *gin.Context)
	UpdateById(ctx *gin.Context)
}

type UserRoleHandler struct {
	userRoleService service.IUserRoleService
}

func InitUserRoleHandler(userRoleService service.IUserRoleService) *UserRoleHandler {
	return &UserRoleHandler{
		userRoleService: userRoleService,
	}
}

/*
	|--------------------------------------------------------------------------
	| to convert entity to response
	|--------------------------------------------------------------------------
	|
	| This function is for return new dto from raw query user
	|

|
*/
func responseUserRole(r *model3.UserRolesRawQueryResult) *model2.UserRoleResponseDTO {
	userRoleResponse := model2.UserRoleResponseDTO{
		ID:           r.ID,
		Name:         r.Name,
		AuthServerId: r.AuthServerId,
		Nip:          r.Nip,
		RoleId:       r.RoleId,
		RoleName:     r.RoleName,
		RoleLabel:    r.RoleLabel,
		ParentId:     r.ParentId,
		NameRoleHead: r.NameRoleHead,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
	return &userRoleResponse
}

func (handler *UserRoleHandler) UpdateById(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	var body model3.UserRoleRequestDTO
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

	httpResponse, errorMessage, dataUserUpdate := handler.userRoleService.UpdateById(&idString, &body)

	if errorMessage != nil && httpResponse == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpResponse == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	dataUserResponse := responseUserRole(dataUserUpdate)
	ctx.SecureJSON(http.StatusOK, res.Success(dataUserResponse))
	return
}

func (handler *UserRoleHandler) GetByID(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}
	httpResponse, errorMessage, dataUser := handler.userRoleService.GetById(&idString)

	if errorMessage != nil && httpResponse == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpResponse == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	dataUserResponse := responseUserRole(dataUser)
	ctx.SecureJSON(http.StatusOK, res.Success(dataUserResponse))
	return
}
