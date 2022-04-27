package eth

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"
	"github.com/okex/exchain/x/evm/watcher"

	"github.com/okex/exchain/x/infura"

	"github.com/okex/exchain/x/infura/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/okex/infura-service/mysql"
	"github.com/okex/infura-service/redis"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
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

func (api *PublicAPI) GetTransactionReceipt(txHash common.Hash) (*watcher.TransactionReceipt, error) {
	receipts, err := api.orm.GetTransactionReceipt(txHash.String())
	if err != nil {
		log.Info("ERROR", err)
		return nil, err
	}
	if len(receipts) == 0 {
		return nil, nil
	}
	receipt := receipts[0]
	result := &watcher.TransactionReceipt{
		Status:            hexutil.Uint64(receipt.Status),
		CumulativeGasUsed: hexutil.Uint64(receipt.CumulativeGasUsed),
		TransactionHash:   receipt.TransactionHash,
		GasUsed:           hexutil.Uint64(receipt.GasUsed),
		BlockHash:         receipt.BlockHash,
		BlockNumber:       hexutil.Uint64(receipt.BlockNumber),
		TransactionIndex:  hexutil.Uint64(receipt.TransactionIndex),
		From:              receipt.From,
		Logs:              []*ethtypes.Log{},
	}

	var contractAddr common.Address
	if len(receipt.ContractAddress) > 0 {
		contractAddr = common.HexToAddress(receipt.ContractAddress)
		result.ContractAddress = &contractAddr
	}
	var to common.Address
	if len(receipt.To) > 0 {
		to = common.HexToAddress(receipt.To)
		result.To = &to
	}
	result.LogsBloom = defaultLogsBloom

	ethLogs := make([]*ethtypes.Log, 0) // 这里为了和以太坊eth_getLogs接口兼容，所以给用make初始化ehtLogs,为空时返回[],而不是null
	for _, v := range receipt.Logs {
		topics := make([]common.Hash, len(v.Topics))
		for i, t := range v.Topics {
			topics[i] = common.HexToHash(t.Topic)
		}
		ethLogs = append(ethLogs, &ethtypes.Log{
			Address:     common.HexToAddress(v.Address),
			Topics:      topics,
			Data:        hexutil.MustDecode(v.Data),
			BlockNumber: uint64(v.BlockNumber),
			TxHash:      common.HexToHash(v.TransactionHash),
			TxIndex:     uint(v.TransactionIndex),
			BlockHash:   common.HexToHash(v.BlockHash),
			Index:       uint(v.LogIndex),
			Removed:     false,
		})

	}
	result.Logs = ethLogs
	return result, nil
}

// GetLogs returns logs matching the given argument that are stored within the state.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getLogs
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
	ethLogs := make([]*ethtypes.Log, 0) // 这里为了和以太坊eth_getLogs接口兼容，所以给用make初始化ehtLogs,为空时返回[],而不是null
	for _, v := range transactionLogs {
		topics := make([]common.Hash, len(v.Topics))
		for i, t := range v.Topics {
			topics[i] = common.HexToHash(t.Topic)
		}
		// match topics
		if len(criteria.Topics) > 0 && !matchTopics(topics, criteria.Topics) {
			continue
		}
		ethLogs = append(ethLogs, &ethtypes.Log{
			Address:     common.HexToAddress(v.Address),
			Topics:      topics,
			Data:        hexutil.MustDecode(v.Data),
			BlockNumber: uint64(v.BlockNumber),
			TxHash:      common.HexToHash(v.TransactionHash),
			TxIndex:     uint(v.TransactionIndex),
			BlockHash:   common.HexToHash(v.BlockHash),
			Index:       uint(v.LogIndex),
			Removed:     false,
		})

	}
	return ethLogs, nil
}

// The Topic list restricts matches to particular event topics. Each event has a list
// of topics. Topics matches a prefix of that list. An empty element slice matches any
// topic. Non-empty elements represent an alternative that matches any of the
// contained topics.
//
// Examples:
// {} or nil          matches any topic list
// {{A}}              matches topic A in first position
// {{}, {B}}          matches any topic in first position AND B in second position
// {{A}, {B}}         matches topic A in first position AND B in second position
// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
func matchTopics(topics []common.Hash, matches [][]common.Hash) bool {
	matchCount := len(matches)
	// 处理传入的参数，topic末位写入null的情况，比如传入：{{A}, nil, nil, nil}，这种情况只要第一个A满足条件，后面的nil忽略即可
	for i := len(matches) - 1; i >= 0; i-- {
		if len(matches[i]) > 0 {
			break
		}
		matchCount--
	}
	// 要求的topic数量不匹配
	if matchCount > len(topics) {
		return false
	}
	// 验证topic
	for i := 0; i < matchCount; i++ {
		if len(matches[i]) == 0 {
			continue
		}
		isMatch := false
		for _, match := range matches[i] {
			if bytes.Equal(topics[i].Bytes(), match.Bytes()) {
				isMatch = true
				break
			}
		}
		if !isMatch {
			return false
		}
	}
	return true
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
