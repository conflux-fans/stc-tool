# 最终交付测试环境
```sh
cli 节点： ssh ubuntu@16.163.233.63
zk 节点： ssh ubuntu@16.163.233.63 -p 2222
 
url:
blockchain: http://manager:8545/
kv: http://manager:6789/
storage: http://node0:5678/,http://node1:5678/,http://node2:5678/,http://node3:5678/
 
具体的 ip 我写到开发环境的 host 文件里了, 包括 zk 节点和 cli 节点
```

# 上海服务器测试环境
blockchain eth rpc 在 14000 端口, storage 在 15000~15003 端口（分片模式，必须全部指定），kv 在 16000 端口
- Mine 0x615d7c6A00335688F35D40Bb27735B7B2DBfDa97
- Flow 0x33f2CFc729Bd870fA54b5032660e06B4CF2a7F94

## `0g-storage-client` 上传命令
```sh
./0g-storage-client upload \
--contract 0x33f2CFc729Bd870fA54b5032660e06B4CF2a7F94 \
--node http://127.0.0.1:15000,http://127.0.0.1:15001,http://127.0.0.1:15002,http://127.0.0.1:15003 \
--url http://127.0.0.1:14000/ \
--key 9a6d3ba2b0c7514b16a006ee605055d71b9edfad183aeb2d9790e9d4ccced471 \
--file bulk
```