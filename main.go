package main

import (
	"flag"
	"fmt"
	"log"

	"prolion.top/saber/bigmap"
	"prolion.top/saber/findfiles"
	"prolion.top/saber/hbaseutils"
)

var cmd string

func main() {
	flag.StringVar(&cmd, "cmd", "findfiles", "命令的类型，当前支持1. bigmap,2. findfiles,3. hbaseutils,4. wiki, 可以通过序号或名称调用命令")
	flag.PrintDefaults()
	flag.Parse()
	fmt.Printf("cmd:%s \n", cmd)
	switch cmd {
	case "findfiles", "1":
		findfiles.FindFilesMain()
	case "bigmap", "2":
		bigmap.BigMapMain()
	case "hbaseutils", "3":
		hbaseutils.HBaseUtilMain()
	default:
		log.Fatalf("未知命令: %s", cmd)
	}
}
