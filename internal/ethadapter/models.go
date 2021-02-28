package ethadapter

import (
	"ethproxy/internal/domain"
	"fmt"
)

type formEthRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type ethError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e ethError) Error() string {
	return fmt.Sprintf("code: %d, %s", e.Code, e.Message)
}

type responseGetBlock struct {
	Id     int64        `json:"id"`
	Error  *ethError    `json:"error"`
	Result domain.Block `json:"result"`
}

type responseBlockNumber struct {
	Id     int       `json:"id"`
	Error  *ethError `json:"error"`
	Result string    `json:"result"`
}
