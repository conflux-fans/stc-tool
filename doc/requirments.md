•10000条/秒

	•KV BATCH，网页 SHOW NUMBER / 命令行 SHOW NUMBER / SEQUENCE NUMBER
	
	- 命令：batch upload
		- 参数：条数
	- 步骤
		- 批量创建随机文本
		- 批量上传
		- 输出log
		- 输出结果，统计时间


•具备数据的存证、加密、验真、分享等治理能力

	•存证调 API, 加密 SDK, 验真 SDK，分享 SCAN 链接？
	
	- 命令：upload
      	- 子命令：append
		- 参数
    		- 名称
			- 文件路径
			- 加密方法
			- 密码
			
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


•具备数据的生命周期全程确认和确权功能

	PLANA: NFT 合约开发，NFT 流转记录。
	用户-＞NFT 合约-＞存储合约
	PLANB: 使用KV 的OWNER, OWNER 转移上链吗？

	- 命令：transferowner
		- 参数 
			- streamID
			- to


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

•支持数据上链和自动重传;

	GO CLIENT，己经有/

•支持数据的增量更新；
	通过目前的KVDB 实现，SDK

	- 命令：append
		- 参数
			- tag


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
