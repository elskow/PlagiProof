package util

type Response struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type EmptyObj struct{}

func ResponseSuccess(message string, data any) Response {
	res := Response{
		Ok:      true,
		Message: message,
		Data:    data,
	}
	return res
}

func ResponseFailed(message string, err string, data any) Response {
	res := Response{
		Ok:      false,
		Message: message,
		Error:   err,
		Data:    data,
	}
	return res
}
