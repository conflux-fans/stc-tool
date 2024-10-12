[toc]
# zk文档
https://gaudy-hub-f32.notion.site/ZK-a15c6dbe12dc4ae28e637a4bd46e4b0b?pvs=4

# 环境搭建

## 部署合约并设置到所有环境中

```sh
git clone https://github.com/wangdayong228/0g-storage-contracts.git
cd zerog-storage-contracts
git checkout add-template
make renew
```
renew 命令做了以下工作
### 1. 部署 zg 相关合约

部署 flow 和 ProMine 合约

### 2. 配置 zg-node

配置文件
```conf
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

<!-- 版本 commithash: 306c43c9dca6645da56c37f3337b08f39eb30cfa -->
**版本 1.0.0-testnet**
### 3. 配置 zg-kv

配置文件
```conf
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

<!-- 版本 commithash: bacf761d0f26af64b6375850ba2e9987ada93dc7 -->
**版本 1.0.0-testnet**

### 4. 配置本工具 storage-tool

```yaml

blockChain:
  url: "https://etest-rpc.nftrainbow.cn/JwtQFtZXar"
  # url: "http://127.0.0.1:8545"
  flowContract: "0xEA6718Cab1eA7aaa61D2a28f0297D6F2Ca194647"
  templateContract: "0x34Ab680c8De93aA0742EF5843520E86239B954EF"
storageNodes:
  - "http://127.0.0.1:11100"
kvNode: "http://127.0.0.1:6789"
privateKeys:
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23" #0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e24" #0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e25" #0xe61646FD48adF644404f373D984B14C877957F7c
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e26" #0xE7b3CafBf258804B867Df17e0AE5238030658a03
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e27" #0x8Faf8127849e4157dD086C923576a4029cA4E2B5
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e28" #0x0513B660EaBb10Ee88b8AC69188d3994f184a3D9
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e29" #0x60E54B5daD7331a85c3408A887588430B19b26D6
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e30" #0xB1b635163C5f58327b2FeD3a83131B6B209082C8
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e31" #0x581773C26661fA73f45516a72a138341F75a4cDD
  - "7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e32" #0xC933adff23Ce870B290C3D59b872855568eBE505
log : info # info,debug
```

# 工具命令

```sh
Storage cli for upload, append, verify, batchupload, owner manager and template manager

Usage:
  storage-cli [command]

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
  zk          generate zk proof and verify

Flags:
  -h, --help   help for storage-cli

Use "storage-cli [command] --help" for more information about a command.
```

## 命令示例
上传 by hash
```sh
storage-cli  upload file --file ~/tmp/random_files/file_1.txt
```

上传 by name
文件长度超过 1KB（config 中配置）, 以pointer方式上传
```sh
storage-cli  upload content --name content1 --file ./go.sum --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
```
指定文件，不超过 1kb， 以 text 方式上传
```sh
storage-cli  upload content --name content2 --file ～/tmp/random_files/file_1.txt  --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
```
直接指定内容，不超过 1kb，以 text 方式上传
```sh
storage-cli  upload content --name content3 --content hello --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE 
```

append
```sh
storage-cli  append --name content3 --data " world" --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE 
```

download
```sh
storage-cli  download content --name content3
```

batch 
```sh
storage-cli run . batch upload -c 100000 
```

owner transfer
```sh
storage-cli  owner transfer --name content2 --from 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE --to 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b 
```

owner check
```sh
storage-cli owner content --name content3 --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
```

zk proof
```sh
storage-cli zk proof -v '{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}' -t 20000101 -k verysecretkey123 -i uniqueiv12345678
```

zk verify
```sh
storage-cli zk verify --proof cbe38399d02b7bb83429c341ddfb726e0153f374953c41f147ae4279d11d879fce65525a19300bb1547781a634388e61601d957580a071b4baaeb2bfccab9c0bf7113f89ec0ab677f595b9667a3fa043641b46439bd396ca49bbeb484f8122147989370ae2f47d81df48c8731dc63a2f3f3c2afb53a4857f93f19270fc649407 --root 0xf0f5e6974aa06c1f6c2dae32b368a8608881d78c65ddc450d0fed9bb14a99431 --birth_threshold 20000101
```