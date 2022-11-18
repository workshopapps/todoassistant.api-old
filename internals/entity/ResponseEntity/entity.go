package ResponseEntity

type ResponseMessage struct {
	Status       string `json:"status,omitempty"`
	ResponseCode int    `json:"code,omitempty"`
	Name         string `json:"name,omitempty"` //name of the error
	Message      string `json:"message,omitempty"`
	Error        any    `json:"error,omitempty"` //for errors that occur even if request is successful
	Data         any    `json:"data,omitempty"`
	Extra        any    `json:"extra,omitempty"`
}

func NewCustomError(code int, message string) *ResponseMessage {
	return &ResponseMessage{ResponseCode: code, Message: message}
}

func NewDecodingError(err error) *ResponseMessage {
	return &ResponseMessage{
		ResponseCode: 400,
		Message:      "Bad request",
		Error:        err,
	}
}

func BuildSuccessResponse(code int, message string, data any, extra ...any) *ResponseMessage {

	return &ResponseMessage{
		Status:       "success",
		ResponseCode: code,
		Name:         "",
		Message:      message,
		Error:        nil,
		Data:         data,
		Extra:        extra,
	}
}

func BuildErrorResponse(code int, message string, err interface{}, data interface{}) *ResponseMessage {

	return &ResponseMessage{
		Status:       "failure",
		ResponseCode: code,
		Name:         "",
		Message:      message,
		Error:        err,
		Data:         data,
	}
}
