# 用户手册

# 功能介绍
## 具备数据的存证、加密、验真、分享等治理能力

### 文件加密
将文件加密，现在支持 aes 加密

```sh
storage-cli file encrypt --source <source_file> --output <output_directory> --cipher <cipher_method> --password <cipher_password>
```

**参数说明**: 
- `<source_file>` 是需要加密的源文件的路径。
- `<output_directory>` 是加密后文件的输出目录路径。
- `<cipher_method>` 是加密方法，例如 aes。
- `<cipher_password>` 是用于加密的密码。密码至少 16 位字符。

**示例**:
```sh
storage-cli file encrypt --source example.txt --output ./encrypted --cipher aes --password mypassword123456
```

### 文件解密
将文件解密，支持与加密相同的加密方法。

```sh
storage-cli file decrypt --source <encrypted_file> --output <output_directory> --cipher <cipher_method> --password <cipher_password>
```

**参数说明**: 

- `<encrypted_file>` 是需要解密的加密文件的路径。
- `<output_directory>` 是解密后文件的输出目录路径。
- `<cipher_method>` 是解密方法，需与加密时使用的方法一致，例如 aes。
- `<cipher_password>` 是用于解密的密码，需与加密时使用的密码一致。

**示例**:
```sh
storage-cli file decrypt --source encrypted/example.enc --output ./decrypted --cipher aes --password mypassword123456
```

### 上传文件
将文件上传到去中心化存储系统。可以选择加密上传。

```sh
storage-cli upload file --file <file_name> [--cipher <cipher_method> --password <cipher_password>]
```

**参数说明**: 

- `<file_name>` 是需要上传的文件的名称或路径。
- `[--cipher <cipher_method>]` 可选参数，指定加密方法，例如 aes。
- `[--password <cipher_password>]` 可选参数，指定用于加密的密码。

**示例**:
```sh
storage-cli upload file --file example.txt
```

如果需要加密上传：
```sh
storage-cli upload file --file example.txt --cipher aes --password mypassword123456
```

### 上传内容
将内容上传到去中心化存储系统。可以通过直接提供内容或指定文件路径来上传。

```sh
storage-cli upload content [--content <content> | --file <file_path>] [--account <account_address>] [--name <content_name>]
```


**参数说明**: 

- `[--content <content>]` 是需要上传的内容，直接以字符串形式提供。
- `[--file <file_path>]` 是需要上传的内容所在文件的路径。
- `[--account <account_address>]` 可选参数，指定内容所有者的账户地址。
- `[--name <content_name>]` 可选参数，指定内容的名称。

**示例**:
```sh
storage-cli upload content --content "Hello, World!" --account 0x1234567890abcdef --name "Greeting"
```


通过文件上传内容：
```sh
storage-cli upload content --file content.txt --account 0x1234567890abcdef --name "FileContent"
```

### 验证内容
验证文件是否与提供的文件匹配。可以选择使用加密方法进行验证。

```sh
storage-cli verify --file <file_path> [--cipher <cipher_method> --password <cipher_password>]
```


**参数说明**: 

- `--file <file_path>` 是需要验证的文件的路径。
- `[--cipher <cipher_method>]` 可选参数，指定用于验证的加密方法，例如 aes。
- `[--password <cipher_password>]` 可选参数，指定用于验证的密码。

**示例**:
```sh
storage-cli verify --file example.txt
```


如果需要使用加密方法进行验证：
```sh
storage-cli verify --file example.txt --cipher aes --password mypassword123456
```
## 支持数据上链和自动重传

上传相关的命令均支持上链和自动重传

**示例**
```sh
storage-cli upload file --file example.txt
```

如执行上述命令时，将会自动重试上传

## 支持数据的增量更新

### 追加内容
将数据追加到指定名称的已上传内容中。

```sh
storage-cli append --name <content_name> [--data <append_data> | --file <file_path>] [--account <account_name>]
```


**参数说明**: 

- `--name <content_name>` 是要追加内容的名称。
- `[--data <append_data>]` 是要追加的内容，直接以字符串形式提供。
- `[--file <file_path>]` 是要追加的内容所在文件的路径。
- `[--account <account_name>]` 可选参数，指定用于追加内容的账户名称。

**示例**:
```sh
storage-cli append --name "Greeting" --data "Hello again!"
```

通过文件追加内容：
```sh
storage-cli append --name "FileContent" --file additional_content.txt
```

## 数据的灵活查询
## 数据的生命周期全程确认和确权功能 （权限转移）
## 为可信数据查询和通用的计算结果提供零知识证明
