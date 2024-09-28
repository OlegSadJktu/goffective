package responses

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func Error(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
