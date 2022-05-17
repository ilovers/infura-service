package cmd

import (
	"log"

	"github.com/okex/infura-service/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagAddress             = "address"
	flagNacosUrl            = "nacos_url"
	flagNacosNamespaceID    = "nacos_namespace_id"
	flagNacosServiceName    = "nacos-service-name"
	flagNacosServiceAddress = "nacos_service_address"
	flagMysqlUrl            = "mysql-url"
	flagMysqlUser           = "mysql-user"
	flagMysqlPass           = "mysql-pass"
	flagMysqlDB             = "mysql-db"
	flagRedisUrl            = "redis-url"
	flagRedisAuth           = "redis-auth"
	flagRedisDB             = "redis-db"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start infura service",
		Run: func(cmd *cobra.Command, args []string) {
			starService()
		},
	}
	bindStartFlags(cmd)
	return cmd
}

func bindStartFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagAddress, ":8080", "Listen address")
	cmd.Flags().String(flagNacosUrl, "", "Nacos server urls for discovery of rpc service")
	cmd.Flags().String(flagNacosNamespaceID, "", "Nacos namespace id for discovery of rpc service")
	cmd.Flags().String(flagNacosServiceName, "", "Rpc service name in nacos")
	cmd.Flags().String(flagNacosServiceAddress, "127.0.0.1:8080", "Rpc service address register to nacos")
	cmd.Flags().String(flagMysqlUrl, "127.0.0.1:3306", "Mysql url(host:port) of rpc service")
	cmd.Flags().String(flagMysqlUser, "root", "Mysql user of rpc service")
	cmd.Flags().String(flagMysqlPass, "root", "Mysql password of rpc service")
	cmd.Flags().String(flagMysqlDB, "infura", "Mysql db name of rpc service")
	cmd.Flags().String(flagRedisUrl, "127.0.0.1:6379", "Redis url(host:port) of infura rpc service")
	cmd.Flags().String(flagRedisAuth, "", "Redis auth of rpc service")
	cmd.Flags().Int(flagRedisDB, 0, "Redis db of rpc service")
	viper.BindPFlags(cmd.Flags())
}

func starService() {
	config := initConfig()
	service, err := rpc.New(config)
	if err != nil {
		log.Fatal(err)
	}
	service.Start()
}

func initConfig() *rpc.Config {
	return &rpc.Config{
		Address:          viper.GetString(flagAddress),
		NacosUrl:         viper.GetString(flagNacosUrl),
		NacosNamespaceId: viper.GetString(flagNacosNamespaceID),
		NacosServiceName: viper.GetString(flagNacosServiceName),
		NacosServiceAddr: viper.GetString(flagNacosServiceAddress),
		MysqlUrl:         viper.GetString(flagMysqlUrl),
		MysqlUser:        viper.GetString(flagMysqlUser),
		MysqlPass:        viper.GetString(flagMysqlPass),
		MysqlDB:          viper.GetString(flagMysqlDB),
		RedisUrl:         viper.GetString(flagRedisUrl),
		RedisAuth:        viper.GetString(flagRedisAuth),
		RedisDB:          viper.GetInt(flagRedisDB),
	}
}
