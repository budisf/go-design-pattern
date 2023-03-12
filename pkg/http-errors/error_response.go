package error

import res "ethical-be/pkg/api-response"

func PageNotFound() res.Meta {
	result := res.Meta{
		Message: "404 page not found",
		Status:  404,
		Code:    "PAGE_NOT_FOUND",
	}
	return result
}

func MethodNotAllowed() res.Meta {
	result := res.Meta{
		Message: "404 page not found",
		Status:  404,
		Code:    "PAGE_NOT_FOUND",
	}
	return result
}
