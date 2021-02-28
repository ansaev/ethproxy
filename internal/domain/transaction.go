package domain

import "strconv"

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumBer      string `json:"blockNumber"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	Value            string `json:"value"`
}

func (t *Transaction) GetTxIndex() uint64 {
	i, _ := strconv.ParseUint(t.TransactionIndex, 0, 64)
	return i
}
