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

## 用户手册
请见[用户手册](./docs/user-manual.md)