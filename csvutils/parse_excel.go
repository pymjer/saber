package csvutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
	"prolion.top/saber/internal/base"
)

var CmdParseExcel = &base.Command{
	Run:       runParseExcel,
	UsageLine: "saber csvutils [flags] [filePath]",
	Short:     "parse excel to csv",
	Long: `
csvutils parse excel to csv.

csvutils acceps one arguments.

The -f flag is filepath
The -o flag is output filepath
The -t flag is convert type 
The -s flag is use stream
The -b flag is batch size
}

Examples:
	saber csvutils -t csv example.xlsx
		parse xlsx file to csv
	saber csvutils -t excel example.xlsx
		parse multi sheet xlsx file to excel with single sheet
	`,
}

var (
	convertType = "excel"
	filePath    string
	outputPath  string
	stream      bool
	batchSize   int
)

func init() {
	CmdParseExcel.Flag.StringVar(&convertType, "t", "excel", "输出文件路径")
	CmdParseExcel.Flag.StringVar(&filePath, "f", "", "解析文件路径")
	CmdParseExcel.Flag.StringVar(&outputPath, "o", ".", "输出文件路径")
	CmdParseExcel.Flag.BoolVar(&stream, "s", false, "是否流式写入数据")
	CmdParseExcel.Flag.IntVar(&batchSize, "b", 1000, "batch size")
}

func runParseExcel(ctx context.Context, cmd *base.Command, args []string) {
	if len(args) >= 2 {
		filePath = args[len(args)-2]
		outputPath = args[len(args)-1]
	}
	log.Printf("begin parse file %s to %s , outputPath: %s, use stream:[%t]\n", filePath, convertType, outputPath, stream)
	ParseExcel(filePath, outputPath, convertType, stream, batchSize)
}

// 将单个 Sheet 转换为新的 Excel 文件
func saveSheetAsNewExcel(f *excelize.File, sheetName, outputPath string) error {
	// 创建一个新的 Excel 文件
	newFile := excelize.NewFile()

	// 将指定的 Sheet 复制到新的 Excel 文件中
	// 如果源文件中有多个 Sheet，应该复制指定的 Sheet 数据
	rows, err := f.GetRows(sheetName)

	log.Printf("total rows[%d]\n", len(rows))

	if err != nil {
		return fmt.Errorf("failed to get rows from sheet: %v", err)
	}
	// 创建新的 Sheet，并写入数据
	newFile.NewSheet(sheetName)

	// 删除默认的空 Sheet 页面
	defaultSheetName := newFile.GetSheetName(0) // 获取默认的第一个 Sheet 页名称
	newFile.DeleteSheet(defaultSheetName)

	for i, row := range rows {
		// 将每一行数据写入新的 Excel 文件
		for j, cellValue := range row {
			// 获取列名，支持从 A 到 Z，AA 到 AZ 等
			cell := getColumnName(j) + fmt.Sprintf("%d", i+1) // A1, B1, C1 ... 表示单元格
			err := newFile.SetCellValue(sheetName, cell, cellValue)
			if err != nil {
				return fmt.Errorf("failed to set cell value: %v", err)
			}
		}
	}

	// 保存新的 Excel 文件
	err = newFile.SaveAs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to save new Excel file: %v", err)
	}

	return nil
}

// 获取列名，如 0 -> A, 1 -> B, 26 -> AA, 27 -> AB 等
func getColumnName(index int) string {
	column := ""
	for index >= 0 {
		column = string(rune(index%26+65)) + column
		index = index/26 - 1
	}
	return column
}

// 将 Sheet 转换为 CSV 格式并保存
func saveSheetAsCSV(f *excelize.File, sheetName, outputPath string) error {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet: %v", err)
	}

	// 打开 CSV 文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	// 写入 CSV 数据
	for _, row := range rows {
		line := strings.Join(row, ",") + "\n"
		_, err := file.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write to CSV file: %v", err)
		}
	}

	return nil
}

func ParseExcel(excelFile string, outputPath string, convertType string, stream bool, batchSize int) {
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		fmt.Println("Error opening Excel file:", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 获取文件名
	fileName := filepath.Base(excelFile)

	// 去掉文件扩展名
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 获取所有 Sheet 名称
	sheetNames := f.GetSheetList()

	fmt.Printf("Excel file has %d sheets\n", len(sheetNames))

	// 确保目录存在
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Printf("failed to create directory: %v", err)
		return
	}
	// 将每个 Sheet 转换为 CSV or Excel 文件
	for i, sheetName := range sheetNames {
		fmt.Printf("parse sheet %d to %s\n", i, convertType)

		var err error
		var resultFile string

		if convertType == "excel" {
			resultFile = fmt.Sprintf("%s/%s-%s.xlsx", outputPath, fileName, sheetName)
			if stream {
				err = SaveSheetAsNewExcelByStream(f, sheetName, resultFile)
			} else {
				err = saveSheetAsNewExcel(f, sheetName, resultFile)
			}
		} else {
			resultFile = fmt.Sprintf("%s/%s-%s.csv", outputPath, fileName, sheetName)
			if stream {
				err = ExcelToCSVByStream(f, sheetName, resultFile, batchSize)
			} else {
				err = saveSheetAsCSV(f, sheetName, resultFile)
			}
		}

		if err != nil {
			fmt.Printf("Error saving sheet '%s' as %s: %v\n", sheetName, resultFile, err)
			continue
		}
		fmt.Printf("Sheet '%s' has been saved as '%s'\n", sheetName, resultFile)
	}
}
