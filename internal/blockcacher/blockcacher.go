package blockcacher

import (
	"encoding/json"
	"ethproxy/internal/domain"
	"fmt"
	"log"
	"time"
)

type blockGetter interface {
	GetBlockByNum(blockID uint64) (*domain.Block, error)
	GetLatestBlockNum() (uint64, error)
	IsBlockCacheable(block *domain.Block) bool
}

type cacheService interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte, expires time.Duration) error
	Del(key string) error
}

type service struct {
	blockService   blockGetter
	cache          cacheService
	blockchainName string
	cachingTime    time.Duration
	isDebug        bool
}

func New(blockService blockGetter, cache cacheService, blockchainName string, cachingTime time.Duration,
	isDebug bool) *service {
	return &service{
		blockService:   blockService,
		cache:          cache,
		blockchainName: blockchainName,
		cachingTime:    cachingTime,
		isDebug:        isDebug,
	}
}

func (s *service) getCacheKeyForBlockNum(blockID uint64) string {
	return fmt.Sprintf("%s_block_num_%d", s.blockchainName, blockID)
}

func (s *service) GetBlockByNum(blockID uint64) (*domain.Block, error) {
	// try get from cash
	cacheKey := s.getCacheKeyForBlockNum(blockID)
	blockData, errGetFromCashe := s.cache.Get(cacheKey)
	if errGetFromCashe == nil {
		block := &domain.Block{}
		errUnmarshalCashedData := json.Unmarshal(blockData, block)
		if errUnmarshalCashedData == nil {
			// return block from cached
			if s.isDebug {
				log.Printf("return block:%d from cache\n", blockID)
			}
			return block, nil
		}
		errDelKey := s.cache.Del(cacheKey)
		if errDelKey != nil {
			log.Printf("failed to delete key \"%s\" from cache: %v\n", cacheKey, errDelKey)
		}
	} else {
		log.Printf("failed to get block %d from cached: %v\n", blockID, errGetFromCashe) // debug comment
	}

	// get block from network
	block, errGetBlock := s.blockService.GetBlockByNum(blockID)
	if errGetBlock != nil {
		return nil, errGetBlock
	}

	if s.blockService.IsBlockCacheable(block) {
		// save block to cache
		blockData, errMarshal := json.Marshal(block)
		if errMarshal != nil {
			log.Printf("failed to marshal block %d: %v", block.GetNumber(), errMarshal)
			return block, nil
		}
		errSetCache := s.cache.Set(cacheKey, blockData, s.cachingTime)
		if errSetCache != nil {
			log.Printf("failed to set cached block %d: %v", block.GetNumber(), errSetCache)
			return block, nil
		}
	}

	if s.isDebug {
		log.Printf("return block:%d from network\n", blockID) // debug comment
	}
	return block, nil
}

func (s *service) GetLatestBlockNum() (uint64, error) {
	// may cashed this too
	return s.blockService.GetLatestBlockNum()
}
