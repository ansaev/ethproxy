package blockcacher

import (
	"testing"
	"time"
)

const (
	blockHash = "0xbbf25a254eb78aab5109c10b8e005423beb5bad5361ab29920de9d0ff82e40be"
)

func TestGetBlockByNum_blockNotInCacheAndNotCacheble(t *testing.T) {
	cacheService := &mockEmptyCacheService{emptyGet: true, calledSet: false}
	blockService := &mockBlockService{
		CalledIsBlockCacheable: false,
		isBlockCacheableResult: false,
		GetBlockByNumCall:      0,
	}
	blockID := uint64(123)

	cacheBlockService := New(blockService, cacheService, "test", time.Second, true)
	block, err := cacheBlockService.GetBlockByNum(blockID)
	if err != nil {
		t.Fatal("failed to get block via cash: ", err)
	}

	// check if cacheBlockService  tried to set cache
	if cacheService.calledSet {
		t.Fatal("blockcacher tried to set cache, but it shouldn't have")
	}
	if !blockService.CalledIsBlockCacheable {
		t.Fatal("blockcacher not call IsBlockCacheable, but it should have")
	}
	if blockService.GetBlockByNumCall != blockID {
		t.Fatalf("blockcacher expected to call GetBlockByNum with block's number: \"%d\", but it is: \"%d\"\n",
			blockID, blockService.GetBlockByNumCall)
	}

	if block.Hash != blockHash {
		t.Fatalf("block's hash expected to be \"%s\", but it is \"%s\"\n", blockHash, block.Hash)
	}
}

func TestGetBlockByNum_blockNotInCacheAndCacheble(t *testing.T) {
	cacheService := &mockEmptyCacheService{emptyGet: true, calledSet: false}
	blockService := &mockBlockService{
		CalledIsBlockCacheable: false,
		isBlockCacheableResult: true,
		GetBlockByNumCall:      0,
	}
	blockID := uint64(123)

	cacheBlockService := New(blockService, cacheService, "test", time.Second, true)
	block, err := cacheBlockService.GetBlockByNum(blockID)
	if err != nil {
		t.Fatal("failed to get block via cash: ", err)
	}

	// check if cacheBlockService  tried to set cache
	if !cacheService.calledSet {
		t.Fatal("blockcacher not tried to set cache, but it should have")
	}
	if !blockService.CalledIsBlockCacheable {
		t.Fatal("blockcacher not call IsBlockCacheable, but it should have")
	}
	if blockService.GetBlockByNumCall != blockID {
		t.Fatalf("blockcacher expected to call GetBlockByNum with block's number: \"%d\", but it is: \"%d\"\n",
			blockID, blockService.GetBlockByNumCall)
	}

	if block.Hash != blockHash {
		t.Fatalf("block's hash expected to be \"%s\", but it is \"%s\"\n", blockHash, block.Hash)
	}
}

func TestGetBlockByNum_blockInCacheAndCacheble(t *testing.T) {
	cacheService := &mockEmptyCacheService{emptyGet: false, calledSet: false}
	blockService := &mockBlockService{
		CalledIsBlockCacheable: false,
		isBlockCacheableResult: true,
		GetBlockByNumCall:      0,
	}
	blockID := uint64(123)

	cacheBlockService := New(blockService, cacheService, "test", time.Second, true)
	block, err := cacheBlockService.GetBlockByNum(blockID)
	if err != nil {
		t.Fatal("failed to get block via cash: ", err)
	}

	// check if cacheBlockService  tried to set cache
	if cacheService.calledSet {
		t.Fatal("blockcacher tried to set cache, but it shouldn't have")
	}
	if blockService.CalledIsBlockCacheable {
		t.Fatal("blockcacher called IsBlockCacheable, but it should't have")
	}
	if blockService.GetBlockByNumCall != 0 {
		t.Fatalf(
			"blockcachercalled GetBlockByNum with block's number: \"%d\", but it should't have\n",
			blockService.GetBlockByNumCall)
	}

	if block.Hash != blockHash {
		t.Fatalf("block's hash expected to be \"%s\", but it is \"%s\"\n", blockHash, block.Hash)
	}
}
