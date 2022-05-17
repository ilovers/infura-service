package eth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	evmtypes "github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/infura"
	"github.com/okex/exchain/x/infura/types"
	"github.com/okex/infura-service/mysql"
	"github.com/okex/infura-service/redis"
)

type PublicAPI struct {
	orm      *mysql.Orm
	redisCli *redis.Client
}

func NewAPI(orm *mysql.Orm, redisCli *redis.Client) (*PublicAPI, error) {
	return &PublicAPI{
		orm:      orm,
		redisCli: redisCli,
	}, nil
}

// GetTransactionReceipt handles eth_getTransactionReceipt
func (api *PublicAPI) GetTransactionReceipt(txHash common.Hash) (*evmtypes.TransactionReceipt, error) {
	receipts, err := api.orm.GetTransactionReceipt(txHash.String())
	if err != nil {
		log.Info("ERROR", err)
		return nil, err
	}
	if len(receipts) == 0 {
		return nil, nil
	}
	receipt := receipts[0]
	result := convertTransactionReceipt(receipt)
	return result, nil
}

// GetLogs returns logs matching the given argument that are stored within the state.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getLogs
// GetLogs handles eth_getLogs
func (api *PublicAPI) GetLogs(ctx context.Context, criteria filters.FilterCriteria) ([]*ethtypes.Log, error) {
	var transactionLogs []types.TransactionLog
	var err error
	// contract address
	addresses := make([]string, len(criteria.Addresses))
	for i, addr := range criteria.Addresses {
		addresses[i] = addr.String()
	}
	// 从mysql查询数据，分两种情况，一种是使用blockHash，另外一种是使用blockNum
	if criteria.BlockHash != nil {
		transactionLogs, err = api.orm.GetLogsByBlockHash(criteria.BlockHash.String(), addresses)
		if err != nil {
			log.Info("ERROR", err)
			return nil, err
		}
	} else {
		var fromBlock, toBlock int64
		if criteria.FromBlock != nil {
			fromBlock = criteria.FromBlock.Int64()
		} else {
			fromBlock = api.latestBlock()
		}
		if criteria.ToBlock != nil {
			toBlock = criteria.ToBlock.Int64()
		} else {
			toBlock = fromBlock
		}

		transactionLogs, err = api.orm.GetLogs(fromBlock, toBlock, addresses)
		if err != nil {
			log.Info("ERROR", err)
			return nil, err
		}
	}
	ethLogs := convertLogs(transactionLogs, criteria.Topics)
	return ethLogs, nil
}

func (api *PublicAPI) latestBlock() int64 {
	value, err := api.redisCli.Get(latestTaskKey)
	if err != nil {
		return 0
	}
	task := infura.Task{}
	err = json.Unmarshal([]byte(value), &task)
	if err != nil {
		return 0
	}
	return task.Height
}

func (api *PublicAPI) GetBlockByNumber(blockNum rpc.BlockNumber, fullTx bool) (*evmtypes.Block, error) {
	height := int64(blockNum)
	if height <= 0 {
		height = api.latestBlock()
	}
	block, err := api.orm.GetBlockByNumber(height)
	if err != nil {
		return nil, nil
	}
	evmBlock := convertBlock(block, fullTx)
	return evmBlock, nil
}

func (api *PublicAPI) GetBlockByHash(blockHash common.Hash, fullTx bool) (*evmtypes.Block, error) {
	block, err := api.orm.GetBlockByHash(blockHash.String())
	if err != nil {
		return nil, err
	}
	evmBlock := convertBlock(block, fullTx)
	return evmBlock, nil
}

func (api *PublicAPI) GetBlockTransactionCountByNumber(blockNum rpc.BlockNumber) *hexutil.Uint {
	height := int64(blockNum)
	if height <= 0 {
		height = api.latestBlock()
	}
	block, err := api.orm.GetBlockByNumber(height)
	if err != nil {
		return nil
	}
	n := hexutil.Uint(len(block.Transactions))
	return &n
}

func (api *PublicAPI) GetBlockTransactionCountByHash(blockHash common.Hash) *hexutil.Uint {
	block, err := api.orm.GetBlockByHash(blockHash.String())
	if err != nil {
		return nil
	}
	n := hexutil.Uint(len(block.Transactions))
	return &n
}

func (api *PublicAPI) GetTransactionByBlockHashAndIndex(blockHash common.Hash, idx hexutil.Uint) (*evmtypes.Transaction, error) {
	block, err := api.orm.GetBlockByHash(blockHash.String())
	if err != nil {
		return nil, nil
	}
	var transaction *evmtypes.Transaction
	for _, t := range block.Transactions {
		if t.Index == uint64(idx) {
			evmTransaction := convertTransaction(t, block.Number, block.Hash)
			transaction = &evmTransaction
		}
	}
	return transaction, nil
}

func (api *PublicAPI) GetTransactionByBlockNumberAndIndex(blockNum rpc.BlockNumber, idx hexutil.Uint) (*evmtypes.Transaction, error) {
	height := int64(blockNum)
	if height <= 0 {
		height = api.latestBlock()
	}
	block, err := api.orm.GetBlockByNumber(height)
	if err != nil {
		return nil, nil
	}
	var transaction *evmtypes.Transaction
	for _, t := range block.Transactions {
		if t.Index == uint64(idx) {
			evmTransaction := convertTransaction(t, block.Number, block.Hash)
			transaction = &evmTransaction
		}
	}
	return transaction, nil
}

func (api *PublicAPI) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {
	receipts, err := api.orm.GetTransactionReceipt(txHash.String())
	if err != nil {
		log.Info("ERROR", err)
		return nil, err
	}
	if len(receipts) == 0 {
		return nil, errors.New(fmt.Sprintf("tx (%s) not found", txHash.String()))
	}
	receipt := receipts[0]
	result := convertLogs(receipt.Logs, nil)
	if len(result) == 0 { // 为空时返回null,不返回[]
		return nil, nil
	}
	return result, nil
}

func (api *PublicAPI) GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	blockNumber, err := api.convertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	contractCode, err := api.orm.GetContractCode(address.String())
	if err != nil || (blockNumber > 0 && contractCode.BlockNumber > blockNumber) {
		return nil, nil // 没有查询结果时返回nil，不返回错误
	}
	return hexutil.MustDecode(contractCode.Code), nil
}

// 参考以太坊源码，返回error信息和以太坊保持一致
func (api *PublicAPI) convertToBlockNumber(blockNrOrHash rpc.BlockNumberOrHash) (int64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return int64(blockNr), nil
	}

	if hash, ok := blockNrOrHash.Hash(); ok {
		block, err := api.orm.GetBlockByHash(hash.String())
		if err != nil {
			return 0, errors.New("header for hash not found")
		}
		return block.Number, nil
	}
	return 0, errors.New("invalid arguments; neither block nor hash specified")
}
