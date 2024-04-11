package http

type ErrorResponse struct {
	Message string `json:"message"`
}

type SendResponse struct {
	URLsCodes map[string]int `json:"urls_codes"`
}
