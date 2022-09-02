package watcher

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// LocalAreaIPReg 内网IP正则
const LocalAreaIPReg = `(^127\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$)|(^10\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$)|(^172\.1[6-9]{1}[0-9]{0,1}\.[0-9]{1,3}\.[0-9]{1,3}$)|(^172\.2[0-9]{1}[0-9]{0,1}\.[0-9]{1,3}\.[0-9]{1,3}$)|(^172\.3[0-1]{1}[0-9]{0,1}\.[0-9]{1,3}\.[0-9]{1,3}$)|(^192\.168\.[0-9]{1,3}\.[0-9]{1,3}$)`

type nginxFields struct {
	Status            string `json:"status"`
	Timestamp         string `json:"timestamp"`
	RemoteAddr        string `json:"remote_addr"`
	HTTPXForwardedFor string `json:"http_x_forwarded_for"`
	HTTPHost          string `json:"http_host"`
	ServerPort        string `json:"server_port"`
	Scheme            string `json:"scheme"`
	RequestMethod     string `json:"request_method"`
	RequestURI        string `json:"request_uri"`
	UpstreamAddr      string `json:"upstream_addr"`
	BodyBytesSent     string `json:"body_bytes_sent"`
	BytesSent         string `json:"bytes_sent"`
	RequestTime       string `json:"request_time"`
	HTTPReferer       string `json:"http_referer"`
	HTTPUserAgent     string `json:"http_user_agent"`
	Host              string `json:"host"`
}

func DecodeLog(line string) {
	//Debug(line)
	var nf = nginxFields{}
	if err := json.Unmarshal([]byte(line), &nf); err != nil {
		Warn(fmt.Sprintf("解析json失败，当前处理行：%s 处理错误原因 %s", line, err.Error()))
		return
	}
	status, _ := strconv.Atoi(nf.Status)
	hostName := nf.HTTPHost
	hostCache := cache.getHostCache(hostName)
	if status < 400 {
		hostCache.cleanHostCError()
		return
	}

	if status >= 400 {
		hostCache.inCreCErrorCount()
		hostCache.inCreTErrorCount()

		re := regexp.MustCompile(LocalAreaIPReg)
		match := re.MatchString(nf.RemoteAddr)
		if !match {
			hostCache.inCreErrorCounter(IpErrType, nf.RemoteAddr)
			hostCache.inCreErrorCounter(UrlErrType, nf.RequestURI)
		}
		cCount := hostCache.getCErrorCounter()
		pCount := hostCache.getPErrorCounter()
		Debug("当前cCount数为: " + strconv.Itoa(cCount) + " pCount: " + strconv.Itoa(pCount))
		if cCount > app.reportRate.continuousErrorCount || pCount > app.reportRate.minutePeriodErrorCount {
			topIpMap := hostCache.getTopTagError(IpErrType, 1)
			var ip, url string
			var ipCount, urlCount int
			for key, value := range topIpMap {
				ip = key
				ipCount = value
			}
			topUrl := hostCache.getTopTagError(UrlErrType, 1)
			for key, value := range topUrl {
				url = key
				urlCount = value
			}
			alarm(hostName, cCount, pCount, ip, ipCount, url, urlCount)
		}
	}
}
