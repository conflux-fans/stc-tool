# Storage CLI 工具文档

## 简介

Storage CLI 是一个用于文件和内容上传、下载、验证、批量操作、所有者管理和模板管理的命令行工具。它支持多种操作，包括文件上传、内容追加、零知识证明生成和验证等。

## 安装

请确保您已安装 Go 语言环境

```sh
go install github.com/wangdayong228/0g-storage-contracts/storage-cli
```

## 配置

在使用工具之前，请确保配置文件已正确设置。以下是一个示例配置：

```yaml

blockChain:
  url: http://127.0.0.1:14000
  flowContract: 0x33f2CFc729Bd870fA54b5032660e06B4CF2a7F94 
  templateContract: 0x34Ab680c8De93aA0742EF5843520E86239B954EF
  pmContract: 0x588D57Fb016CEE89513B9B7ee78AeB8b56BAc85D
storageNodes:
  - http://127.0.0.1:15000
  - http://127.0.0.1:15001
  - http://127.0.0.1:15002
  - http://127.0.0.1:15003
kvNode: http://127.0.0.1:16000
kvStreamId: 000000000000000000000000000000000000000000000000000000000000f009
zkNode: http://127.0.0.1:3030
privateKeys:
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e23
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e24
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e25
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e26
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e27
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e28
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e29
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e30
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e31
  - 7c5da44cf462b81e0b61a582f8c9b23ca78fc23e7104138f4e4329a9b2076e32
log : info # info,debug
extendData:
  textMaxSize: 1024 # 大于该长度的数据将以 pointer 扩展数据类型存储
```

## 使用方法

### 基本命令

- **上传文件**

  上传文件到存储节点：

  ```sh
  storage-cli upload file --file ~/tmp/random_files/file_1.txt
  ```

- **上传内容**

  上传内容到存储节点：

  ```sh
  storage-cli upload content --name content1 --file ./go.sum --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
  ```

- **追加内容**

  向指定内容追加内容：

  ```sh
  storage-cli append --name content3 --data " world" --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
  ```

- **下载内容**

  下载指定内容：

  ```sh
  storage-cli download content --name content3
  ```

- **所有者转移**

  转移内容的所有权：

  ```sh
  storage-cli owner transfer --name content2 --from 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE --to 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b
  ```
- **所有权查询**

  查询账户是否拥有内容所有权：

  ```sh
  storage-cli owner content --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE --name content2 
  ```

### 零知识证明

- **生成证明**

  生成零知识证明：

  ```sh
  storage-cli zk proof -v '{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}' -t 20000101 -k verysecretkey123 -i uniqueiv12345678
  ```

- **验证证明**

  验证零知识证明：

  ```sh
  storage-cli zk verify --proof <proof> --root <root_hash> --birth_threshold 20000101
  ```

## 高级功能

### 批量上传

批量上传 N 条随机数据：

```sh
storage-cli run . batch upload -c 100000
```

### 模板管理

- **创建模板**

  创建新的存证模板：

  ```sh
  storage-cli template create --name test4 --keys name,age
  ```

- **下载模板**

  下载存证模板：

  ```sh
  storage-cli template download --name test4
  ```

## 常见问题

1. **如何配置多个存储节点？**

   在配置文件中，`storageNodes` 字段可以接受一个节点列表。

2. **如何处理上传失败的情况？**

   请检查网络连接和配置文件中的节点地址是否正确。

3. **如何确保数据的安全性？**

   使用加密方法上传文件，并在配置中设置合适的加密参数。

## 版本信息

当前版本：1.0.0-testnet

## 贡献

欢迎提交问题和贡献代码！请访问我们的 [GitHub 仓库](https://github.com/wangdayong228/0g-storage-contracts) 了解更多信息。