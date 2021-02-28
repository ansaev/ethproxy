package cachefinder

import (
	"encoding/json"
	"ethproxy/internal/domain"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type blockGetter interface {
	GetBlockByNum(blockID uint64) (*domain.Block, error)
	GetLatestBlockNum() (uint64, error)
	IsBlockCacheable(block *domain.Block) bool
}

type service struct {
	blockService blockGetter
	// todo: move it to inteface casher
	cache          *redis.Client
	blockchainName string
	cachingTime    time.Duration
}

func New(blockService blockGetter, cache *redis.Client, blockchainName string) *service {
	return &service{
		blockService:   blockService,
		cache:          cache,
		blockchainName: blockchainName,
	}
}

func (s *service) getCacheKeyForBlockNum(blockID uint64) string {
	return fmt.Sprintf("%s_block_num_%d", s.blockchainName, blockID)
}

func (s *service) GetBlockByNum(blockID uint64) (*domain.Block, error) {
	// try get from cash
	cacheKey := s.getCacheKeyForBlockNum(blockID)
	blockData, errGetFromCashe := s.cache.Get(cacheKey).Bytes()
	if errGetFromCashe == nil {
		block := &domain.Block{}
		errUnmarshalCashedData := json.Unmarshal(blockData, block)
		if errUnmarshalCashedData == nil {
			// return block from cached
			return block, nil
		}
		s.cache.Del(cacheKey)
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

	return block, nil
}

func (s *service) GetLatestBlockNum() (uint64, error) {
	// may cashed this too
	return s.blockService.GetLatestBlockNum()
}
