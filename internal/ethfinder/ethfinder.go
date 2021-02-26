package ethfinder

import (
	"bytes"
	"encoding/json"
	"errors"
	"ethproxy/internal/domain"
	"fmt"
	"net/http"
	"strconv"
)

const (
	methodGetBlockByNum = "eth_getBlockByNumber"
	methodBlockNum      = "eth_blockNumber"
	jsonRPCVersion      = "2.0"
	headerContentType   = "Content-Type"
	applicationJson     = "application/json"
)

var (
	errBadStatus = errors.New("get response with bad status")
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func New(address string, client httpClient) *service {
	return &service{
		address: address,
		client:  client,
	}
}

type service struct {
	address string
	client  httpClient
}

func ethHexFromInt(num int64) string {
	return fmt.Sprintf("0x%x", num)
}

func (srv *service) GetLatestBlockNum() (int64, error) {
	// prepare data for request
	reqData := &formEthRequest{
		Jsonrpc: jsonRPCVersion,
		Method:  methodBlockNum,
		Id:      1,
		Params:  []interface{}{},
	}
	result := &responseBlockNumber{}
	err := srv.doReq(reqData, result)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block num: %w", err)
	}
	if result.Error != nil {
		return 0, fmt.Errorf("got response with error: %w", result.Error)
	}

	// parse result
	blockNum, errParseBlockNum := strconv.ParseInt(result.Result, 0, 64)
	if errParseBlockNum != nil {
		return 0, fmt.Errorf("failed to parse block num from hex: %w, hex block number: %s",
			errParseBlockNum, result.Result)
	}
	return blockNum, nil
}

func (srv *service) GetBlockByNum(blockID int64) (*domain.Block, error) {
	reqData := &formEthRequest{
		Jsonrpc: jsonRPCVersion,
		Method:  methodGetBlockByNum,
		Id:      1,
		Params:  []interface{}{ethHexFromInt(blockID), true},
	}
	blockResp := &responseGetBlock{}

	err := srv.doReq(reqData, blockResp)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by num: %w", err)
	}

	if blockResp.Error != nil {
		return nil, fmt.Errorf("got response with error: %w", blockResp.Error)
	}

	return &blockResp.Result, nil
}

func (srv *service) doReq(reqData *formEthRequest, result interface{}) error {
	// prepare data for request
	jsonData, errMarshalReq := json.Marshal(reqData)
	if errMarshalReq != nil {
		return fmt.Errorf("failed to encode reqest data: %w", errMarshalReq)
	}

	// create http request
	req, errCreateReq := http.NewRequest(http.MethodPost, srv.address, bytes.NewBuffer(jsonData))
	if errCreateReq != nil {
		return fmt.Errorf("failed to create request: %w", errCreateReq)
	}
	req.Header.Add(headerContentType, applicationJson)

	// do req
	resp, errDoReq := srv.client.Do(req)
	if errDoReq != nil {
		return fmt.Errorf("failed to execute request: %w", errDoReq)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w, status: %d", errBadStatus, resp.StatusCode)
	}

	// parse response
	errDecodeResp := json.NewDecoder(resp.Body).Decode(result)
	if errDecodeResp != nil {
		return fmt.Errorf("failed to decode response: %w", errDecodeResp)
	}

	return nil
}
