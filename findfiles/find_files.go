package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func FindFiles(content, dirPath, filter string) {
	fmt.Printf("开始查找目录[%s]下文件中包含字符[%s]的情况。", dirPath, content)
	files, dirs, _ := GetFilesAndDirs(dirPath, filter)
	for _, dir := range dirs {
		fmt.Printf("读取到文件夹：[%s]\n", dir)
	}
	var containFiles, notContainFiles [][]string
	for _, file := range files {
		contain, line := IsFileContainStr(file, content)
		if contain {
			containFiles = append(containFiles, []string{file, line})
		} else {
			notContainFiles = append(notContainFiles, []string{file})
		}
	}
	fmt.Printf("包含字符串[%s]的文件如下：\n", content)
	for _, f := range containFiles {
		fmt.Printf("文件：[%s], 内容: %s \n", f[0], f[1])
	}

	fmt.Printf("不包含字符串[%s]的文件如下：\n", content)
	for _, f := range notContainFiles {
		fmt.Printf("文件：[%s]\n", f[0])
	}
}

func GetFilesAndDirs(dirPath string, filter string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPath+PthSep+fi.Name())
			GetFilesAndDirs(dirPath+PthSep+fi.Name(), filter)
		} else {
			// 过滤文件类型
			fileName := dirPath + PthSep + fi.Name()
			if filter != "" {
				ok := strings.HasSuffix(fi.Name(), filter)
				if ok {
					files = append(files, fileName)
				}
			} else {
				files = append(files, fileName)
			}
		}
	}
	return
}

func IsFileContainStr(filePath string, str string) (bool, string) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
		panic(err.Error())
	}
	content := string(buf)
	idx := strings.Index(content, str)
	if idx < 0 {
		return false, ""
	}
	end := strings.Index(content[idx:], "\n")
	line := content[idx : idx+end]
	return true, line
}

func main() {
	var content, dirPath, filter string
	fmt.Println("请输入想查询的字符串：")
	fmt.Scanln(&content)
	fmt.Println("请输入想查询的目录：")
	fmt.Scanln(&dirPath)
	fmt.Println("输入想要查找的文件后缀，默认所有文本文件：")
	fmt.Scanln(&filter)
	FindFiles(content, dirPath, filter)
}
