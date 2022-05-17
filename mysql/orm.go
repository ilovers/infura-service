package mysql

import (
	"fmt"

	"github.com/okex/exchain/x/infura/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const maxSize = 10000

type Orm struct {
	db *gorm.DB
}

func NewOrm(url, user, pass, dbName string) (*Orm, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, url, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return &Orm{
		db: db,
	}, nil
}

func (orm *Orm) GetTransactionReceipt(txHash string) (receipts []types.TransactionReceipt, err error) {
	err = orm.db.Preload("Logs.Topics").Preload("Logs").Where("transaction_hash =?",
		txHash).Limit(1).Find(&receipts).Error // 这里使用Find而不是First的理由是：如果没有查询结果First会返回error
	return
}

func (orm *Orm) GetLogs(fromBlock, toBlock int64, addresses []string) (logs []types.TransactionLog, err error) {
	if len(addresses) == 0 {
		err = orm.db.Preload("Topics").Where("block_number >=? AND block_number<=?",
			fromBlock, toBlock).Limit(maxSize).Find(&logs).Error
	} else if len(addresses) == 1 {
		err = orm.db.Preload("Topics").Where("block_number >=? AND block_number<=? AND address=?",
			fromBlock, toBlock, addresses[0]).Limit(maxSize).Find(&logs).Error
	} else {
		err = orm.db.Preload("Topics").Where("block_number >=? AND block_number<=? AND address IN ?",
			fromBlock, toBlock, addresses).Limit(maxSize).Find(&logs).Error
	}
	return
}

func (orm *Orm) GetLogsByBlockHash(blockHash string, addresses []string) (logs []types.TransactionLog, err error) {
	if len(addresses) == 0 {
		err = orm.db.Preload("Topics").Where("block_hash=?",
			blockHash).Limit(maxSize).Find(&logs).Error
	} else if len(addresses) == 1 {
		err = orm.db.Preload("Topics").Where("block_hash=? AND address=?",
			blockHash, addresses[0]).Limit(maxSize).Find(&logs).Error
	} else {
		err = orm.db.Preload("Topics").Where("block_hash=? AND address IN ?",
			blockHash, addresses).Limit(maxSize).Find(&logs).Error
	}
	return
}

func (orm *Orm) GetBlockByNumber(blockNum int64) (block types.Block, err error) {
	err = orm.db.Preload("Transactions").Where("number=?", blockNum).First(&block).Error
	return
}

func (orm *Orm) GetBlockByHash(blockHash string) (block types.Block, err error) {
	err = orm.db.Preload("Transactions").Where("hash=?", blockHash).First(&block).Error
	return
}

func (orm *Orm) GetContractCode(address string) (code types.ContractCode, err error) {
	err = orm.db.Where("address=?", address).First(&code).Error
	return
}
