// simulate包用于查找包含指定字符的文件
package simulate

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"prolion.top/saber/internal/base"
)

var CmdSimulateData = &base.Command{
	Run:       runSimulateData,
	UsageLine: "saber sdata [flags] [fileName]",
	Short:     "simulate data",
	Long: `
Simulate Data.

Column is a comma-separated field, the field defaults to string type, if the type is int or float, use the field$type

The -c flag is columns
The -n flag is the row numbers, default 100
The -t flag is file type, support csv,json, default is csv
The -f flag is the result file name. default is {current_time_stamp}.json

Examples:
	saber sdata -c id,name,age 
	saber sdata -c id,name,age$int
	saber sdata -c id,name -t json -f user.json
	`,
}

var (
	columns  string
	nums     int
	ftype    string
	fileName string
)

func init() {
	CmdSimulateData.Flag.StringVar(&columns, "c", "", "字段名，逗号分隔")
	CmdSimulateData.Flag.IntVar(&nums, "n", 100, "结果行数")
	CmdSimulateData.Flag.StringVar(&ftype, "t", "csv", "文件类型，json/csv")
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Column struct {
	Name string
	Type string
}

func runSimulateData(ctx context.Context, cmd *base.Command, args []string) {
	if columns == "" {
		base.Fatalf("columns is not specified.")
	}
	if len(args) == 1 {
		fileName = args[0]
	} else {
		fileName = strconv.Itoa(time.Now().Nanosecond()) + "." + ftype
	}
	SimulateData(columns, fileName, ftype, nums)
}

func SimulateData(columns string, fileName string, ftype string, nums int) {
	fmt.Printf("开始生成文件[%s], 文件类型[%s]\n", fileName, ftype)
	// 写入表头
	rand.Seed(time.Now().UTC().UnixNano())
	header := strings.Split(columns, ",")
	fields := make([]Column, len(header))
	for i, v := range header {
		sep := strings.Split(v, "@")
		t := "string"
		if len(sep) == 2 {
			t = sep[1]
		}
		fields[i] = Column{sep[0], t}
	}
	if ftype == "csv" {
		CsvWriter(fields, fileName, nums)
	} else if ftype == "json" {
		JsonWriter(fields, fileName, nums)
	} else {
		panic("只支持csv或者json类型")
	}
}

func CsvWriter(fields []Column, fileName string, nums int) {
	// 创建一个csv文件
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 创建CSV writer
	writer := csv.NewWriter(f)
	defer writer.Flush()
	header := make([]string, len(fields))
	for i, v := range fields {
		header[i] = v.Name
	}
	writer.Write(header)

	// 循环写入
	record := make([]string, len(header))
	for i := 0; i < nums; i++ {
		for j, v := range fields {
			if v.Type == "string" {
				record[j] = RandStringBytes(rand.Intn(10) + 1)
			} else if v.Type == "int" {
				record[j] = strconv.Itoa(rand.Intn(10000))
			} else if v.Type == "float" {
				record[j] = fmt.Sprintf("%.2f", rand.Float64()*10000)
			} else {
				panic("not support type" + v.Type)
			}
		}
		writer.Write(record)
	}
}

func JsonWriter(fields []Column, fileName string, nums int) {
	data := make([]map[string]interface{}, nums)
	for i := 0; i < nums; i++ {
		record := map[string]interface{}{}
		for _, v := range fields {
			name := v.Name
			if v.Type == "string" {
				record[name] = RandStringBytes(rand.Intn(10) + 1)
			} else if v.Type == "int" {
				record[name] = rand.Intn(10000)
			} else if v.Type == "float" {
				record[name] = toFixed(rand.Float64()*10000, 2)
			} else {
				panic("not support type" + v.Type)
			}
		}
		data[i] = record
	}
	file, _ := json.MarshalIndent(data, "", " ")
	_ = os.WriteFile(fileName, file, 0644)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(math.Round(num*output)) / output
}
