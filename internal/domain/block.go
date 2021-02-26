package domain

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
