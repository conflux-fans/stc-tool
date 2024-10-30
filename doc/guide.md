# Storage CLI 用户手册

# 准备工作

storage-cli 是一个基于区块链技术的综合数据治理平台命令行工具，提供文件管理、数据访问控制和隐私保护等功能。作为去中心化系统，它采用区块链私钥进行身份验证和权限管理，而非传统的用户名密码机制。

使用 storage-cli 的第一步是配置区块链私钥。系统通过私钥进行身份验证、交易签名和访问控制。

### 私钥文件位置

私钥文件 `private_keys` 默认存放位置：
- Linux/macOS: `$HOME/.storage-cli/private_keys`
- Windows: `%USERPROFILE%\.storage-cli\private_keys`

### 私钥文件格式

私钥文件采用纯文本格式，每行包含一个私钥：

```text
0x123456789abcdef...
0x987654321fedcba...
```

注意事项：
- 每个私钥必须是有效的十六进制字符串
- 私钥长度必须是 64 字节（128 个十六进制字符）
- 私钥可以包含或省略 "0x" 前缀
- 文件权限建议设置为 600（仅所有者可读写）

### 私钥安全建议

1. **安全存储**
   - 使用安全的方式生成私钥
   - 保管好私钥的备份
   - 不要在多个设备间共享私钥文件

2. **权限控制**
   ```bash
   # Linux/macOS 设置私钥文件权限
   chmod 600 ~/.storage-cli/private_keys
   ```

3. **定期更新**
   - 建议定期更换私钥
   - 更换私钥后及时更新相关的数据访问权限


# 一、文件加密与解密

storage-cli 工具提供文件加密和解密功能，使用 AES 加密算法对文件内容进行保护。加密后的文件只能通过正确的密码和加密算法进行解密，可有效防止文件内容泄露。


## 1.1 文件加密

文件加密功能使用 AES-256-CBC 加密算法将普通文件转换为加密文件。加密后的文件将以二进制形式存储，无法直接查看或编辑其内容。

### 命令语法

```shell
storage-cli file encrypt --source <SOURCE_FILE> --output <OUTPUT_DIR> --cipher <CIPHER_METHOD> --password <PASSWORD>
```

### 选项说明

| 选项                       | 必填 | 说明                          |
| -------------------------- | ---- | ----------------------------- |
| `--source <SOURCE_FILE>`   | 是   | 指定需要加密的源文件路径      |
| `--output <OUTPUT_DIR>`    | 是   | 指定加密后文件的输出目录      |
| `--cipher <CIPHER_METHOD>` | 是   | 指定加密算法，当前支持：`aes` |
| `--password <PASSWORD>`    | 是   | 加密密码，至少16位字符        |

### 示例

**示例 1.1**：将 `example.txt` 文件通过 AES 加密后输出到 `./encrypted` 目录：

```shell
storage-cli file encrypt \
   --source example.txt \
   --output ./encrypted \
   --cipher AES_CBC \
   --password mypassword123456
```

### 注意事项
- 加密会产生额外的元数据，因此加密后的文件大小会略大于源文件
- 加密密码无法找回，密码丢失将导致永久无法解密文件内容
- 对于重要文件建议先进行备份，再执行加密操作
- 加密过程会占用系统资源，大文件加密耗时较长

## 1.2 文件解密

文件解密功能用于将加密文件还原为原始内容。解密时需要提供加密时使用的密码和加密算法，两者必须与加密时完全一致。解密成功后可以获得与源文件相同的文件内容

### 命令语法

```shell
storage-cli file decrypt --source <ENCRYPTED_FILE> --output <OUTPUT_DIR> --cipher <CIPHER_METHOD> --password <PASSWORD>
```

### 选项说明

| 选项                        | 必填 | 说明                                                    |
| --------------------------- | ---- | ------------------------------------------------------- |
| `--source <ENCRYPTED_FILE>` | 是   | 指定需要解密的加密文件路径                              |
| `--output <OUTPUT_DIR>`     | 是   | 指定解密后文件的输出目录                                |
| `--cipher <CIPHER_METHOD>`  | 是   | 指定解密算法，需与加密时使用的方法一致，当前支持：`aes` |
| `--password <PASSWORD>`     | 是   | 解密密码，需与加密时使用的密码一致                      |

### 示例

**示例 1.2**：将 `encrypted/example.txt.encrypt` 文件通过 AES 解密后输出到 `./decrypted` 目录：

```shell
storage-cli file decrypt \
   --source encrypted/example.txt.encrypt \
   --output ./decrypted \
   --cipher AES_CBC \
   --password mypassword123456
```

### 注意事项
- 输入错误的密码将导致解密失败，但不会损坏加密文件
- 如果输出目录下存在同名文件，解密操作会覆盖已有文件，建议使用空目录作为输出目录
- 解密过程中请勿中断操作，以免生成不完整的文件
- 解密大文件时可能需要较长时间，请耐心等待

# 二、文件数据管理

storage-cli 工具支持将文件上传到去中心化存储系统、从系统中下载文件，以及验证已上传文件的完整性。上传时可以选择对文件进行加密保护，下载时支持通过文件共享代码或根哈希值获取文件。

## 2.1 文件上传

文件上传功能支持将本地文件上传到去中心化存储系统。上传过程中可以选择对文件进行 AES 加密，加密后的文件内容在存储和传输过程中都将保持加密状态。上传成功后会生成文件共享代码和根哈希值，可用于后续下载文件。

### 命令语法

```shell
storage-cli upload file --file <FILE_NAME> [--cipher <CIPHER_METHOD> --password <PASSWORD>]
```

### 选项说明

| 选项                       | 必填 | 说明                               |
| -------------------------- | ---- | ---------------------------------- |
| `--file <FILE_NAME>`       | 是   | 需要上传的文件名称或路径。支持相对路径和绝对路径 |
| `--cipher <CIPHER_METHOD>` | 否   | 加密算法。当需要加密上传时使用，当前支持：`aes` |
| `--password <PASSWORD>`    | 否   | 加密密码。使用加密上传时必须提供，建议使用强密码 |

### 注意事项

- 上传大文件时会自动进行分片处理
- 网络状况会影响上传速度
- 上传成功后请妥善保存文件共享代码和根哈希值
- 使用加密上传时，密码丢失将导致无法解密下载的文件

### 示例

**示例 2.1**：将 `example.txt` 文件上传到去中心化存储系统：

```shell
storage-cli upload file --file example.txt
```

**示例 2.2**：将 `example.txt` 文件通过 AES 加密后上传到去中心化存储系统：

```shell
storage-cli upload file \
   --file example.txt \
   --cipher AES_CBC \
   --password mypassword123456
```

## 2.2 文件下载

文件下载功能用于从去中心化存储系统获取已上传的文件。下载时可以使用文件的共享代码或根哈希值，两者指定其一即可。如果下载的是加密文件，需要提供正确的解密密码才能获得原始文件内容。

### 命令语法

```shell
storage-cli download file (--code <FILE_SHARE_CODE> | --root <ROOT_HASH>)
```

### 选项说明

| 选项                       | 必填               | 说明                         |
| -------------------------- | ------------------ | ---------------------------- |
| `--code <FILE_SHARE_CODE>` | 与 `--root` 二选一 | 文件共享代码。上传后生成的便于分享的短代码 |
| `--root <ROOT_HASH>`       | 与 `--code` 二选一 | 文件根哈希值。上传后生成的唯一标识符 |

### 注意事项

- 下载前确认本地磁盘空间充足
- 大文件下载过程中请保持网络连接稳定
- 如果下载加密文件，请准备好解密所需的密码
- 下载文件会自动进行完整性校验

### 示例

**示例 2.3**：使用文件共享代码下载文件：

```shell
storage-cli download file --code abc123def456
```

**示例 2.4**：使用根哈希值下载文件：

```shell
storage-cli download file --root 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4
```

## 2.3 文件验证

文件验证功能用于检查本地文件是否与去中心化存储系统中的文件一致。验证过程会计算文件哈希值并进行比对。对于加密文件，需要提供加密时使用的算法和密码才能进行验证。

### 命令语法

```shell
storage-cli verify --file <FILE_PATH> [--cipher <CIPHER_METHOD> --password <PASSWORD>]
```

### 选项说明

| 选项                       | 必填 | 说明                                           |
| -------------------------- | ---- | ---------------------------------------------- |
| `--file <FILE_PATH>`       | 是   | 需要验证的文件路径。支持相对路径和绝对路径 |
| `--cipher <CIPHER_METHOD>` | 否   | 文件加密算法。验证加密文件时需要提供，当前支持：`aes` |
| `--password <PASSWORD>`    | 否   | 加密密码。验证加密文件时必须提供原始密码 |

### 注意事项

- 验证过程会读取整个文件内容
- 大文件验证可能需要较长时间
- 加密文件验证需要使用原始密码

### 示例

**示例 2.5**：验证非加密文件：

```shell
storage-cli verify --file example.txt
```

**示例 2.6**：验证加密文件：

```shell
storage-cli verify \
   --file example.txt \
   --cipher AES_CBC \
   --password mypassword123456
```

# 三、富标记数据管理

富标记数据管理功能支持将内容作为带有所有者标记的数据上传到去中心化存储系统。与普通文件上传不同，上传的内容会绑定所有者信息，并支持后续追加内容。

## 3.1 数据上传

数据上传功能支持两种上传方式：直接提供字符串内容，或指定包含内容的文件路径。每个上传的内容都需要指定所有者账户地址和内容名称作为标识。

### 命令语法

```shell
storage-cli upload content --account <ACCOUNT_ADDRESS> --name <CONTENT_NAME> (--content <CONTENT> | --file <FILE_PATH>)
```

### 选项说明

| 选项                          | 必填                | 说明                                 |
| ----------------------------- | ------------------- | ------------------------------------ |
| `--content <CONTENT>`         | 与 `file` 二选一    | 需要上传的内容，直接以字符串形式提供 |
| `--file <FILE_PATH>`          | 与 `content` 二选一 | 需要上传的内容所在文件的路径         |
| `--account <ACCOUNT_ADDRESS>` | 是                  | 指定内容所有者的账户地址             |
| `--name <CONTENT_NAME>`       | 是                  | 指定内容的名称，用于后续访问         |

### 示例

**示例 3.1**：将文本内容 "Hello, World!" 上传到去中心化存储系统：

```shell
storage-cli upload content \
   --content "Hello, World!" \
   --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
   --name "Greeting"
```

**示例 3.2**：将文件 `content.txt` 的内容上传到去中心化存储系统：

```shell
storage-cli upload content \
   --file content.txt \
   --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
   --name "FileContent"
```

## 3.2 数据下载

数据下载功能用于获取已上传的内容。系统支持将内容直接输出到控制台查看，或获取内容的元数据信息。下载时只需提供内容的名称即可。

### 命令语法

```shell
storage-cli download content [flags]
```

### 选项说明

| 选项                  | 必填 | 说明                                          |
| --------------------- | ---- | --------------------------------------------- |
| `--console`           | 否   | 将内容直接输出到控制台显示                    |
| `--metadata`          | 否   | 输出内容的元数据信息                          |
| `-n, --name <string>` | 是   | 要下载的内容名称，必须是上传时指定的内容名称  |

### 示例

**示例 3.3**：下载名称为 "ExampleData" 的内容并输出到控制台：

```shell
storage-cli download content --name "ExampleData" --console
```

**示例 3.4**：下载名称为 "ExampleData" 的内容并输出元数据：

```shell
storage-cli download content --name "ExampleData" --metadata
```

## 3.3 数据追加

数据追加功能允许为已存在的内容添加新的数据。追加操作仅支持小于 1KB 的直接存储内容。追加时可以选择直接提供追加内容，或指定包含追加内容的文件路径。

### 命令语法

```shell
storage-cli append --name <CONTENT_NAME> (--data <APPEND_DATA> | --file <FILE_PATH>) [--account <ACCOUNT_NAME>]
```

### 选项说明

| 选项                       | 必填             | 说明                               |
| -------------------------- | ---------------- | ---------------------------------- |
| `--name <CONTENT_NAME>`    | 是               | 要追加内容的名称                   |
| `--data <APPEND_DATA>`     | 与 `file` 二选一 | 要追加的内容，直接以字符串形式提供 |
| `--file <FILE_PATH>`       | 与 `data` 二选一 | 要追加的内容所在文件的路径         |
| `--account <ACCOUNT_NAME>` | 否               | 指定用于追加内容的账户名称         |

### 示例

**示例 3.5**：将文本内容 "Hello again!" 追加到名称为 Greeting 的内容后：

```shell
storage-cli append \
   --name "Greeting" \
   --data "Hello again!"
```

**示例 3.6**：将文件 `additional_content.txt` 中的内容追加到名称为 FileContent 的内容后：

```shell
storage-cli append \
   --name "FileContent" \
   --file additional_content.txt
```

# 四、数据所有权管理

数据所有权管理功能提供对[富标记数据](#三-富标记数据管理)的所有权进行管理，包括查询、转移和历史记录查看。这些操作仅适用于通过富标记数据管理功能上传的内容，普通文件上传不支持所有权管理。

## 4.1 所有权查询

所有权查询功能用于验证特定账户是否拥有某个内容的所有权。查询时需要提供账户地址和内容名称，系统将返回该账户是否为内容的当前所有者。

### 命令语法

```shell
storage-cli owner content --account <ACCOUNT_ADDRESS> --name <CONTENT_NAME>
```

### 选项说明

| 选项                          | 必填 | 说明               |
| ----------------------------- | ---- | ------------------ |
| `--account <ACCOUNT_ADDRESS>` | 是   | 需要检查的账户地址 |
| `--name <CONTENT_NAME>`       | 是   | 需要检查的内容名称 |

### 示例

**示例 4.1**：检查账户 `0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE` 是否拥有名称为 Greeting 的内容：

```shell
storage-cli owner content \
  --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
  --name "Greeting"
```

## 4.2 所有权转移

所有权转移功能允许内容的当前所有者将内容的所有权转移给其他账户。转移后，新所有者将获得对内容的完整控制权，包括内容的更新、追加和后续转移权限。

### 命令语法

```shell
storage-cli owner transfer --from <CURRENT_OWNER> --to <TARGET_OWNER> --name <CONTENT_NAME>
```

### 选项说明

| 选项                     | 必填 | 说明                     |
| ------------------------ | ---- | ------------------------ |
| `--from <CURRENT_OWNER>` | 是   | 当前内容所有者的账户地址 |
| `--to <TARGET_OWNER>`    | 是   | 目标所有者的账户地址     |
| `--name <CONTENT_NAME>`  | 是   | 要转移所有权的内容名称   |

### 示例

**示例 4.2**：将名称为 Greeting 的内容的所有权从一个账户转移到另一个账户：

```shell
storage-cli owner transfer \
  --from 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
  --to 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b \
  --name "Greeting"
```

## 4.3 所有权转移历史查询

所有权转移历史查询功能用于追踪内容所有权的完整变更记录。系统记录了每次所有权转移的时间、转出方和接收方信息，方便审计和追溯。

### 命令语法

```shell
storage-cli owner history --name <CONTENT_NAME>
```

### 选项说明

| 选项                    | 必填 | 说明                         |
| ----------------------- | ---- | ---------------------------- |
| `--name <CONTENT_NAME>` | 是   | 需要查询所有权历史的内容名称 |

### 示例

**示例 4.3**：查询名称为 "Greeting" 的内容的所有权转移历史：

```shell
storage-cli owner history --name "Greeting"
```

# 五、可验证证书与零知识证明

可验证证书（Verifiable Credentials，VC）与零知识证明（Zero-Knowledge Proof，ZKP）模块提供了一套完整的隐私保护数据验证解决方案。该模块支持用户提交符合系统定义数据模型规范的可验证证书，并利用零知识证明技术实现对证书数据的**可信查询**，在保护数据隐私的同时确保数据真实性。

通过系统提供的二次开发接口（详见[附录 A](附录a-可信证书技术原理与开发手册)），用户可以使用 circom 脚本语言自定义查询逻辑和验证规则。circom **支持通用计算**，能够实现复杂的数据处理和验证要求，系统将根据定制化的查询内容生成相应的零知识证明。

本章以 `check_age.circom` 验证电路为例，演示系统的完整交互流程。使用前请确保已按附录 A 的步骤完成验证电路的部署。

## 5.1 VC加密上传

VC 加密上传功能用于安全地将可验证证书数据存储到去中心化系统。系统采用对称加密算法保护数据，并自动生成用于零知识证明的必要参数。上传成功后，系统将在本地生成包含以下关键信息的配置文件：

- `key`: 数据加密密钥
- `iv`: 加密初始化向量
- `submission_tx_hash`: 区块链上的提交交易哈希
- `vc_data_root`: 加密数据的 Merkle 树根哈希

### 命令语法

```shell
storage-cli zk upload --vc <VC_JSON_STRING>
```

### 选项说明

| 选项               | 必填 | 说明                               |
| ------------------ | ---- | ---------------------------------- |
| `--vc <VC_STRING>` | 是   | 验证凭证字符串，需以 JSON 格式提供 |

### 示例

**示例 5.1**：上传VC数据到去中心化存储系统：

```shell
storage-cli zk upload \
   --vc '{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}'
```

### 输出示例

上传后，输出文件将包含以下信息：

- `key`: verysecretkey123
- `iv`: initialvector123
- `submission_tx_hash`: "0x276b14f314e7d3583c6718c75f8fc2e1d89b0f13446bf1ee5a02ab8457325343"
- `vc_data_root`: 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4

## 5.2 生成零知识证明

零知识证明生成功能将 VC 数据和验证条件作为输入，输出可用于验证的零知识证明。这个过程在用户自行部署的节点完成，不会泄露原始数据。生成的证明包含 Merkle 证明和数据根哈希，可用于后续的链上验证。

### 命令语法

```shell
storage-cli zk proof --input <INPUT_FILE_PATH> --threshold <BIRTH_DATE_THRESHOLD>
```

### 选项说明

| 选项                                 | 必填 | 说明                                                         |
| ------------------------------------ | ---- | ------------------------------------------------------------ |
| `--input <INPUT_FILE_PATH>`          | 是   | 包含输入值（VC、密钥、初始化向量、提交交易哈希）的文件路径，该文件在运行上传命令时自动生成 |
| `--threshold <BIRTH_DATE_THRESHOLD>` | 是   | 出生日期阈值，格式为年/月/日（如：20000101）                 |

### 示例

**示例 5.2**：为 Alice 生成零知识证明：

```shell
storage-cli zk proof \
   --input input_values.json \
   --threshold 20000101
```

其中 `input_values.json` 文件为[上传步骤](#51-vc加密上传)时的输出文件。

### 输出示例

```
✅ SUCCESS: == Successfully generated zk proof ==
   - VC Proof : c9c3da6512e9f20ef8dd07df85fd6831d6a8dc82f58f88a4d2f3163941345b9a5ab2e38a717ae078b1bb2c576878f3ed8f24161f4693ef2b0891ae9fb97d1103608d86f2697fc3336966effee5516460067463761cc5004ca2a113fbc0183099ca74cb260b27e0bc97bd15e9e1a8339b1e56d73d640d504dc65b94d55d087c28
   - Flow Root: 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4
```

## 5.3 零知识证明验证

零知识证明验证功能用于验证生成的证明是否有效。验证过程完全在链上进行，无需访问原始数据。验证器通过检查证明的密码学正确性和约束条件的满足情况，确定证明的有效性。

### 命令语法

```shell
storage-cli zk verify --proof <PROOF> --root <ROOT_HASH> --birth_threshold <BIRTH_DATE_THRESHOLD>
```

### 选项说明

| 选项                                       | 必填 | 说明                                                 |
| ------------------------------------------ | ---- | ---------------------------------------------------- |
| `--proof <PROOF>`                          | 是   | 零知识证明字符串                                     |
| `--root <ROOT_HASH>`                       | 是   | 去中心化文件系统的根哈希值                           |
| `--birth_threshold <BIRTH_DATE_THRESHOLD>` | 是   | 要验证的出生日期阈值，格式为年/月/日（如：20240101） |

### 示例

**示例 5.3**：用零知识证明的方式验证 Alice 的生日是否为 20000101：

```shell
storage-cli zk verify \
   --proof c9c3da6512e9f20ef8dd07df85fd6831d6a8dc82f58f88a4d2f3163941345b9a5ab2e38a717ae078b1bb2c576878f3ed8f24161f4693ef2b0891ae9fb97d1103608d86f2697fc3336966effee5516460067463761cc5004ca2a113fbc0183099ca74cb260b27e0bc97bd15e9e1a8339b1e56d73d640d504dc65b94d55d087c28 \
   --root 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4 \
   --birth_threshold 20000101
```

> **注意**：示例中使用的证明和根哈希值来自[生成零知识证明](#52-生成零知识证明)的输出结果。

# 附录A、可信证书技术原理与开发手册

## A.1 可信证书数据结构规范

### 原始证书结构

| 起始字节 | 终止字节 | 字段 | 格式规范 |
| --- | --- | --- | --- |
| 0 | 3 | 字面量`name` | ASCII 字符串 |
| 4 | 19 | 姓名 | UTF-8 编码，末尾补 0 |
| 20 | 22 | 字面量`age` | ASCII 字符串 |
| 23 | 23 | 年龄 | 无符号 8 位整数 |
| 24 | 28 | 字面量`birth` | ASCII 字符串 |
| 29 | 36 | 出生日期 | UTC 时间戳，无符号 64 位整数，小端序 |
| 37 | 39 | 字面量`edu` | ASCII 字符串 |
| 40 | 40 | 教育程度 | 无符号 8 位整数 |
| 41 | 46 | 字面量`serial` | ASCII 字符串 |
| 47 | 78 | 序列号 | 256 位哈希值 |

总字节数：79

### 证书加密流程

链上存储的可验证凭证（VC）为加密后的内容。加密过程所需输入包括：

- `key`：16 字节私钥
- `iv`：16 字节初始化向量，每次加密均不同
- `vc`：79 字节原始可验证凭证

加密步骤为

```bash
明文 = <79 字节 VC><32 字节 VC 的 Keccak 哈希值>
密文 = 使用 AES-CTR 模式，以 key、iv 和明文作为输入进行加密
```

### 分布式存储系统数据格式

上链数据结构如下，上链后由分布式存储系统确保数据可信性：

```bash
<16 字节初始化向量><111 字节密文><129 字节填充（全为 0）>
```

## A.2 自定义通用计算规范

应用开发者可通过 Circom 模板文件为链上可信证书提供可信数据查询的零知识证明，并实现自定义通用计算逻辑。

```
pragma circom 2.0.0;

template CustomCheck() {
    var name_len = 16;
    var serial_len = 32;
    var num_extensions = 16;

    signal input name[name_len];
    signal input age;
    signal input eduLevel;
    signal input serialNo[serial_len];
    signal input birthDateInt;

    signal input extensions[num_extensions];

    // 开发者自定义通用计算检查逻辑
    //
    // 示例：验证生日阈值
    // ```
    // signal birthdayOutput <== LessThan(64)([birthDateInt, extensions[0]]);
    // birthdayOutput === 1;
    // ```
}
```

## A.3 自定义通用计算构建与部署

### 环境准备

1. **Rust 环境安装**

   使用官方安装脚本进行安装:

   ```bash
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   ```

   配置当前 shell 环境以启用 Rust:

   ```bash
   source $HOME/.cargo/env
   ```

2. **Node.js 环境安装**

   安装 Node.js (要求版本 >= v18):

   ```bash
   curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
   sudo apt-get install -y nodejs
   ```

3. **Yarn 包管理器安装**

   通过 npm 安装 Yarn:

   ```bash
   npm install --global yarn
   ```

### 构建流程

1. **获取项目代码**

   从软件发布包中获取名称为 `0g-vc.tar.gz` 的项目代码，并切换至代码所在目录

   ```bash
   tar -xvf 0g-vc.tar.gz
   cd 0g-vc
   ```

2. **初始化项目**

   安装所有 JavaScript 依赖:

   ```bash
   yarn
   ```

3. **构建自定义逻辑**

   1. 将自定义文件保存至 `customized/<name>.circom`
   2. 执行以下命令:

   ```bash
   yarn build:custom <name>
   yarn setup <name>
   ```

**注意事项:**
- 构建过程可能需要 10 分钟或更长时间,具体取决于 CPU 性能
- 过程可能消耗超过 70GB 内存，请确保有足够的内存

## A.3 部署

构建完成后，在 `output` 目录下会生成以下文件:

- `<name>_js` 目录
- `<name>.pk` 文件
- `<name>.vk` 文件
- `<name>.r1cs` 文件

将上述文件替换零知识证明服务 `output` 目录下对应文件，重启服务。
