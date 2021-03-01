package api

import "ethproxy/internal/domain"

type TxResponse struct {
	Error *ErrorForm          `json:"error"`
	Data  *domain.Transaction `json:"data"`
	Ok    bool                `json:"ok"`
}

type ErrorForm struct {
	Message string
	// TODO: add internal error's code
}
