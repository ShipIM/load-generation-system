package web

type status string

const (
	OK    = status("OK")
	ERROR = status("ERROR")
)

type (
	Response struct {
		Status status `json:"status"`
		Data   any    `json:"data,omitempty"`
	}
)

func OKResponse(body any) Response {
	return Response{
		Status: OK,
		Data:   body,
	}
}

func ErrorResponse(body any) Response {
	return Response{
		Status: ERROR,
		Data:   body,
	}
}
