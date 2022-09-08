# nginx 日志监控处理转发系统

## nginx 可以按照如下配置

```shell
http{
  ...
  log_format json escape=json '{"timestamp":"$msec",'
                    '"remote_addr":"$remote_addr",'
					'"http_x_forwarded_for":"$http_x_forwarded_for",'
					'"http_host":"$server_name",'
					'"server_port":"$server_port",'
					'"scheme":"$scheme",'
                    '"request_method":"$request_method",'
                    '"request_uri":"$request_uri",'
                    '"status":"$status",'
                    '"upstream_addr":"$upstream_addr",'
                    '"body_bytes_sent":"$body_bytes_sent",'
                    '"bytes_sent":"$bytes_sent",'
					'"request_time":"$request_time",'
#					'"request_body":"$request_body",' #UDP协议有长度限制会被截包
					'"http_referer":"$http_referer",'
					'"http_user_agent":"$http_user_agent",'
                    '"host":"$host"'
                    '}';
  ...
  
                    
}

server{
  ...
  access_log syslog:server=127.0.0.1:9999,facility=local7,tag=nginx_access_log,severity=info json;
  ...
}
```

## 重启nginx

## 运行该程序，即可将nginx告警消息推送到 telegram 上面去

# 该程序放弃了 `go tail`类库来监控nginx日志，在nginx日志过大的场景下会消耗巨量的内存