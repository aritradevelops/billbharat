package handlers

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

func NewResponse(message string, data any, error any) Response {
	return Response{
		Message: message,
		Data:    data,
		Error:   error,
	}
}
