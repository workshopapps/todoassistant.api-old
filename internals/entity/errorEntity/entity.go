package errorEntity

type ErrorRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewCustomError(code int, message string) *ErrorRes {
	return &ErrorRes{Code: code, Message: message}
}

func NewDecodingError() *ErrorRes {
	return &ErrorRes{
		Code:    400,
		Message: "Bad request",
	}
}
