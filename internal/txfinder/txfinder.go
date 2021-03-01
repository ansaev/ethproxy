package txfinder

import (
	"errors"
	"ethproxy/internal/domain"
	"fmt"
	"strconv"
)

const (
	latest = "latest"
)

var (
	ErrTxNotFoundInBlock = errors.New("cannot find tx in block")
	// nolint:gosimple // always use errors.New for producing new errors
	ErrInvalidBlockID = errors.New(fmt.Sprintf("blockID must be num or \"%s\"", latest))
)

type BlockGetter interface {
	GetBlockByNum(blockID uint64) (*domain.Block, error)
	GetLatestBlockNum() (uint64, error)
}

type service struct {
	blockService BlockGetter
}

func New(blockService BlockGetter) *service {
	return &service{
		blockService: blockService,
	}
}

func (srv *service) FindTx(blockID string, txID string) (*domain.Transaction, error) {
	// get block by id
	var blockNum uint64
	// get block id
	if blockID != latest {
		var errParseBlockID error
		blockNum, errParseBlockID = strconv.ParseUint(blockID, 10, 64)
		if errParseBlockID != nil {
			return nil, ErrInvalidBlockID
		}
	} else {
		var errGetLatestBlockNum error
		blockNum, errGetLatestBlockNum = srv.blockService.GetLatestBlockNum()
		if errGetLatestBlockNum != nil {
			return nil, fmt.Errorf("failed to get latest block num: %w", errGetLatestBlockNum)
		}
	}
	// get block
	block, err := srv.blockService.GetBlockByNum(blockNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}

	// parse tx ID
	txNum, errParseTxNum := strconv.ParseUint(txID, 10, 64)
	if errParseTxNum == nil {
		// search tx by index
		return srv.findTxByIndex(block, txNum)
	}

	// search tx by hash
	tx, errFindTxByHash := srv.findTxByHash(block, txID)
	if errFindTxByHash != nil {
		return nil, fmt.Errorf("failed to get tx by hash: %w. "+
			"tx_id must be index of transaction or transaction's hash", errFindTxByHash)
	}

	return tx, nil
}

func (srv *service) findTxByHash(block *domain.Block, txHash string) (*domain.Transaction, error) {
	for i := range block.Transactions {
		if block.Transactions[i].Hash == txHash {
			return &block.Transactions[i], nil
		}
	}

	return nil, fmt.Errorf("%w, blockHash: %s, txHash: %s", ErrTxNotFoundInBlock, block.Hash, txHash)
}

func (srv *service) findTxByIndex(block *domain.Block, txIndex uint64) (*domain.Transaction, error) {
	for i := range block.Transactions {
		if block.Transactions[i].GetTxIndex() == txIndex {
			return &block.Transactions[i], nil
		}
	}

	return nil, fmt.Errorf("%w, blockHash: %s, txIndex: %d", ErrTxNotFoundInBlock, block.Hash, txIndex)
}
