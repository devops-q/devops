package json_models

type ErrorResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error_msg"`
}
