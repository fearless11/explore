#!/bin/bash
#######################################################
# date: 2019-05-15
# desc: 检测es中日志的状态
# (1). 5分钟内，状态码(500、502、504)出现超过20个告警
# (2). 5分钟内，慢响应(超过10s)出现超过20个告警
#######################################################

if ! pgrep -f '/usr/local/logmonitor/log_exporter' &> /dev/null ;then
    nohup  /usr/local/logmonitor/log_exporter >/usr/local/logmonitor/log 2>&1 &
fi
