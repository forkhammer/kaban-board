package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorsResponse struct {
	Errors []string `json:"errors"`
}

type ConflictResponse struct {
	Error string `json:"error"`
	Data  any    `json:"data"`
}
