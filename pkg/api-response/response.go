package apiresponse

import (
	"net/http"
	"strings"
)

func BadRequest(errorMessages interface{}) Meta {
	result := Meta{
		Message: errorMessages,
		Status:  400,
		Code:    "BAD_REQUEST",
	}
	return result
}

func StatusOK(message string) Meta {
	result := Meta{
		Message: message,
		Status:  200,
		Code:    strings.ToUpper(strings.Replace(message, " ", "_", -1)),
	}
	return result
}

func Success(data interface{}) Response {
	result := Response{
		Message: "Success",
		Status:  200,
		Code:    "SUCCESS",
		Data:    data,
	}
	return result
}

func DataNotFound() Meta {
	result := Meta{
		Message: "Data Not Found",
		Status:  http.StatusNotFound,
		Code:    "DATA_NOT_FOUND",
	}
	return result
}

func NotFound(message string) Meta {
	result := Meta{
		Message: message + " Not Found",
		Status:  http.StatusNotFound,
		Code:    strings.ToUpper(strings.Replace(message, " ", "_", -1)) + "_NOT_FOUND",
	}
	return result
}

func StatusForbidden(message string) Meta {
	result := Meta{
		Message: "Forbidden, " + message,
		Status:  http.StatusForbidden,
		Code:    "FORBIDDEN" + strings.ToUpper(strings.Replace(message, " ", "_", -1)),
	}
	return result
}

func DataFound(data string) Meta {
	result := Meta{
		Message: data + " Exist",
		Status:  http.StatusFound,
		Code:    strings.ToUpper(data) + "_EXIST",
	}
	return result
}

func ServerError(errorMessages interface{}) Meta {
	result := Meta{
		Message: errorMessages,
		Status:  http.StatusInternalServerError,
		Code:    "STATUS_INTERNAL_SERVER_ERROR",
	}
	return result
}

func UnAuthorized(errorMessage string) Meta {
	var result Meta
	switch errorMessage {
	case "Token is expired":
		result = Meta{
			Message: errorMessage,
			Status:  http.StatusUnauthorized,
			Code:    "TOKEN_IS_EXPIRED",
		}
	case "Invalid token":
		result = Meta{
			Message: errorMessage,
			Status:  http.StatusUnauthorized,
			Code:    "INVALID_TOKEN",
		}
	default:
		result = Meta{
			Message: errorMessage,
			Status:  http.StatusUnauthorized,
			Code:    "UNAUTHORIZED",
		}
	}
	return result

}

// @NOTED: empty object is used when data doesnt want to be null on json
type EmptyObj struct{}
