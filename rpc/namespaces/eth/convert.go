package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/infura/types"
)

func convertBlock(block types.Block, fullTx bool) *evmtypes.Block {
	evmBlock := &evmtypes.Block{
		Number:           hexutil.Uint64(block.Number),
		Hash:             common.HexToHash(block.Hash),
		ParentHash:       common.HexToHash(block.ParentHash),
		Nonce:            evmtypes.BlockNonce{},
		UncleHash:        ethtypes.EmptyUncleHash,
		LogsBloom:        defaultLogsBloom,
		TransactionsRoot: common.HexToHash(block.TransactionsRoot),
		StateRoot:        common.HexToHash(block.StateRoot),
		Miner:            common.HexToAddress(block.Miner),
		MixHash:          common.Hash{},
		Difficulty:       hexutil.Uint64(0),
		TotalDifficulty:  hexutil.Uint64(0),
		ExtraData:        hexutil.Bytes{},
		Size:             hexutil.Uint64(block.Size),
		GasLimit:         hexutil.Uint64(block.GasLimit),
		Timestamp:        hexutil.Uint64(block.Timestamp),
		Uncles:           []common.Hash{},
		ReceiptsRoot:     ethtypes.EmptyRootHash,
	}
	gasUsed := hexutil.Big(*big.NewInt(int64(block.GasUsed)))
	evmBlock.GasUsed = &gasUsed

	if fullTx {
		transactions := make([]evmtypes.Transaction, len(block.Transactions))
		for i, t := range block.Transactions {
			transactions[i] = convertTransaction(t, block.Number, block.Hash)
		}
		evmBlock.Transactions = transactions
	} else {
		transactions := make([]common.Hash, len(block.Transactions))
		for i, t := range block.Transactions {
			transactions[i] = common.HexToHash(t.Hash)
		}
		evmBlock.Transactions = transactions
	}
	return evmBlock
}

func convertTransaction(t types.Transaction, blockNumber int64, blockHash string) evmtypes.Transaction {
	number := hexutil.Big(*big.NewInt(blockNumber))
	hash := common.HexToHash(blockHash)
	gasPrice := hexutil.Big(*hexutil.MustDecodeBig(t.GasPrice))
	index := hexutil.Uint64(t.Index)
	value := hexutil.Big(*hexutil.MustDecodeBig(t.Value))
	V := hexutil.Big(*hexutil.MustDecodeBig(t.V))
	R := hexutil.Big(*hexutil.MustDecodeBig(t.R))
	S := hexutil.Big(*hexutil.MustDecodeBig(t.S))
	result := evmtypes.Transaction{
		BlockHash:        &hash,
		BlockNumber:      &number,
		From:             common.HexToAddress(t.From),
		Gas:              hexutil.Uint64(t.Gas),
		GasPrice:         &gasPrice,
		Hash:             common.HexToHash(t.Hash),
		Input:            hexutil.MustDecode(t.Input),
		Nonce:            hexutil.Uint64(t.Nonce),
		TransactionIndex: &index,
		Value:            &value,
		V:                &V,
		R:                &R,
		S:                &S,
	}
	var to common.Address
	if len(t.To) > 0 {
		to = common.HexToAddress(t.To)
		result.To = &to
	}
	return result
}

func convertTransactionReceipt(receipt types.TransactionReceipt) *evmtypes.TransactionReceipt {
	result := &evmtypes.TransactionReceipt{
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
	result.Logs = convertLogs(receipt.Logs, nil)
	return result
}

func convertLogs(transactionLogs []types.TransactionLog, filterTopics [][]common.Hash) []*ethtypes.Log {
	ethLogs := make([]*ethtypes.Log, 0) // 这里为了和以太坊eth_getLogs接口兼容，所以给用make初始化ehtLogs,为空时返回[],而不是null
	for _, v := range transactionLogs {
		topics := make([]common.Hash, len(v.Topics))
		for i, t := range v.Topics {
			topics[i] = common.HexToHash(t.Topic)
		}
		// match topics
		if len(filterTopics) > 0 && !matchTopics(topics, filterTopics) {
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
	return ethLogs
}
