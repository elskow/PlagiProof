package dto

type ResponseDTO struct {
	Status  int         `json:"status"`
	Error   bool        `json:"error"`
	Message interface{} `json:"message"`
}
