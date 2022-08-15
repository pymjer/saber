# Saber项目
使用Go写的小些工具，供日常使用

> 提示：使用命令`godoc -http=localhost:6060`启动文档服务器，打开浏览器输入`localhost:6060`查看详细文档

## 安装方法
* 安装好Go环境
* 进入工具目录，执行`go install`命令
*  执行命令

如，安装find_files工具：
```shell
$ cd findfiles
$ go install
$ findfiles
```

## find_files工具
在文件夹中查找包含指定内容的文件，输入想要查找的字符串和文件夹，输出包含字符串的文件列表
使用示例
```shell
$ findfiles
请输入想查询的字符串：
func
请输入想查询的目录：
.
开始查找目录[.]下文件中包含字符[func]的情况。包含字符串[func]的文件如下：
文件：[./0.go], 内容: func max(x, y int) int {
文件：[./1.go], 内容: func twoSum(nums []int, target int) []int {
文件：[./10.go], 内容: func isMatch(s
```

## hbaseutils工具
操作hbase库，目前支持如下功能：
根据前缀查询列值

```shell
$ hbaseutils
请输入zk地址：
zk.example.com
请输入表名和前缀名，用空格分隔：
test1201_0601 00
row: 000:info:name/testvalue:0
row: 001:info:name/testvalue:1
row: 002:info:name/testvalue:2
row: 003:info:name/testvalue:3
row: 004:info:name/testvalue:4
```

## bigmap工具
基于badger开发的kv文件数据库，用于存储kv数据
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