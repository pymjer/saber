package csvutils

import "testing"

func TestParseExcel(t *testing.T) {
	excelFile := "../../../clouda-data/example.xlsx"
	ParseExcel(excelFile, "../../../clouda-data")
}
