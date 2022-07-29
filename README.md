# Saber项目
使用Go写的小些工具，供日常使用

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