package httpresponse

type Success struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Error struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

func NewSuccess(message string, data any) Success {
	return Success{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewError(message string, errors any) Error {
	return Error{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}