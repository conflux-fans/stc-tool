# Storage CLI 用户手册

## 1. 文件加密

对文件进行加密操作。目前支持 AES 加密算法。

### 命令语法

```shell
storage-cli file encrypt --source <SOURCE_FILE> --output <OUTPUT_DIR> --cipher <CIPHER_METHOD> --password <PASSWORD>
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--source <SOURCE_FILE>` | 是 | 指定需要加密的源文件路径 |
| `--output <OUTPUT_DIR>` | 是 | 指定加密后文件的输出目录 |
| `--cipher <CIPHER_METHOD>` | 是 | 指定加密算法，当前支持：`aes` |
| `--password <PASSWORD>` | 是 | 加密密码，至少16位字符 |

### 示例

**示例 1**：将 `example.txt` 文件通过 AES 加密后输出到 `./encrypted` 目录：

```shell
storage-cli file encrypt \
   --source example.txt \
   --output ./encrypted \
   --cipher aes \
   --password mypassword123456
```

## 2. 文件解密

将加密文件解密，支持与加密相同的加密方法。

### 命令语法

```shell
storage-cli file decrypt --source <ENCRYPTED_FILE> --output <OUTPUT_DIR> --cipher <CIPHER_METHOD> --password <PASSWORD>
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--source <ENCRYPTED_FILE>` | 是 | 指定需要解密的加密文件路径 |
| `--output <OUTPUT_DIR>` | 是 | 指定解密后文件的输出目录 |
| `--cipher <CIPHER_METHOD>` | 是 | 指定解密算法，需与加密时使用的方法一致，当前支持：`aes` |
| `--password <PASSWORD>` | 是 | 解密密码，需与加密时使用的密码一致 |

### 示例

**示例 2**：将 `encrypted/example.enc` 文件通过 AES 解密后输出到 `./decrypted` 目录：

```shell
storage-cli file decrypt \
   --source encrypted/example.enc \
   --output ./decrypted \
   --cipher aes \
   --password mypassword123456
```

## 3. 上传文件

将文件上传到去中心化存储系统。支持可选的加密上传。

### 命令语法

```shell
storage-cli upload file --file <FILE_NAME> [--cipher <CIPHER_METHOD> --password <PASSWORD>]
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--file <FILE_NAME>` | 是 | 指定需要上传的文件名称或路径 |
| `--cipher <CIPHER_METHOD>` | 否 | 指定加密算法，当前支持：`aes` |
| `--password <PASSWORD>` | 否 | 加密密码，当指定 `--cipher` 时必填 |

### 示例

**示例 3**：将 `example.txt` 文件上传到去中心化存储系统：

```shell
storage-cli upload file --file example.txt
```

**示例 4**：将 `example.txt` 文件通过 AES 加密后上传到去中心化存储系统：

```shell
storage-cli upload file \
   --file example.txt \
   --cipher aes \
   --password mypassword123456
```

## 4. 内容上传

将内容上传到去中心化存储系统。支持通过直接提供内容或指定文件路径两种方式上传。

> **注意**：与[文件上传](#3-上传文件)的区别在于，文件上传时该文件没有所有者，而内容上传时文件拥有所有者。需要通过 `--name` 参数为内容设置名称，同时必须通过 `--account` 参数指定所有者。
>
> 当内容字节数不大于 1k 时，内容会直接存储到去中心化存储系统。当内容字节数大于 1k 时，则会以文件的形式上传。注意以文件形式上传时不能追加内容。

### 命令语法

```shell
storage-cli upload content --account <ACCOUNT_ADDRESS> --name <CONTENT_NAME> (--content <CONTENT> | --file <FILE_PATH>)
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--content <CONTENT>` | 与 `file` 二选一 | 需要上传的内容，直接以字符串形式提供 |
| `--file <FILE_PATH>` | 与 `content` 二选一 | 需要上传的内容所在文件的路径 |
| `--account <ACCOUNT_ADDRESS>` | 是 | 指定内容所有者的账户地址 |
| `--name <CONTENT_NAME>` | 是 | 指定内容的名称 |

### 示例

**示例 5**：将文本内容 "Hello, World!" 上传到去中心化存储系统：

```shell
storage-cli upload content \
   --content "Hello, World!" \
   --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
   --name "Greeting"
```

**示例 6**：将文件 `content.txt` 的内容上传到去中心化存储系统：

```shell
storage-cli upload content \
   --file content.txt \
   --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
   --name "FileContent"
```

## 5. 验证文件

验证文件是否与提供的文件匹配。当文件是加密上传时，该命令也需要指定相同的加密方法和密码。

### 命令语法

```shell
storage-cli verify --file <FILE_PATH> [--cipher <CIPHER_METHOD> --password <PASSWORD>]
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--file <FILE_PATH>` | 是 | 需要验证的文件路径 |
| `--cipher <CIPHER_METHOD>` | 否 | 文件上传时使用的加密方法，当前支持：`aes` |
| `--password <PASSWORD>` | 否 | 文件上传时使用的密码，当指定 `--cipher` 时必填 |

### 示例

**示例 7**：验证非加密上传的文件：

```shell
storage-cli verify --file example.txt
```

**示例 8**：验证加密上传的文件：

```shell
storage-cli verify \
   --file example.txt \
   --cipher aes \
   --password mypassword123456
```

## 6. 追加内容

将数据追加到已上传的对应名称的内容后。

### 命令语法

```shell
storage-cli append --name <CONTENT_NAME> (--data <APPEND_DATA> | --file <FILE_PATH>) [--account <ACCOUNT_NAME>]
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--name <CONTENT_NAME>` | 是 | 要追加内容的名称 |
| `--data <APPEND_DATA>` | 与 `file` 二选一 | 要追加的内容，直接以字符串形式提供 |
| `--file <FILE_PATH>` | 与 `data` 二选一 | 要追加的内容所在文件的路径 |
| `--account <ACCOUNT_NAME>` | 否 | 指定用于追加内容的账户名称 |

### 示例

**示例 9**：将文本内容 "Hello again!" 追加到名称为 Greeting 的内容后：

```shell
storage-cli append \
   --name "Greeting" \
   --data "Hello again!"
```

**示例 10**：将文件 `additional_content.txt` 中的内容追加到名称为 FileContent 的内容后：

```shell
storage-cli append \
   --name "FileContent" \
   --file additional_content.txt
```

## 7. 所有权查询

查询指定用户是否拥有某个内容的所有权。该命令用于查询[内容上传](#4-内容上传)时指定的用户账户地址。

### 命令语法

```shell
storage-cli owner content --account <ACCOUNT_ADDRESS> --name <CONTENT_NAME>
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--account <ACCOUNT_ADDRESS>` | 是 | 需要检查的账户地址 |
| `--name <CONTENT_NAME>` | 是 | 需要检查的内容名称 |

### 示例

**示例 11**：检查账户 `0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE` 是否拥有名称为 Greeting 的内容：

```shell
storage-cli owner content \
   --account 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
   --name "Greeting"
```

## 8. 所有权转移

将内容的所有权从一个账户转移到另一个账户。

### 命令语法

```shell
storage-cli owner transfer --from <CURRENT_OWNER> --to <TARGET_OWNER> --name <CONTENT_NAME>
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--from <CURRENT_OWNER>` | 是 | 当前内容所有者的账户地址 |
| `--to <TARGET_OWNER>` | 是 | 目标所有者的账户地址 |
| `--name <CONTENT_NAME>` | 是 | 要转移所有权的内容名称 |

### 示例

**示例 12**：将名称为 Greeting 的内容的所有权从 `0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE` 转移给 `0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b`：

```shell
storage-cli owner transfer \
   --from 0x26154DF6A79a6C241b46545D672A3Ba6AE8813bE \
   --to 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b \
   --name "Greeting"
```

## 9. 生成零知识证明

生成零知识证明以验证特定条件下的数据。输出结果包括 Merkle 证明和去中心化文件系统的 `root hash`，可用于零知识证明验证。

### 命令语法

```shell
storage-cli zk proof --vc <VC_STRING> --threshold <BIRTH_DATE_THRESHOLD> --key <KEY> --iv <IV>
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--vc <VC_STRING>` | 是 | 验证凭证字符串，需以 JSON 格式提供 |
| `--threshold <BIRTH_DATE_THRESHOLD>` | 是 | 出生日期阈值，格式为年/月/日（如：20240101）|
| `--key <KEY>` | 是 | 用于加密的密钥 |
| `--iv <IV>` | 是 | 初始化向量 |

### 示例

**示例 13**：为 Alice 生成零知识证明：

```shell
storage-cli zk proof \
   --vc '{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}' \
   --threshold 20240101 \
   --key verysecretkey123
```

### 输出示例

```
✅ SUCCESS: == Successfully generated zk proof ==
   - VC Proof : c9c3da6512e9f20ef8dd07df85fd6831d6a8dc82f58f88a4d2f3163941345b9a5ab2e38a717ae078b1bb2c576878f3ed8f24161f4693ef2b0891ae9fb97d1103608d86f2697fc3336966effee5516460067463761cc5004ca2a113fbc0183099ca74cb260b27e0bc97bd15e9e1a8339b1e56d73d640d504dc65b94d55d087c28
   - Flow Root: 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4
```

## 10. 零知识证明验证

验证零知识证明以确保数据的真实性和完整性。验证通过时结果为 `true`，否则为 `false`。

### 命令语法

```shell
storage-cli zk verify --proof <PROOF> --root <ROOT_HASH> --birth_threshold <BIRTH_DATE_THRESHOLD>
```

### 选项说明

| 选项 | 必填 | 说明 |
|------|------|------|
| `--proof <PROOF>` | 是 | 零知识证明字符串 |
| `--root <ROOT_HASH>` | 是 | 去中心化文件系统的根哈希值 |
| `--birth_threshold <BIRTH_DATE_THRESHOLD>` | 是 | 要验证的出生日期阈值，格式为年/月/日（如：20240101）|

### 示例

**示例 14**：用零知识证明的方式验证 Alice 的生日是否为 20000101：

```shell
storage-cli zk verify \
   --proof c9c3da6512e9f20ef8dd07df85fd6831d6a8dc82f58f88a4d2f3163941345b9a5ab2e38a717ae078b1bb2c576878f3ed8f24161f4693ef2b0891ae9fb97d1103608d86f2697fc3336966effee5516460067463761cc5004ca2a113fbc0183099ca74cb260b27e0bc97bd15e9e1a8339b1e56d73d640d504dc65b94d55d087c28 \
   --root 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4 \
   --birth_threshold 20000101
```

> **注意**：示例中使用的证明和根哈希值来自[生成零知识证明](#9-生成零知识证明)的输出结果。
