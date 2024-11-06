package csvutils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// 将 Excel Sheet 转换为 CSV 文件，逐行流式写入并定期刷新
func ExcelToCSVByStream(f *excelize.File, sheetName string, csvFile string, batchSize int) error {
	// 打开 CSV 文件
	csvFileHandle, err := os.Create(csvFile)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer csvFileHandle.Close()

	// 创建 CSV 写入器
	csvWriter := csv.NewWriter(csvFileHandle)

	// 获取指定的 Sheet 中的所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet: %v", err)
	}
	// 获取总行数
	totalRows := len(rows)
	totalBatchs := (totalRows + batchSize - 1) / batchSize

	// 初始化计数器
	rowCount := 0
	batchCount := 0

	// 创建进度通道
	progressChan := make(chan float64)

	// 启动一个 goroutine 来处理进度打印
	go printProgress(progressChan)

	// 逐行读取 Excel 数据并写入 CSV
	for _, row := range rows {
		// 写入当前行
		err := csvWriter.Write(row)
		if err != nil {
			return fmt.Errorf("failed to write to CSV: %v", err)
		}

		// 增加行计数
		rowCount++

		// 每隔 batchSize 行刷新一次 CSV 文件
		if rowCount >= batchSize {
			// 刷新 CSV 内容到磁盘
			csvWriter.Flush()

			if err := csvWriter.Error(); err != nil {
				return fmt.Errorf("failed to flush CSV data: %v", err)
			}
			// 更新批次计数
			batchCount++

			// 计算进度百分比
			percentage := float64(batchCount) / float64(totalBatchs) * 100

			// 将进度百分比发送到进度通道
			progressChan <- percentage

			// 重置行计数
			rowCount = 0
		}
	}

	// 确保文件内容被写入到磁盘
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("failed to flush CSV data: %v", err)
	}
	// 最终更新进度为 100%
	progressChan <- 100.0
	return nil
}

// 将单个 Sheet 转换为新的 Excel 文件
func SaveSheetAsNewExcelByStream(f *excelize.File, sheetName, outputPath string) error {
	// 创建一个新的 Excel 文件
	fmt.Printf("save file[%s] by stream", sheetName)
	newFile := excelize.NewFile()

	// 将指定的 Sheet 复制到新的 Excel 文件中
	// 如果源文件中有多个 Sheet，应该复制指定的 Sheet 数据
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet: %v", err)
	}

	// 删除默认的空 Sheet 页面
	defaultSheetName := newFile.GetSheetName(0) // 获取默认的第一个 Sheet 页名称
	newFile.DeleteSheet(defaultSheetName)

	sw, err := newFile.NewStreamWriter(sheetName)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i, row := range rows {
		// 将每一行数据写入新的 Excel 文件
		b := make([]interface{}, len(row)) // 创建一个长度与a相同的[]interface{}切片
		for i := range row {
			b[i] = row[i] // 将a切片中的每个元素转换为interface{}并赋值给b切片
		}
		err := sw.SetRow("A"+strconv.Itoa(i), b)
		if err != nil {
			fmt.Println(err)
		}
	}
	if err := sw.Flush(); err != nil {
		fmt.Println(err)
	}

	// 保存新的 Excel 文件
	if err := newFile.SaveAs(outputPath); err != nil {
		return fmt.Errorf("failed to save new Excel file: %v", err)
	}
	return nil
}

// 打印进度条，监听 progressChan
func printProgress(progressChan chan float64) {
	for percentage := range progressChan {
		// 更新进度，使用退格符 \r 更新同一行
		fmt.Printf("Progress: %.2f%% \n", percentage)
	}
}
