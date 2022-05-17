package rpc

import (
	"log"

	"github.com/okex/infura-service/redis"

	"github.com/okex/infura-service/mysql"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/okex/infura-service/rpc/namespaces/eth"
)

const (
	ethNamespace = "eth"
	apiVersion   = "1.0"
)

// getAPIs returns the list of all APIs from the Ethereum namespaces
func getAPIs(config *Config) []rpc.API {
	orm, err := mysql.NewOrm(config.MysqlUrl, config.MysqlUser, config.MysqlPass, config.MysqlDB)
	if err != nil {
		log.Fatal(err)
	}
	redisCli := redis.NewClient(config.RedisUrl, config.RedisAuth, config.RedisDB)
	ethAPI, err := eth.NewAPI(orm, redisCli)
	if err != nil {
		log.Fatal(err)
	}
	apis := []rpc.API{
		{
			Namespace: ethNamespace,
			Version:   apiVersion,
			Service:   ethAPI,
			Public:    true,
		},
	}
	return apis
}
