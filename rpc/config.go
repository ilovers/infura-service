package rpc

import "errors"

type Config struct {
	Address          string
	NacosUrl         string
	NacosNamespaceId string
	NacosServiceName string
	NacosServiceAddr string
	MysqlUrl         string
	MysqlUser        string
	MysqlPass        string
	RedisUrl         string
	RedisAuth        string
	RedisDB          int
}

func validateConfig(config *Config) error {
	if config.MysqlUser == "" || config.MysqlUser == "" {
		return errors.New("must set mysql url or user")
	}
	return nil
}
