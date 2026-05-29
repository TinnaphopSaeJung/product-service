package response

const (
	CodeValidationError     = "VALIDATION_ERROR"
	CodeInvalidProductID    = "INVALID_PRODUCT_ID"
	CodeProductNotFound     = "PRODUCT_NOT_FOUND"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

type Response struct {
	Successful bool        `json:"successful"`
	ErrorCode  string      `json:"error_code"`
	Data       interface{} `json:"data"`
}

func Success(data interface{}) Response {
	return Response{
		Successful: true,
		ErrorCode:  "",
		Data:       data,
	}
}

func Error(errorCode string) Response {
	return Response{
		Successful: false,
		ErrorCode:  errorCode,
		Data:       nil,
	}
}
