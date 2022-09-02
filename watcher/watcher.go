package watcher

import (
	"strconv"
	"time"
)

func alarm(hostName string, cErrorCount, pErrorCount int, topIp string, topIpCnt int, topUrl string, topUrlCnt int) {
	send := &sendStruct{}
	send.name = "● 告警节点: " + hostName
	send.time = time.Now().Format("2006-01-02 15:04:05 +08:00")
	send.level = " 问题等级: 警告"
	send.detail = " 问题详情:  \n" +
		"        nginx错误请求达到阈值: \n" +
		"        连续请求错误数为：" + strconv.Itoa(cErrorCount) + " \n" +
		"        每分钟内请求错误数达到：" + strconv.Itoa(pErrorCount) + "\n" +
		"        最高请求IP为：" + topIp + "，数量为：" + strconv.Itoa(topIpCnt) + "\n" +
		"        最高请求地址为：" + topUrl + "，最高请求数量为：" + strconv.Itoa(topUrlCnt)
	Warn(send.detail)
	sendAlarm(send)
}
