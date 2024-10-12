#!/bin/bash
set -x
source .profile

trap "trap - SIGTERM && kill 0" SIGINT SIGTERM EXIT
# ssh -vvv -qCN -p 9022 -L 14000:127.0.0.1:14000 wangdayong@$SHANGHAI_SERVER 安静模式
autossh -M 0 -vvv -CN -p 9022 -L 14000:127.0.0.1:14000 wangdayong@$SHANGHAI_SERVER &
autossh -M 0 -vvv -CN -p 9022 -L 15000:127.0.0.1:15000 wangdayong@$SHANGHAI_SERVER &
autossh -M 0 -vvv -CN -p 9022 -L 15001:127.0.0.1:15001 wangdayong@$SHANGHAI_SERVER &
autossh -M 0 -vvv -CN -p 9022 -L 15002:127.0.0.1:15002 wangdayong@$SHANGHAI_SERVER &
autossh -M 0 -vvv -CN -p 9022 -L 15003:127.0.0.1:15003 wangdayong@$SHANGHAI_SERVER &
autossh -M 0 -vvv -CN -p 9022 -L 16000:127.0.0.1:16000 wangdayong@$SHANGHAI_SERVER 

# 保持脚本运行
# while true; do
#     sleep 1
# done