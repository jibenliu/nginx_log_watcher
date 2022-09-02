package watcher

import (
	"fmt"
	"net"
	"strings"
)

func InitServer() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9999,
	})
	if err != nil {
		Error("listen failed, err:" + err.Error())
	}
	defer listen.Close()
	for {
		var data [1024]byte
		n, _, err := listen.ReadFromUDP(data[:]) // 接收数据
		if err != nil {
			fmt.Println("read udp failed, err:", err)
			continue
		}
		msg := string(data[:n])
		Debug("入参信息为：" + msg)
		logMsg := data[strings.Index(msg, "nginx_access_log: ")+18 : n]
		DecodeLog(string(logMsg))
	}
}
