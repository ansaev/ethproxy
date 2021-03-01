package blockcacher

import (
	"encoding/json"
	"errors"
	"ethproxy/internal/domain"
	"io/ioutil"
	"time"
)

const (
	pathData             = "internal/blockcacher/data/"
	pathBlockFile        = "block.json"
	lastBlockNum  uint64 = 345
)

type mockBlockService struct {
	CalledIsBlockCacheable bool
	isBlockCacheableResult bool
	GetBlockByNumCall      uint64
}

func (bs *mockBlockService) GetBlockByNum(blockID uint64) (*domain.Block, error) {
	bs.GetBlockByNumCall = blockID
	data, errReadFile := ioutil.ReadFile(pathData + pathBlockFile)
	if errReadFile != nil {
		return nil, errReadFile
	}

	block := &domain.Block{}
	errUnmarshal := json.Unmarshal(data, block)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return block, nil
}

func (bs *mockBlockService) GetLatestBlockNum() (uint64, error) {
	return lastBlockNum, nil
}

func (bs *mockBlockService) IsBlockCacheable(block *domain.Block) bool {
	bs.CalledIsBlockCacheable = true
	return bs.isBlockCacheableResult
}

type mockEmptyCacheService struct {
	emptyGet  bool
	calledSet bool
}

func (cache *mockEmptyCacheService) Get(key string) ([]byte, error) {
	if cache.emptyGet {
		return nil, errors.New("no such key")
	}

	data, errReadFile := ioutil.ReadFile(pathData + pathBlockFile)
	if errReadFile != nil {
		return nil, errReadFile
	}

	return data, nil
}

func (cache *mockEmptyCacheService) Set(key string, data []byte, expires time.Duration) error {
	cache.calledSet = true
	return nil
}

func (cache *mockEmptyCacheService) Del(key string) error {
	return nil
}
