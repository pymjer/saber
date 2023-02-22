// findfiles包用于查找包含指定字符的文件
package findfiles

import (
	"context"
	"fmt"
	"os"
	"strings"

	"prolion.top/saber/internal/base"
)

var CmdFindFiles = &base.Command{
	Run:       runFindFiles,
	UsageLine: "saber findfiles [flags] [dir]",
	Short:     "find files in a dir",
	Long: `
FindFiles prints the file in the dir.

FindFiles acceps zero, one, or two arguments.

Given no arguments, that is, when run as 

	saber findfiles

it prints the current dir all files.

The -c flag is content
The -f flag is file suffix

Examples:
	saber findfiles -c content
		Show files in current dir contain content
	saber findfiles -c content -f filter
		Show files 
	`,
}

var (
	content string
	filter  string
	dirPath = "."
)

func init() {
	CmdFindFiles.Flag.StringVar(&content, "c", "", "查询内容")
	CmdFindFiles.Flag.StringVar(&filter, "f", "", "查询文件后缀")
}

func runFindFiles(ctx context.Context, cmd *base.Command, args []string) {
	if len(args) >= 1 {
		dirPath = args[len(args)-1]
	}
	FindFiles(content, dirPath, filter)
}

func FindFilesMain() {
	fmt.Println("请输入想查询的字符串：")
	fmt.Scanln(&content)
	fmt.Println("请输入想查询的目录(默认是当前目录)：")
	fmt.Scanln(&dirPath)
	fmt.Println("输入想要查找的文件后缀，默认所有文本文件：")
	fmt.Scanln(&filter)
}

// 在给定目录中查找包含指定内容的文件，可以根据文件名过滤
func FindFiles(content, dirPath, filter string) {
	fmt.Printf("开始查找目录[%s]下文件中包含字符[%s]的情况...\n", dirPath, content)
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
	for i, f := range containFiles {
		fmt.Printf("%d [%s], 内容: %s \n", i, f[0], f[1])
	}

	if len(notContainFiles) > 1 {
		fmt.Printf("不包含字符串[%s]的文件如下：\n", content)
		for _, f := range notContainFiles {
			fmt.Printf("文件：[%s]\n", f[0])
		}
	} else {
		fmt.Printf("没有不包含字符串[%s]的文件\n", content)
	}

}

func GetFilesAndDirs(dirPath string, filter string) (files []string, dirs []string, err error) {
	dir, err := os.ReadDir(dirPath)
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
	buf, err := os.ReadFile(filePath)
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
