package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"record_srv/global"
	"strconv"
)

// GetFreePort 生成可用端口号
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// PortCheck 检查端口是否可用
func PortCheck() bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(global.ServerConfig.Port)))

	if err != nil {
		return false
	}
	defer l.Close()
	return true
}

// GetPort 获取端口
func GetPort() {
	//如果端口可用，则不操作
	if PortCheck() {
		return
	}
	//如果端口不可用，则随机获取端口
	zap.S().Infof("默认端口 %d 不可用,正在获取空闲端口", global.ServerConfig.Port)
	port, err := GetFreePort()
	if err != nil {
		panic("获取空闲端口失败")
	}
	global.ServerConfig.Port = port
	zap.S().Info("获取空闲端口成功")
}
