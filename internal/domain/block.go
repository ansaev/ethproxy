package domain

import (
	"strconv"
	"time"
)

type Block struct {
	Difficulty   string        `json:"difficulty"`
	ExtraData    string        `json:"extraData"`
	GasLimit     string        `json:"gasLimit"`
	GasUsed      string        `json:"gasUsed"`
	Hash         string        `json:"hash"`
	MixHash      string        `json:"mixHash"`
	Nonce        string        `json:"Nonce"`
	Number       string        `json:"number"`
	ParentHash   string        `json:"parentHash"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

func (b *Block) GetTime() time.Time {
	timestamp, _ := strconv.ParseInt(b.Timestamp, 0, 64)
	return time.Unix(timestamp, 0)
}

func (b *Block) GetNumber() uint64 {
	num, _ := strconv.ParseUint(b.Timestamp, 0, 64)
	return num
}
