package controller

type ErrorResponse struct {
	Code int `json:"code"`
	Msg string `json:"message"`
}
