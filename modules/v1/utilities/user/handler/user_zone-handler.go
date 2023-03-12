package handler

import (
	"errors"
	model2 "ethical-be/modules/v1/utilities/user/model"
	model "ethical-be/modules/v1/utilities/user/model/user_zone"
	"ethical-be/modules/v1/utilities/user/service"
	res "ethical-be/pkg/api-response"
	"ethical-be/pkg/helpers"
	helper "ethical-be/pkg/helpers"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IUserZoneHandler interface {
	AssignUserZone(ctx *gin.Context)
	GetUserZoneByUserID(ctx *gin.Context) // delete endpoint ini
	UpdateFinishedAssigment(ctx *gin.Context)
	GetListUserByZoneId(ctx *gin.Context)
	GetChildByUser(ctx *gin.Context)
	GetUserZoneByUserIDZoneIDZoneType(ctx *gin.Context)
	GetChildVacantByUserId(ctx *gin.Context)
	GetSubordinateEmployeesByUserIDZoneIDZoneType(ctx *gin.Context)
	GetChildNonVacantByUserId(ctx *gin.Context)
	GetZoneByUserID(ctx *gin.Context)
	dryCheckGetZoneByUserID(userIDAccessToken string, params model.GetZoneByUserID) (uint, error, model.GetZoneByUserID)
	// algorithms customize
	binarySearch(data []int, search int) (result int)
	extractSliceUserResponseIntoSliceInt(data []model.UserZoneWithEpochEntity) []int
}

type UserZoneHandler struct {
	userZoneService service.IUserZoneService
	userRoleService service.IUserRoleService
}

func InitUserZoneHandler(userZoneService service.IUserZoneService, userRoleService service.IUserRoleService) *UserZoneHandler {
	return &UserZoneHandler{
		userZoneService: userZoneService,
		userRoleService: userRoleService,
	}
}

func responseUserZoneDetail(value *model.GetBySalesZoneIdUserIdRawQuery) *model.UserZoneDetailResponseDTO {
	var userZoneDetail model.UserZoneDetailResponseDTO

	userZoneDetail = model.UserZoneDetailResponseDTO{
		ID:            value.ID,
		UserId:        value.UserId,
		Username:      value.Username,
		UserNip:       value.Username,
		NameSalesZone: value.NameSalesZone,
		SalesZoneId:   value.SalesZoneId,
		SalesZoneType: value.SalesZoneType,
		AssignedDate:  value.AssignedDate,
	}

	return &userZoneDetail
}

func responseUserZoneJoinUserDetail(user *model.UserZone) *model.UserZoneDetailJoinUserResponseDTO {
	var userZoneJoinUserDetail model.UserZoneDetailJoinUserResponseDTO
	var assignedDataUint64 uint64 = uint64(helper.ConvertDateToUnix(*user.AssignedDate))
	var createdUserUint64 uint64 = uint64(helper.ConvertDateToUnix(user.Users.CreatedAt))
	var updatedUserUint64 uint64 = uint64(helper.ConvertDateToUnix(user.Users.UpdatedAt))
	userZoneJoinUserDetail = model.UserZoneDetailJoinUserResponseDTO{
		ID:            user.ID,
		UserId:        user.UserId,
		SalesZoneId:   user.SalesZoneId,
		SalesZoneType: user.SalesZoneType,
		AssignedDate:  &assignedDataUint64,
		UserResponseEpochDTO: model2.UserResponseEpochDTO{
			ID:           user.Users.ID,
			Name:         user.Users.Name,
			AuthServerId: user.Users.AuthServerId,
			Nip:          user.Users.Nip,
			RoleId:       user.Users.RoleId,
			CreatedAt:    &createdUserUint64,
			UpdatedAt:    &updatedUserUint64,
		},
	}
	return &userZoneJoinUserDetail
}

func (handler *UserZoneHandler) AssignUserZone(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	var body model.UserZoneRequestDTO

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

	httpStatus, errorMessage := handler.userZoneService.AssignUserZone(&idString, &body)
	if errorMessage != nil && httpStatus == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatus == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	ctx.SecureJSON(http.StatusOK, res.StatusOK("Successfully Assign User"))
	return
}

func (handler *UserZoneHandler) GetUserZoneByUserID(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	httpStatus, errorMessage, data := handler.userZoneService.GetUserZoneByUserId(&idString)
	if errorMessage != nil && httpStatus == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatus == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}

	userZoneDetail := responseUserZoneDetail(data)

	ctx.SecureJSON(http.StatusOK, res.Success(userZoneDetail))
	return
}

func (handler *UserZoneHandler) UpdateFinishedAssigment(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}
	var body model.UserZoneRequestDTO

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

	salesZoneIdString := strconv.FormatUint(uint64(*body.SalesZoneId), 10)

	httpStatus, errorMessage := handler.userZoneService.UpdateFinishedAssigment(&idString, &body)
	if errorMessage != nil && httpStatus == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatus == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	message := fmt.Sprintf("Successfully finished assigment user_id %v with sales zone id %v", idString, salesZoneIdString)
	ctx.SecureJSON(http.StatusOK, res.StatusOK(message))
	return
}

func (handler *UserZoneHandler) GetListUserByZoneId(ctx *gin.Context) {
	var body model.UserZoneRequestQueryParamDTO

	if err := ctx.ShouldBindQuery(&body); err != nil {
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

	salesZoneIdUint64, errConvert := helper.StringToUint64(body.SalesZoneId)

	if errConvert != nil {
		helper.HandleError(errConvert)
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	salesZoneIdString := strconv.FormatUint(*salesZoneIdUint64, 10)

	httpStatus, errorMessage, data := handler.userZoneService.GetListUserByZoneId(&salesZoneIdString, body.SalesZoneType)

	if errorMessage != nil && httpStatus == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatus == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}

	var statusVacant bool
	var totalRecord uint
	var responseResultQuery interface{}

	if data.ID == nil {
		statusVacant = true
		totalRecord = 0
		responseResultQuery = res.EmptyObj{}
	} else {
		statusVacant = false
		totalRecord = 1
		responseResultQuery = model.GetListUserByZoneResponseDTO{
			ID:            data.ID,
			UserId:        data.UserId,
			Username:      data.Username,
			UserNip:       data.UserNip,
			SalesZoneId:   data.SalesZoneId,
			SalesZoneType: data.SalesZoneType,
			AssignedDate:  data.AssignedDate,
			FinishedDate:  data.FinishedDate,
			CreatedAt:     data.CreatedAt,
			UpdatedAt:     data.UpdatedAt,
			UserRoleId:    data.UserRoleId,
			NameSalesZone: data.NameSalesZone,
		}
	}

	response := model.ResponseListUserByZoneStatusResponseDTO{
		Vacant:       &statusVacant,
		TotalRecords: &totalRecord,
		Records:      responseResultQuery,
	}

	ctx.SecureJSON(http.StatusOK, res.Success(response))
	return
}

func (handler *UserZoneHandler) GetChildByUser(ctx *gin.Context) {

	userId := helper.GetUserIdFromMiddleware(ctx)
	userIdString := fmt.Sprintf("%v", userId)
	roleName := ctx.Query("role_name")

	result, err, statusCode := handler.userZoneService.GetChildByUserId(userIdString, roleName)
	helpers.ResponseError(statusCode, err, *ctx)

	ctx.SecureJSON(http.StatusOK, res.Success(result))
	return
}

func (handler *UserZoneHandler) GetUserZoneByUserIDZoneIDZoneType(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}

	var body model.UserZoneRequestQueryParamDTO

	if err := ctx.ShouldBindQuery(&body); err != nil {
		helper.HandleError(err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(errorMessages[0]))
			return
		}
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}

	httpStatus, errorMessage, data := handler.userZoneService.GetUserZoneByUserIDZoneIDZoneType(&idString, body.SalesZoneType, body.SalesZoneId)
	if errorMessage != nil && httpStatus == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatus == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}
	response := responseUserZoneJoinUserDetail(data)
	ctx.SecureJSON(http.StatusOK, res.Success(response))
	return
}

func (handler *UserZoneHandler) GetChildVacantByUserId(ctx *gin.Context) {

	userId := ctx.Query("user_id")

	if userId == "" {
		userIdAuth := helper.GetUserIdFromMiddleware(ctx)
		userId = fmt.Sprintf("%v", userIdAuth)
	}

	result, err, statusCode := handler.userZoneService.GetChildVacantByUserId(userId)
	helpers.ResponseError(statusCode, err, *ctx)

	ctx.SecureJSON(http.StatusOK, res.Success(result))
	return
}

func (handler *UserZoneHandler) GetSubordinateEmployeesByUserIDZoneIDZoneType(ctx *gin.Context) {
	idString := ctx.Param("user_id")

	// just to check params is number only
	_, errConvert := helper.StringToUint64(&idString)
	helper.HandleError(errConvert)
	if errConvert != nil {
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest("type data param isn't allow"))
		return
	}
	var params model.UserZoneRequestParamsDTO

	if err := ctx.ShouldBindQuery(&params); err != nil {
		fmt.Println(err)
		helper.HandleError(err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(errorMessages[0]))
			return
		}
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}
	fmt.Print("idString: ", idString)
	salesZoneIdString := strconv.FormatUint(uint64(*params.SalesZoneId), 10)

	httpStatus, errorMessage, data := handler.userZoneService.GetSubordinateEmployeesByUserIDZoneIDZoneType(&idString, params.SalesZoneType, &salesZoneIdString, params.RoleName)
	if errorMessage != nil && httpStatus == 500 {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatus == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}

	ctx.SecureJSON(http.StatusOK, res.Success(data))
	return
}

func (handler *UserZoneHandler) GetZoneChildVacantByUserId(ctx *gin.Context) {
	userIdAuth := helper.GetUserIdFromMiddleware(ctx)
	userId := fmt.Sprintf("%v", userIdAuth)

	result, err, statusCode := handler.userZoneService.GetZoneChildVacantByUserId(userId)
	if err != nil {
		// if statusCode == http.StatusBadRequest {
		// 	ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		// 	return
		// }
		// if statusCode == http.StatusNotFound {
		// 	ctx.SecureJSON(http.StatusNotFound, res.NotFound(err.Error()))
		// 	return
		// }
		if statusCode == http.StatusOK {
			ctx.SecureJSON(http.StatusOK, res.Success(err.Error()))
			return
		}
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(err.Error()))
		return
	}

	// if result == nil {
	// 	ctx.SecureJSON(http.StatusNotFound, res.DataNotFound())
	// 	return
	// }

	ctx.SecureJSON(http.StatusOK, res.Success(result))
	return
}

func (handler *UserZoneHandler) ImpersonateAccessControlSales(ctx *gin.Context) {
	userIdAuth := helper.GetUserIdFromMiddleware(ctx)
	userId := fmt.Sprintf("%v", userIdAuth)

	result, err, statusCode := handler.userZoneService.ImpersonateAccessControlSales(userId)
	if err != nil {
		if statusCode == 400 {
			ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(err.Error()))
			return
		}
		if statusCode == 404 {
			ctx.SecureJSON(http.StatusNotFound, res.NotFound(err.Error()))
			return
		}
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(err.Error()))
		return
	}

	// if result == nil {
	// 	ctx.SecureJSON(http.StatusNotFound, res.DataNotFound())
	// 	return
	// }

	ctx.SecureJSON(http.StatusOK, res.Success(result))
	return
}

func (handler *UserZoneHandler) GetChildNonVacantByUserId(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	if userId == "" {
		userIdAuth := helper.GetUserIdFromMiddleware(ctx)
		userId = fmt.Sprintf("%v", userIdAuth)
	}

	result, err, statusCode := handler.userZoneService.GetChildNonVacantByUserId(userId)
	helpers.ResponseError(statusCode, err, *ctx)

	ctx.SecureJSON(http.StatusOK, res.Success(result))
	return
}

func (handler *UserZoneHandler) extractSliceUserResponseIntoSliceInt(data []model.UserZoneWithEpochEntity) []int {
	var newSliceInt []int
	for _, value := range data {
		newSliceInt = append(newSliceInt, int(*value.UserId))
	}
	return newSliceInt
}

func (handler *UserZoneHandler) binarySearch(data []int, search int) (result int) {
	mid := len(data) / 2
	switch {
	case len(data) == 0:
		result = -1 // not found
	case data[mid] > search:
		result = handler.binarySearch(data[:mid], search)
	case data[mid] < search:
		result = handler.binarySearch(data[mid+1:], search)
		if result >= 0 { // if anything but the -1 "not found" result
			result += mid + 1
		}
	default: // a[mid] == search
		result = mid // found
	}
	return
}

func (handler *UserZoneHandler) dryCheckGetZoneByUserID(userIDAccessToken string, params model.GetZoneByUserID) (uint, error, *model.GetZoneByUserID) {
	fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, UserID Access Token: ", userIDAccessToken)

	if params.UserId != nil {
		/*
			validation for query params user_id and to check user_id from params is Subordinate employee from user access token or not
		*/
		fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, UserID By Filter has been requested: ", *params.UserId)
		httpStatusCheckUserRoleAccessToken, errCheckUserRoleAccessToken, dataUserAccessToken := handler.userRoleService.GetById(&userIDAccessToken)
		if httpStatusCheckUserRoleAccessToken == 500 && errCheckUserRoleAccessToken != nil {
			fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, Error when get User Role By User ID Access Token")
			return 500, errCheckUserRoleAccessToken, nil
		} else if dataUserAccessToken.ID == nil || httpStatusCheckUserRoleAccessToken == 404 {
			fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, Not Found when get User Role By User ID Access Token")
			return 404, errCheckUserRoleAccessToken, nil
		}

		// get sales zone type and sales zone id default by access token
		httpStatusGetUserZoneByUserID, errGetUserZoneByUserID, dataUserZoneByUserID := handler.userZoneService.GetUserZoneByUserId(&userIDAccessToken)
		if httpStatusGetUserZoneByUserID == 500 && errGetUserZoneByUserID != nil {
			fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 500, when Get Data User Zone by User ID Access Token")
			return 500, errGetUserZoneByUserID, nil
		} else if httpStatusGetUserZoneByUserID == 404 {
			fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 404, when Get Data User Zone by User ID Access Token")
			return 404, errGetUserZoneByUserID, nil
		}

		userIDParamString := strconv.FormatUint(uint64(*params.UserId), 10)
		dataUserZoneByUserIDWithoutPointer := *dataUserZoneByUserID
		salesZoneIdString := strconv.FormatUint(uint64(*dataUserZoneByUserIDWithoutPointer.SalesZoneId), 10)

		if *dataUserAccessToken.RoleName == "nsm" {
			fmt.Println("---  USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, user_id params isn't nil, NSM")
			salesZoneType := "districts"
			roleName := "nsm"
			httpStatus, errorMessage, dataUserZone := handler.userZoneService.GetSubordinateEmployeesByUserIDZoneIDZoneType(&userIDAccessToken, &salesZoneType, &salesZoneIdString, &roleName)
			if httpStatus == 500 && errorMessage != nil {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 500 (1) - NSM, when get the data Subordinate Employees")
				return 500, errorMessage, nil
			} else if httpStatus == 404 || len(*dataUserZone) == 0 {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 404 (1) - NSM, when get the data Subordinate Employees")
				return 404, errorMessage, nil
			}
			sliceDataUserID := handler.extractSliceUserResponseIntoSliceInt(*dataUserZone)
			fmt.Println("---  USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, data user id slice, NSM: ", sliceDataUserID)
			resultCheckUserIDSubordinateEmployee := handler.binarySearch(sliceDataUserID, int(*params.UserId))
			if resultCheckUserIDSubordinateEmployee == -1 {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, Error 404 (2) - NSM, USER ID PARAMS ISN'T EXISTS")
				return 404, errors.New(fmt.Sprintf("the data User ID %v isn't exist on subordinate from User ID %v", userIDParamString, userIDAccessToken)), nil
			}
		} else if *dataUserAccessToken.RoleName == "sm" {
			fmt.Println("---  USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, user_id params isn't nil, SM")
			salesZoneType := "regions"
			roleName := "sm"
			httpStatus, errorMessage, dataUserZone := handler.userZoneService.GetSubordinateEmployeesByUserIDZoneIDZoneType(&userIDAccessToken, &salesZoneType, &salesZoneIdString, &roleName)
			if httpStatus == 500 && errorMessage != nil {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 500 (1) - SM, when get the data Subordinate Employees")
				return 500, errorMessage, nil
			} else if httpStatus == 404 || len(*dataUserZone) == 0 {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 404 (1) - SM, when get the data Subordinate Employees")
				return 404, errorMessage, nil
			}
			sliceDataUserID := handler.extractSliceUserResponseIntoSliceInt(*dataUserZone)
			fmt.Println("---  USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, data user id slice, SM: ", sliceDataUserID)
			resultCheckUserIDSubordinateEmployee := handler.binarySearch(sliceDataUserID, int(*params.UserId))
			if resultCheckUserIDSubordinateEmployee == -1 {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, Error 404 (2) - SM, USER ID PARAMS ISN'T EXISTS")
				return 404, errors.New(fmt.Sprintf("the data User ID %v isn't exist on subordinate from User ID %v", userIDParamString, userIDAccessToken)), nil
			}
		} else if *dataUserAccessToken.RoleName == "asm" {
			fmt.Println("---  USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, user_id params isn't nil, ASM")
			salesZoneType := "areas"
			roleName := "asm"
			httpStatus, errorMessage, dataUserZone := handler.userZoneService.GetSubordinateEmployeesByUserIDZoneIDZoneType(&userIDAccessToken, &salesZoneType, &salesZoneIdString, &roleName)
			if httpStatus == 500 && errorMessage != nil {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 500 (1) - ASM, when get the data Subordinate Employees")
				return 500, errorMessage, nil
			} else if httpStatus == 404 || len(*dataUserZone) == 0 {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID  ---, Error 404 (1) - ASM, when get the data Subordinate Employees")
				return 404, errorMessage, nil
			}
			sliceDataUserID := handler.extractSliceUserResponseIntoSliceInt(*dataUserZone)
			fmt.Println("---  USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, data user id slice, ASM: ", sliceDataUserID)
			resultCheckUserIDSubordinateEmployee := handler.binarySearch(sliceDataUserID, int(*params.UserId))
			if resultCheckUserIDSubordinateEmployee == -1 {
				fmt.Println("--- USER ZONE HANDLER VALIDATION / GET ZONE BY USER ID ---, Error 404 (2) - ASM, USER ID PARAMS ISN'T EXISTS")
				return 404, errors.New(fmt.Sprintf("the data User ID %v isn't exist on subordinate from User ID %v", userIDParamString, userIDAccessToken)), nil
			}
		}
	}

	return 200, nil, &params
}

func (handler *UserZoneHandler) GetZoneByUserID(ctx *gin.Context) {
	var params model.GetZoneByUserID

	if err := ctx.ShouldBindQuery(&params); err != nil {
		fmt.Println(err)
		helper.HandleError(err)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := helper.ErrorMessage(err)
			ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(errorMessages[0]))
			return
		}
		ctx.SecureJSON(http.StatusBadRequest, res.BadRequest(err.Error()))
		return
	}

	userIDAccessTokenInt := helper.GetUserIdFromMiddleware(ctx)
	userIDAccessToken := strconv.FormatUint(uint64(userIDAccessTokenInt), 10)

	httpStatusValidation, errorStatusValidation, newParams := handler.dryCheckGetZoneByUserID(userIDAccessToken, params)
	if httpStatusValidation == 500 && errorStatusValidation != nil {
		helper.HandleError(errorStatusValidation)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorStatusValidation.Error()))
		return
	} else if httpStatusValidation == 404 && errorStatusValidation != nil {
		helper.HandleError(errorStatusValidation)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorStatusValidation.Error()))
		return
	}

	var userIDString string
	if newParams.UserId != nil {
		userIDString = strconv.FormatUint(uint64(*newParams.UserId), 10)
	} else {
		userIDString = userIDAccessToken
	}

	httpStatus, errorMessage, data := handler.userZoneService.GetZoneByUserID(params.SalesZoneType, &userIDString)
	if httpStatus == 500 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusInternalServerError, res.ServerError(errorMessage.Error()))
		return
	} else if httpStatus == 404 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	} else if httpStatus == 400 && errorMessage != nil {
		helper.HandleError(errorMessage)
		ctx.SecureJSON(http.StatusNotFound, res.NotFound(errorMessage.Error()))
		return
	}

	ctx.SecureJSON(http.StatusOK, res.Success(data))
	return
}
