# 9.2 日需求整理
## 数据类型
原始数据、扩展数据
- 原始数据被 rpc 支持, id 是 hash
- 扩展数据具有 Append, Owner transfer 的权限，id 是 name

### 扩展数据存储格式 （kv）
类型可为： 文本(text) 或 指针(pointer)
```
<name>:size
<name>:type = "text" | "pointer"
<name>:_ = "nft token id"
<name>:<number> key, value; pointer类型时value是原始数据的指针；text类型时value是一行的文本内容（行是按字符串长度截取的，比如 1k 一行）
```
## 命令所用数据类型
- append      扩展数据实现(text) Append content to specified file  
    - 只支持 by name 且 类型为 text 的文件
- batch       扩展数据实现(text) Batch operations
- upload      
    - by hash: 原始数据实现
    - by name: 根据长度自适应两种扩展数据(pointer, text) 
        - 类型为 pointer 的不支持 append
- download    
    - by hash: 原始数据实现
    - by name: 自适应两种扩展数据(pointer, text) 
- file        本地操作
- owner       
    - by hash: 无需支持
    - by name: 扩展数据实现(text)
- zk prove    原始数据
- zk verify   原始数据
- verify      原始数据实现（课题二）
- template   （课题二）
