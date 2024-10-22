# 用户手册

## 文件加密
将文件加密，现在支持 aes 加密

### 命令格式
```sh
storage-cli file encrypt --source <source_file> --output <output_directory> --cipher <cipher_method> --password <cipher_password>
```

### 参数 
- `<source_file>` 需要加密的源文件的路径。
- `<output_directory>` 加密后文件的输出目录路径。
- `<cipher_method>` 加密方法，例如 aes。
- `<cipher_password>` 用于加密的密码。密码至少 16 位字符。

### 示例
示例 1： 将 example.txt 文件通过 aes 加密后输出到 ./encrypted 目录
```sh
storage-cli file encrypt --source example.txt --output ./encrypted --cipher aes --password mypassword123456
```

## 文件解密
将文件解密，支持与加密相同的加密方法。

```sh
storage-cli file decrypt --source <encrypted_file> --output <output_directory> --cipher <cipher_method> --password <cipher_password>
```

### 参数 

- `<encrypted_file>` 需要解密的加密文件的路径。
- `<output_directory>` 解密后文件的输出目录路径。
- `<cipher_method>` 解密方法，需与加密时使用的方法一致，例如 aes。
- `<cipher_password>` 用于解密的密码，需与加密时使用的密码一致。

### 示例
示例 2： 将 encrypted/example.enc 文件通过 aes 解密后输出到 ./decrypted 目录
```sh
storage-cli file decrypt --source encrypted/example.enc --output ./decrypted --cipher aes --password mypassword123456
```

## 上传文件
将文件上传到去中心化存储系统。可以选择加密上传。

```sh
storage-cli upload file --file <file_name> [--cipher <cipher_method> --password <cipher_password>]
```

### 参数 

- `<file_name>` 需要上传的文件的名称或路径。
- `[--cipher <cipher_method>]` 可选参数，指定加密方法，例如 aes。
- `[--password <cipher_password>]` 可选参数，指定用于加密的密码。

### 示例
示例 3： 将 example.txt 文件上传到去中心化存储系统
```sh
storage-cli upload file --file example.txt
```

示例 4： 将 example.txt 文件通过 aes 加密后上传到去中心化存储系统
```sh
storage-cli upload file --file example.txt --cipher aes --password mypassword123456
```

## 内容上传
将内容上传到去中心化存储系统。可以通过直接提供内容或指定文件路径来上传。

内容上传与[文件上传](#文件上传)有以下区别，[文件上传](#文件上传)时该文件没有所有者。而以该内容方式上传时，该文件拥有所有者。需要通过`--name`参数为内容设置名称，同时必须通过`--account`参数指定所有者。

所以如果需要将文件赋予所有权，则可以通过该方式上传。

> 需要注意的是当内容字节数不大于 1k 时，内容会直接存储到去中心化存储系统。当内容字节数大于 1k 时，则会以文件的形式上传。
而以文件的形式上传时则不能追加内容。

```sh
storage-cli upload content [--content <content> | --file <file_path>] [--account <account_address>] [--name <content_name>]
```

### 参数 
- `[--content <content>]` 需要上传的内容，直接以字符串形式提供。
- `[--file <file_path>]` 需要上传的内容所在文件的路径。
- `[--account <account_address>]` 可选参数，指定内容所有者的账户地址。
- `[--name <content_name>]` 可选参数，指定内容的名称。

### 示例
示例 5： 将 "Hello, World!" 以内容方式上传到去中心化存储系统，并指定所有者为 0x1234567890abcdef，内容名称为 Greeting
```sh
storage-cli upload content --content "Hello, World!" --account 0x1234567890abcdef --name "Greeting"
```

示例 6： 将 content.txt 文件的内容上传到去中心化存储系统，并指定所有者为 0x1234567890abcdef，内容名称为 FileContent
```sh
storage-cli upload content --file content.txt --account 0x1234567890abcdef --name "FileContent"
```

## 验证文件

验证文件是否与提供的文件匹配。当文件时加密上传时，该命令也需要指定相同的加密方法和密码。

```sh
storage-cli verify --file <file_path> [--cipher <cipher_method> --password <cipher_password>]
```

### 参数 

- `--file <file_path>` 需要验证的文件的路径。
- `[--cipher <cipher_method>]` 可选参数，指定该文件上传时使用的加密方法，当前支持 aes。
- `[--password <cipher_password>]` 可选参数，指定该文件上传时使用的密码。

### 示例
示例 7： 验证非加密上传的文件
```sh
storage-cli verify --file example.txt
```

示例 8： 验证加密上传的文件
```sh
storage-cli verify --file example.txt --cipher aes --password mypassword123456
```

## 追加内容
将数据追加到已上传的对应名称的内容后。

```sh
storage-cli append --name <content_name> [--data <append_data> | --file <file_path>] [--account <account_name>]
```

### 参数 

- `--name <content_name>` 要追加内容的名称。
- `[--data <append_data>]` 要追加的内容，直接以字符串形式提供。
- `[--file <file_path>]` 要追加的内容所在文件的路径。
- `[--account <account_name>]` 可选参数，指定用于追加内容的账户名称。

### 示例
示例 9： 将 "Hello again!" 追加到名称为 Greeting 的内容后
```sh
storage-cli append --name "Greeting" --data "Hello again!"
```

示例 10： 将 additional_content.txt 文件中的内容 追加到名称为 FileContent 的内容后
```sh
storage-cli append --name "FileContent" --file additional_content.txt
```

## 所有权查询

[内容上传](#内容上传) 时，会指定用户账户地址。该命令用于查询某个用户是否拥有该内容的所有权

### 命令格式

```sh
storage-cli owner content --account <account_address> --name <content_name>
```

### 参数

- `--account <account_address>` 需要检查的账户地址。
- `--name <content_name>` 需要检查的内容名称。

### 示例
示例 11： 检查账户 0x1234567890abcdef 是否拥有名称为 Greeting 的内容
```sh
storage-cli owner content --account 0x1234567890abcdef --name "Greeting"
```

## 所有权转移

将内容的所有权从一个账户转移到另一个账户。

### 命令
```sh
storage-cli owner transfer --from <current_owner> --to <target_owner> --name <content_name>
```

### 参数

- `--from <current_owner>` 当前内容所有者的账户地址。
- `--to <target_owner>` 目标所有者的账户地址。
- `--name <content_name>` 要转移所有权的内容名称。

### 示例
示例 12： 将名称为 Greeting 的内容的所有权从 0x1234567890abcdef 转移给 0xfedcba0987654321
```sh
storage-cli owner transfer --from 0x1234567890abcdef --to 0xfedcba0987654321 --name "Greeting"
```

## 生成零知识证明

生成零知识证明以验证特定条件下的数据。输出结果为 merkel 证明和 去中心化文件系统的 `root hash`，可用于零知识证明验证。

### 命令格式
```sh
storage-cli zk proof --vc <vc_string> --threshold <birth_date_threshold> --key <key> --iv <iv>
```

### 参数说明

- `--vc <vc_string>` 是以 JSON 格式提供的验证凭证字符串。
- `--threshold <birth_date_threshold>` 是出生日期的阈值，格式为年/月/日，例如 20240101。
- `--key <key>` 是用于加密的密钥。
- `--iv <iv>` 是初始化向量。

### 示例
示例 13： 为 Alice 生成零知识证明
```sh
storage-cli zk proof --vc '{"name": "Alice", "age": 25, "birth_date": "20000101", "edu_level": 4, "serial_no": "1234567890"}' --threshold 20240101 --key verysecretkey123
```
生成结果如下，`VC Proof`为 merkel 证明，`Flow Root`为去中心化文件系统的 `root hash`
```
✅ SUCCESS: == Successfully generated zk proof ==
    - VC Proof : c9c3da6512e9f20ef8dd07df85fd6831d6a8dc82f58f88a4d2f3163941345b9a5ab2e38a717ae078b1bb2c576878f3ed8f24161f4693ef2b0891ae9fb97d1103608d86f2697fc3336966effee5516460067463761cc5004ca2a113fbc0183099ca74cb260b27e0bc97bd15e9e1a8339b1e56d73d640d504dc65b94d55d087c28
    - Flow Root: 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4
```

## 零知识证明验证

验证零知识证明以确保数据的真实性和完整性。验证通过时结果为 `true`，否则为 `false`。

### 命令格式
```sh
storage-cli zk verify --proof <proof> --root <root_hash> --birth_threshold <birth_date_threshold>
```

### 参数说明

- `--proof <proof>` 零知识证明。
- `--root <root_hash>` 去中心化文件系统根哈希。
- `--birth_threshold <birth_date_threshold>` 要验证的出生日期。

### 示例
示例 14： 用零知识证明的方式验证 Alice 的生日是否为 20000101
```sh
storage-cli zk verify --proof c9c3da6512e9f20ef8dd07df85fd6831d6a8dc82f58f88a4d2f3163941345b9a5ab2e38a717ae078b1bb2c576878f3ed8f24161f4693ef2b0891ae9fb97d1103608d86f2697fc3336966effee5516460067463761cc5004ca2a113fbc0183099ca74cb260b27e0bc97bd15e9e1a8339b1e56d73d640d504dc65b94d55d087c28 --root 0x032303d969d3f271abfba865e159aa67e45ed406621c301e99c0643498eba7e4 --birth_threshold 20000101
```


