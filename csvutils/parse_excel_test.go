package csvutils

import (
	"testing"
)

func TestParseExcel(t *testing.T) {
	excelFile := "../../../clouda-data/example.xlsx"
	ParseExcel(excelFile, "../../../clouda-data", "excel", false, 1000)
}

func TestParseExcelToCsv(t *testing.T) {
	excelFile := "../../../clouda-data/example.xlsx"
	ParseExcel(excelFile, "../../../clouda-data", "csv", false, 1000)
}

func TestExcelToCSV(t *testing.T) {
	// 输入 Excel 文件路径和目标 Sheet 名称
	excelFile := "../../../clouda-data/example.xlsx"
	// excelFile := "../../../clouda-data/sample_2w.xlsx"
	outputPath := "output"
	batchSize := 10000 // 每次写入 1000 行后刷新

	// 调用函数进行转换
	ParseExcel(excelFile, outputPath, "csv", false, batchSize)
}

func TestExcelToCSVByStream(t *testing.T) {
	// 输入 Excel 文件路径和目标 Sheet 名称
	// excelFile := "../../../clouda-data/example.xlsx"
	excelFile := "../../../clouda-data/sample_2w.xlsx"
	outputPath := "output"
	batchSize := 10000 // 每次写入 1000 行后刷新

	// 调用函数进行转换
	ParseExcel(excelFile, outputPath, "csv", true, batchSize)
}
