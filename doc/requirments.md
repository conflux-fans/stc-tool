•10000条/秒

	•KV BATCH，网页 SHOW NUMBER / 命令行 SHOW NUMBER / SEQUENCE NUMBER
	
	- 命令：batch upload
		- 参数：条数
	- 步骤
		- 批量创建随机文本
		- 批量上传
		- 输出log
		- 输出结果，统计时间

```sh
go run . batch upload -c 100000
```


•具备数据的存证、加密、验真、分享等治理能力

	•存证调 API, 加密 SDK, 验真 SDK，分享 SCAN 链接？
	
	- 命令：upload 
      	- 子命令：file
		- 参数
			- 文件路径
			- 加密方法
			- 密码
    	- 子命令：data
    	- 参数
    		- 名称
    		- 内容
			- 或文件路径
	- 命令：append
    	- 参数
        	- 名称
        	- 内容
			- 或文件路径
			
	- 命令：verify
		- 参数
			- 文件路径
			- 加密方法
			- 密码
    	- 步骤
        	- 如果指定加密，根据文件生成加密文件
        	- 计算文件root哈希
        	- 查询文件是否存在
        	- 需要查所有者？

```sh
go run . upload content --name data2 --content hello --account 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b

go run . upload file --file ./main.go --cipher aes --password 1234567812345678

go run . append --name data2 --data " world" --account 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b

go run . verify --file ./main.go --cipher aes --password 1234567812345678
```

•具备数据的生命周期全程确认和确权功能

	PLANA: NFT 合约开发，NFT 流转记录。
	用户-＞NFT 合约-＞存储合约
	PLANB: 使用KV 的OWNER, OWNER 转移上链吗？

	- 命令：transferowner
		- 参数 
			- streamID
			- to

```sh
go run . owner content --account 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b --name data2 
go run . owner transfer --name data1 --from 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b --to 0xe61646FD48adF644404f373D984B14C877957F7c 
```

•支持上链存证模板自定义；

	• 文本模板，存合约里，SDK包一下
	- 命令：template
		- 子命令：创建模版
			- 参数
    			- 模版名称
				- 字段数组
		- 子命令：列出模版

		- 子命令：下载模版
			- 参数
				- 模版名称
			- 步骤
				- 从合约读取对应字段数组创建csv文件

```sh
 go run . template list  

 go run . template create --name test4 --keys name,age    

 go run . template download --name test4
```

•支持数据上链和自动重传;

	GO CLIENT，己经有/

•支持数据的增量更新；
	通过目前的KVDB 实现，SDK

	- 命令：append
		- 参数
			- tag

```sh
go run . append --name data1 --data " world" --account 0xd68D7A9639FaaDed2a6002562178502fA3b3Af9b
```


•支持数据的灵活查询；

	PLANA: MERKLE PROOF
	PLANB: KV查询

	- 命令：query
		- 参数
			- key
			- tag

•支持异常数据的链上验真

•支持为可信数据查询和通用的计算结果提供零知识证明



stream 每次都是全量上传
upload: 1.上传文件 2.创建stream，FILE0=>ROOTHASH
append: 1.上传文件 2.在stream上追加 FILE1=>ROOTHASH


# 3.6会议
1. 批量上传现在以发送完交易就算完成，需要等到文件上传完成吗？
   1. 用kv做，k：ID，value：文本
   2. 一个batch 500条，20个batch
   3. 等到文件上传完成
   4. 计算时间
   5. 批量查询是否上链
2. 验真现在是只根据root查文件是否存在，是否还需其它步骤？ 辰星研究
3. 文件追加使用stream，stream只支持配置文件中的？如果这样，也可以用同一个stream管理所有文件，会有什么问题。 用kv做；value就是行
4. 待讨论方案
   1. 支持异常数据的链上验真
   2. 支持为可信数据查询和通用的计算结果提供零知识证明

# 5.9会议

1. 优化日志
2. bug 修复
3. 确定卡住的问题是否能修复，是否需要现场部署环境