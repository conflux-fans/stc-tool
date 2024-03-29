# 环境搭建
## 1. 部署 zg 相关合约

``` 
git clone https://github.com/wangdayong228/0g-storage-contracts.git
cd zerog-storage-contracts
git checkout 12cecf88752e706413c1da33d476123e244402b0
npx hardhat run --network ecfx_test ./scripts/deploy.ts 
```
将部署 flow 和 ProMine 合约

版本 commithash: 12cecf88752e706413c1da33d476123e244402b0

## 2. zg-rust

配置文件
```
log_config_file = "log_config"
network_libp2p_port = 11000
network_discovery_port = 11000
rpc_listen_address = "127.0.0.1:11100"
log_contract_address = "0xF5587366B9bDA86854f765A4B6184bDd5dBdFa8E"
mine_contract_address = "0x4AF117e7B516969488EDee50a8f7afFb48C62bCb"
blockchain_rpc_endpoint = "https://evmtestnet.confluxrpc.com/6XWH2kDUX4wcKVN1VThMpjhhwerkTMZR8GYjk3S8Ti6GhM8qw7TJXDuT4sJWsM8MNmz2oxLsWAbjDUELaeAG4QA9Y"
network_libp2p_nodes = []
log_sync_start_block_number = 164900000
```
- log_contract_address 配置为步骤 1 部署的 flow 合约
- mine_contract_address 配置为不走 1 部署的 ProMine 合约
- log_sync_start_block_number 配置为不大于 flow 合约部署区块高度即可

版本 commithash: 306c43c9dca6645da56c37f3337b08f39eb30cfa
## 3. zg-kv

配置文件
```
stream_ids = ["000000000000000000000000000000000000000000000000000000000000f2bd", "000000000000000000000000000000000000000000000000000000000000f009", "0000000000000000000000000000000000000000000000000000000000016879", "0000000000000000000000000000000000000000000000000000000000002e3d"]


db_dir = "db"
kv_db_dir = "kv.DB"

rpc_enabled = true
rpc_listen_address = "127.0.0.1:6789"
zgs_node_urls = "http://127.0.0.1:11100"

log_config_file = "log_config"

blockchain_rpc_endpoint = "https://evmtestnet.confluxrpc.com/6XWH2kDUX4wcKVN1VThMpjhhwerkTMZR8GYjk3S8Ti6GhM8qw7TJXDuT4sJWsM8MNmz2oxLsWAbjDUELaeAG4QA9Y"
log_contract_address = "0xF5587366B9bDA86854f765A4B6184bDd5dBdFa8E"
log_sync_start_block_number = 164900000

```
- log_contract_address 配置为步骤 1 部署的 flow 合约
- log_sync_start_block_number 配置为不大于 flow 合约部署区块高度即可
- stream_ids 每个 stream 可以看成是一个数据库，写文件只能使用配置好的 stream，否则不生效（也不报错）

版本 commithash: bacf761d0f26af64b6375850ba2e9987ada93dc7

# 工具命令

```
Zerog storage tool for upload, append, verify, batchupload, owner manager and template manager

Usage:
  zerog-storage-tool [command]

Available Commands:
  append      Append content to specified file
  batch       Batch operations
  download    Download file or content
  file        File operations
  help        Help about any command
  owner       Owner operations
  template    Template opertaions
  upload      Upload file or content
  verify      Verify file

Flags:
  -h, --help   help for zerog-storage-tool

Use "zerog-storage-tool [command] --help" for more information about a command.
```