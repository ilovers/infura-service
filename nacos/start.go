package nacos

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// Register start nacos client and register rpc service in nacos
func Register(urls string, namespace string, name string, externalAddr string) {
	ip, port, err := resolveIPAndPort(externalAddr)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to resolve %s error: %s", externalAddr, err.Error()))
		return
	}

	serverConfigs, err := getServerConfigs(urls)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to resolve nacos server url %s: %s", urls, err.Error()))
		return
	}
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig": constant.ClientConfig{
			TimeoutMs:           5000,
			ListenInterval:      10000,
			NotLoadCacheAtStart: true,
			NamespaceId:         namespace,
			LogDir:              "/dev/null",
			LogLevel:            "error",
		},
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create nacos client. error: %s", err.Error()))
		return
	}

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: name,
		Weight:      10,
		ClusterName: "DEFAULT",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"preserved.register.source": "GO",
			"app_registry_tag":          strconv.FormatInt(time.Now().Unix(), 10),
		},
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to register instance in nacos server. error: %s", err.Error()))
		return
	}
	log.Println("register application instance in nacos successfully")
}
