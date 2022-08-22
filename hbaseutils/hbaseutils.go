// hbaseutils包实现了操作hbase库的常用工具
package hbaseutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tsuna/gohbase/pb"
	"prolion.top/saber/internal/base"
)

var CmdHBaseUtil = &base.Command{
	UsageLine: "saber hutil [-l] [-s] [-d] [-n]",
	Short:     "hbase utils",
	Long: `
saber hutil	-h zk.example.com -l

The -h flag is zookeeper 

The -l flag show tables

The -s flag show tables values

The -d flag delete tables

The -n flag limit the number of rows

Examples:
	saber hutil -l
		Show all tables
	saber hutil -l [table_prefix]
		Show all tables with prefix, prefix is reg, eg "test_.*"
	saber hutil -s table_name [keyprefix]
		Show tables values
	saber hutil -s -n 20 table_name [keyprefix]
		Show tables values with 20 rows
	saber hutil -d table_name1 table_name2
		Delete tables
	`,
}

var (
	host string
	list bool
	show bool
	rows int
	del  bool
)

func init() {
	CmdHBaseUtil.Run = runHUtil

	CmdHBaseUtil.Flag.StringVar(&host, "h", "", "zk host address")
	CmdHBaseUtil.Flag.BoolVar(&list, "l", false, "list tables")
	CmdHBaseUtil.Flag.BoolVar(&show, "s", false, "show tables values")
	CmdHBaseUtil.Flag.BoolVar(&del, "d", false, "delete tables")
	CmdHBaseUtil.Flag.IntVar(&rows, "n", 10, "limit the number of rows, default is 10")
}

func runHUtil(ctx context.Context, cmd *base.Command, args []string) {
	HBaseUtilMain(args)
}

func HBaseUtilMain(args []string) {
	if host == "" {
		host = os.Getenv("ZK")
		if host == "" {
			base.Fatalf("unable to locate zookeeper host, use -h flag or saber env -w set zk")
		}
	}
	fmt.Printf("connnect to zookeeper host: %s\n", host)

	if list && show {
		base.Fatalf("saber hutil: cannot use -l with -s")
	}

	u := NewHBaseUtils(host)

	if list {
		if len(args) < 1 {
			PrintTables(u.ListTables())
			return
		}

		tabPrefix := args[len(args)-1]
		PrintTables(u.FindTables(tabPrefix))
		return
	}

	if show {
		if len(args) < 1 {
			base.Fatalf("saber hutil tablename: no table name")
		}

		table := args[0]
		if len(args) == 1 {
			cells := u.ScanTable(table, rows)
			PrintCells(cells)
			return
		}

		prefix := args[1]
		cells, err := u.ScanWithPrefixFilter(table, prefix, rows)
		if err != nil {
			log.Fatal(err)
		}
		PrintCells(cells)
		return
	}

	if del {
		if len(args) < 1 {
			base.Fatalf("saber hutil -d tablename: no table name")
		}
		fmt.Printf("确认要删除表[%s]吗？y/n, 操作不可恢复 \n", strings.Join(args, ","))
		var ans string
		fmt.Scanln(&ans)
		if ans == "y" || ans == "yes" {
			for _, tableName := range args {
				fmt.Printf("begin delete table:[%s] \n", tableName)
				err := u.DisableTable(tableName)
				if err != nil {
					if !strings.Contains(err.Error(), "TableNotEnabledException") {
						base.Fatalf("disable table[%s] error: %v", tableName, err)
					} else {
						fmt.Printf("table[%s] is disabled, continue.\n", tableName)
					}
				}
				err = u.DeleteTable(tableName)
				if err != nil {
					base.Fatalf("delete table[%s] error: %v", tableName, err)
				}
				fmt.Printf("table[%s] delete success!\n", tableName)
			}
			return
		}
		fmt.Printf("no table delete. \n")
		return
	}

	PrintTables(u.ListTables())
}

func PrintCells(cells []*Cell) {
	for _, cell := range cells {
		fmt.Printf("row: %s\n", cell)
	}
}

func PrintTables(names []*pb.TableName) {
	for i, n := range names {
		fmt.Printf("%d: %s \n", i, string(n.Qualifier))
	}
}
