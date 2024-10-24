# Storage CLI 工具文档

## 简介

Storage CLI 是一个用于文件和内容上传、下载、验证、批量操作、所有者管理和模板管理的命令行工具。它支持多种操作，包括文件上传、内容追加、零知识证明生成和验证等。

## 安装

请确保您已安装 Go 语言环境

```sh
go install github.com/wangdayong228/0g-storage-contracts/storage-cli
```

## 配置

在使用工具之前，请确保配置文件已正确设置。配置文件名称为`config.yaml`且应与工作目录同级。

编写`config.yaml`可参考`config.yaml.example`文件。

## 准备工作

运行命令前需要确保账户中有足够的 CFX 余额，否则无法正常运行。

该工具提供了初始化账户余额的命令，可以用来为指定账户充值。
```sh
cd ./scripts/init && go run .
```

## 使用方法

### 基本命令

- **上传文件**

  上传文件到存储节点：

  ```sh
  storage-cli upload file --file ./file_1.txt
  ```

- **上传内容**

  上传内容为 pointer数据类型 到存储节点：

  ```sh
  storage-cli upload content --name content1 --file ./go.sum --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
  ```

  上传内容为 text数据类型 到存储节点：

  ```sh
  storage-cli upload content --name content1 --content "hello" --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
  ```

- **追加内容**

  向指定内容追加内容：

  ```sh
  storage-cli append --name content1 --data " world" --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE
  ```

- **下载内容**

  下载指定内容：

  ```sh
  storage-cli download content --name content1
  ```

  下载内容并输出到终端：
  ```sh
  storage-cli download content --name content2 --console
  ```

- **所有者转移**

  转移内容的所有权：

  ```sh
  storage-cli owner transfer --name content1 --from 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE --to 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b
  ```
- **所有权查询**

  查询账户是否拥有内容所有权：

  ```sh
  storage-cli owner content --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE --name content1 
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
storage-cli batch upload -c 100000
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

2. **配置中的节点信息怎么填？**

   需要分别搭建 `storage-node`,`kv-node`,`blockchain` 节点服务。

3. **如何确保数据的安全性？**

   upload 命令可以使用加密方法上传文件，并在配置中设置合适的加密参数。
