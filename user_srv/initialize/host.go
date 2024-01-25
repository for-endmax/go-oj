package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"user_srv/global"
)

func GetHost() {
	host := GetCheckHost()
	if host != "" {
		global.ServerConfig.CheckHost = host
		zap.S().Infof("更新，本地host为:%s", host)
	}
}

func GetCheckHost() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取本机IP地址时发生错误:", err)
		panic(err)
	}

	for _, addr := range addrs {
		// 检查IP地址是否是IPv4或IPv6，并排除一些特殊情况
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				return ipnet.IP.String()
			} else if ipnet.IP.To16() != nil {
				fmt.Println("IPv6地址:", ipnet.IP.String())
			}
		}
	}
	return ""
}
