package main

import (
	"fmt"
	"log"
)

func main() {
	// 查找表指定前缀的值
	var host, table, prefix string
	fmt.Println("请输入zk地址：")
	fmt.Scanln(&host)
	fmt.Println("请输入表名和前缀名，用空格分隔：")
	fmt.Scanln(&table, &prefix)
	u := NewHBaseUtils(host)
	cells, err := u.ScanWithPrefixFilter(table, prefix)
	if err != nil {
		log.Fatal(err)
	}
	for _, cell := range cells {
		fmt.Printf("row: %s\n", cell)
	}
}
