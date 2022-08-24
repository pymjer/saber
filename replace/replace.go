package replace

import (
	"context"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"strings"

	"prolion.top/saber/internal/base"
)

var CmdReplaceFile = &base.Command{
	Run:       runReplaceFile,
	UsageLine: "saber replace [src_str] [diststr]",
	Short:     "replace file",
	Long: `
Replace replace file.

Examples:
	saber replace aa bb
		replace aa to bb
	`,
}

//go:embed tpl.mht
var tpl string

func runReplaceFile(ctx context.Context, cmd *base.Command, args []string) {
	var name string
	var phoneNum string
	fmt.Printf("请输入用户姓名： \n")
	fmt.Scanln(&name)

	fmt.Printf("请输入用户电话号码： \n")
	fmt.Scanln(&phoneNum)

	if name == "" || phoneNum == "" {
		base.Fatalf("缺少参数，需要两个参数，名称和电话。")
	}
	ReplaceFile(name, phoneNum)
}

func ReplaceFile(name, phoneNum string) {
	// inputFile := "tpl.mht"
	// buf, err := ioutil.ReadFile(inputFile)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
	// 	panic(err.Error())
	// }
	// content := string(buf)

	outputFile := phoneNum + ".mht"
	content := tpl

	content = strings.ReplaceAll(content, "{{nickName}}", strings.ReplaceAll(url.QueryEscape(name[:3]), "%", "="))
	content = strings.ReplaceAll(content, "{{phoneNum}}", phoneNum[:3]+" **** "+phoneNum[7:])
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		base.Fatalf("写文件失败: %v", err)
	}
}
