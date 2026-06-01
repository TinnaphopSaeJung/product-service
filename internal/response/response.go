package response

const (
	CodeValidationError     = "VALIDATION_ERROR"
	CodeInvalidProductID    = "INVALID_PRODUCT_ID"
	CodeProductNotFound     = "PRODUCT_NOT_FOUND"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

type Response struct {
	Successful bool        `json:"successful"`
	ErrorCode  *string     `json:"error_code"`
	Data       interface{} `json:"data"`
}

type ResponseNoData struct {
	Successful bool    `json:"successful"`
	ErrorCode  *string `json:"error_code"`
}

func Success(data interface{}) Response {
	return Response{
		Successful: true,
		ErrorCode:  nil,
		Data:       data,
	}
}

func Error(errorCode string) Response {
	return Response{
		Successful: false,
		ErrorCode:  &errorCode,
		Data:       nil,
	}
}

func SuccessNoData() ResponseNoData {
	return ResponseNoData{
		Successful: true,
		ErrorCode:  nil,
	}
}

func ErrorNoData(errorCode string) ResponseNoData {
	return ResponseNoData{
		Successful: false,
		ErrorCode:  &errorCode,
	}
}
