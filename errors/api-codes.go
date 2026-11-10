package errors

type ApiCode string

const (
	ApiInternalServerError ApiCode = "E_INTERNAL_SERVER_ERROR"
	ApiUnauthorized        ApiCode = "E_UNAUTHORIZED"
	ApiForbidden           ApiCode = "E_FORBIDDEN"
	ApiNotFound            ApiCode = "E_NOT_FOUND"
	ApiResourceLocked      ApiCode = "E_RESOURCE_LOCKED"
	ApiServiceUnavailable  ApiCode = "E_SERVICE_UNAVAILABLE"
)
