package api

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  *Error `json:"error,omitempty"` // pointer for omitempty in case of nil
}

type Error struct {
	Message string      `json:"message"`
	Details []ErrDetail `json:"details,omitempty"`
}

type ErrDetail struct {
	Field string `json:"field,omitempty"`
	Info  string `json:"info"`
}

func Ok() Response {
	return Response{
		Status: StatusOk,
		Error:  nil,
	}
}

func Err(msg string) Response {
	return Response{
		Status: StatusError,
		Error: &Error{
			Message: msg,
			Details: nil,
		},
	}
}

func ErrD(msg string, details []ErrDetail) Response {
	return Response{
		Status: StatusError,
		Error: &Error{
			Message: msg,
			Details: details,
		},
	}
}
