package mysql

import (
	"fmt"

	"gorm.io/gorm/logger"

	"github.com/okex/exchain/x/infura/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const maxSize = 10000

type Orm struct {
	db *gorm.DB
}

func NewOrm(url, user, pass string) (*Orm, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/infura?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, url)

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
		txHash).Limit(1).Find(&receipts).Error
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
