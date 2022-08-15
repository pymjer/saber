// hbaseutils包实现了操作hbase库的常用工具
package hbaseutils

import (
	"context"
	"fmt"
	"log"

	"prolion.top/saber/internal/base"
)

var CmdHBaseUtil = &base.Command{
	UsageLine: "saber hutil [flags]",
	Short:     "hbase utils",
	Long: `
saber hutil	-h zk.example.com -p keyprefix table_name

The -h flag is zookeeper address

The -p is key prefix
	`,
}

var (
	host   string
	prefix string
	table  string
)

func init() {
	CmdHBaseUtil.Run = runHUtil

	CmdHBaseUtil.Flag.StringVar(&host, "h", "", "zk host address")
	CmdHBaseUtil.Flag.StringVar(&prefix, "p", "", "key prefix")
}

func runHUtil(ctx context.Context, cmd *base.Command, args []string) {
	if len(args) < 1 {
		base.Fatalf("saber hutil tablename: no table name")
	}
	table = args[len(args)-1]
	HBaseUtilMain()
}

func HBaseUtilMain() {
	u := NewHBaseUtils(host)
	cells, err := u.ScanWithPrefixFilter(table, prefix)
	if err != nil {
		log.Fatal(err)
	}
	for _, cell := range cells {
		fmt.Printf("row: %s\n", cell)
	}
}
