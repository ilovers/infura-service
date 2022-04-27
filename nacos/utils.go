package nacos

import (
	"net"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func getServerConfigs(urls string) ([]constant.ServerConfig, error) {
	// nolint
	var configs []constant.ServerConfig
	for _, url := range strings.Split(urls, ",") {
		addrs := strings.Split(url, ":")
		serverPort, err := strconv.Atoi(addrs[1])
		if err != nil {
			return nil, err
		}
		configs = append(configs, constant.ServerConfig{
			IpAddr: addrs[0],
			Port:   uint64(serverPort),
		})
	}
	return configs, nil
}

func resolveIPAndPort(addr string) (string, int, error) {
	lAddr := strings.Split(addr, ":")
	ip := lAddr[0]
	if ip == "127.0.0.1" {
		return getLocalIP(), 8080, nil
	}
	port, err := strconv.Atoi(lAddr[1])
	if err != nil {
		return "", 0, err
	}
	return ip, port, nil
}

// getLocalIP get local ip
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
