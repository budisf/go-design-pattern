package handler

import (
	"errors"
	"ethical-be/modules/v1/utilities/user/model"
	"ethical-be/modules/v1/utilities/user/service"
	res "ethical-be/pkg/api-response"
	helper "ethical-be/pkg/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IUserHandler interface {
	Save(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	UpdateById(ctx *gin.Context)
	GetAllPaginate(ctx *gin.Context)
	DeleteById(ctx *gin.Context)
	GetByAuthID(ctx *gin.Context)
}

type UserHandler struct {
	userService service.IUserService
}

func InitUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func responseUser(r *model.Users) *model.UserResponseEpochDTO {
	createdAtEpoch := uint64(helper.ConvertDateToUnix(r.CreatedAt))
	var updatedAtEpoch *uint64
	if r.UpdatedAt.IsZero() {
		updatedAtEpoch = nil
	} else {
		updatedAtEpochTemp := uint64(helper.ConvertDateToUnix(r.UpdatedAt))
		updatedAtEpoch = &updatedAtEpochTemp
	}
	userRoleResponse := model.UserResponseEpochDTO{
		ID:           r.ID,
		Name:         r.Name,
		AuthServerId: r.AuthServerId,
		Nip:          r.Nip,
		RoleId:       r.RoleId,
		Email:        r.Email,
		RoleName:     r.Role.Name,
		Label:        r.Role.Label,
		CreatedAt:    &createdAtEpoch,
		UpdatedAt:    updatedAtEpoch,
	}
	return &userRoleResponse
}

func (handler *UserHandler) Save(ctx *gin.Context) {
	var body model.UserRequestDTO
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
	httpResponse, errorMessage := handler.userService.Save(&body)

	if errorMessage != nil && httpResponse == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpResponse == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	ctx.SecureJSON(http.StatusOK, res.StatusOK("Successfully create the data user"))
	return
}

func (handler *UserHandler) GetByID(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	httpStatusResponse, errorMessage, dataUserRole := handler.userService.GetById(&idString)

	if httpStatusResponse == 500 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatusResponse == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	role := responseUser(dataUserRole)
	ctx.SecureJSON(http.StatusOK, res.Success(role))
	return
}

func (handler *UserHandler) UpdateById(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	var body model.UserRequestUpdateDTO
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

	httpResponse, errorMessage, dataUserUpdate := handler.userService.UpdateById(&idString, &body)

	if errorMessage != nil && httpResponse == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpResponse == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	dataUserResponse := responseUser(dataUserUpdate)
	ctx.SecureJSON(http.StatusOK, res.Success(dataUserResponse))
	return
}

func (handler *UserHandler) GetAllPaginate(ctx *gin.Context) {
	var httpStatus int
	// validate query params
	resultQueryParamPaginate, errResultQueryParamPaginate := helper.QueryParamPaginateTransform(ctx)
	if errResultQueryParamPaginate != nil {
		helper.HandleError(errResultQueryParamPaginate)
		httpStatus = http.StatusBadRequest
		ctx.SecureJSON(httpStatus, res.BadRequest(errResultQueryParamPaginate.Error()))
		return
	}
	fmt.Printf("resultQueryParamPaginate : %+v\n", *resultQueryParamPaginate.OrderBy)
	httpStatusResponse, errorMessage, responseJson := handler.userService.GetAllByPaginate(resultQueryParamPaginate)
	if errorMessage != nil || httpStatusResponse == 500 {
		helper.HandleError(errorMessage)
		httpStatus = http.StatusInternalServerError
		ctx.SecureJSON(httpStatus, res.ServerError(errorMessage.Error()))
		return
	}
	ctx.SecureJSON(httpStatus, res.Success(responseJson))
	return
}

func (handler *UserHandler) DeleteById(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}
	httpStatusResponse, errorMessage := handler.userService.DeleteById(&idString)

	if httpStatusResponse == 500 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatusResponse == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	ctx.SecureJSON(http.StatusOK, res.StatusOK("Successfully delete the data user "+idString))
	return

}

func (handler *UserHandler) GetByAuthID(ctx *gin.Context) {
	id, exist := ctx.Get("user_id")
	if exist == false {
		errMassage := errors.New("UserId authentication requested but middleware not registered on routes")
		helper.HandleError(errMassage)
		ctx.SecureJSON(http.StatusUnauthorized, res.UnAuthorized("Unauthorized"))
		return
	}
	idString := fmt.Sprintf("%v", id)

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	if errConvert != nil {
		helper.HandleError(errConvert)
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	httpStatusResponse, errorMessage, dataUserRole := handler.userService.GetById(&idString)

	if httpStatusResponse == 500 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatusResponse == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	role := responseUser(dataUserRole)
	ctx.SecureJSON(http.StatusOK, res.Success(role))
	return
}
