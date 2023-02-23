# Saber项目
使用Go写的小些工具，供日常使用

> 提示：使用命令`godoc -http=localhost:6060`启动文档服务器，打开浏览器输入`localhost:6060`查看详细文档

## 安装方法
* 安装好Go环境
* 构建代码
* 进入工具目录，执行`go install`命令
* 执行命令

跨平台构建方法

```shell
$ env GOOS=target-OS GOARCH=target-architecture go build package-import-path
```
示例
```
// mac平台 x86核心
$ env GOOS=darwin GOARCH=amd64 go build .
// mac平台 arm核心
$ env GOOS=darwin GOARCH=arm go build .
// window平台
$ env GOOS=windows GOARCH=amd64 go build .
// Linux平台
$ env GOOS=linux GOARCH=amd64 go build .
```

## 使用方法
执行`saber help`查看当前支持的命令
当前支持四种工具
* simulate 一个模拟数据的小工具
* findfiles
* bigmap 基于badger开发的kv文件数据库，用于存储kv数据
* hutil
* env
* wiki 一个简单的wiki页面

使用`saber help <command>`查看某个工具的使用方法

## simulate 工具
Simulate 是一个模拟数据的小工具

使用方法如下：
```
$ saber sdata -c id@int,name,age@int,amount@float -t json users.json
```
帮助文档
```
Column is a comma-separated field, the field defaults to string type, if the type is int or float, use the field@type

The -c flag is columns
The -n flag is the row numbers, default 100
The -t flag is file type, support csv,json, default is csv
The -f flag is the result file name. default is {current_time_stamp}.json
```

## hutil 工具
hutil 用户操作hbase表
```
## 查看test前缀表 
$ saber hutil -h zk.example.com -l "test.*"

## 删除表
$ saber hutil -h zk.example.com -d table1 table2

## 如果想省略-h的配置，可以设置为环境变量
$ saber env -w ZK=zk.example.com
```

## bigmap工具
使用示例：
```
$ saber bigmap data
当前数据目录为: data 
badger 2022/08/05 15:03:38 INFO: All 0 tables opened in 0s
> h
quit,q 退出
help,h 查看帮助
list,l 查看所有key
seek,s 根据指定前缀查找所有kv对
get key 获取指定key的值
put key value 添加/跟新指定key的值
del key 删除指定key的值
> put aa 11
> put bb 22
> get aa
11
> get bb
22
> q
```