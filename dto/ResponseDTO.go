package dto

type ResponseDTO struct {
	Error   bool        `json:"error"`
	Message interface{} `json:"message"`
}
