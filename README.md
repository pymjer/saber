# Saber项目
使用Go写的小些工具，供日常使用

> 提示：使用命令`godoc -http=localhost:6060`启动文档服务器，打开浏览器输入`localhost:6060`查看详细文档

## 安装方法
* 安装好Go环境
* 进入工具目录，执行`go install`命令
* 执行命令

## 使用方法
执行`saber help`查看当前支持的命令
当前支持四种工具
* findfiles
* bigmap 基于badger开发的kv文件数据库，用于存储kv数据
* hutil
* env

使用`saber help <command>`查看某个工具的使用方法


## bigmap工具
使用示例：
```
$ bigmap data
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