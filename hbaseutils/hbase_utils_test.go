package tools

import (
	"fmt"
	"testing"
	"time"
)

var host string
var testTableName string
var testCf string
var testQualifier string
var u *HBaseUtils

func init() {
	host = "cdh1.iwellmass.com"
	// testTableName = "test1_" + getTimestampString()
	testTableName = "test1201_0601"
	testCf = "info"
	testQualifier = "name"
	u = NewHBaseUtils(host)
}

func TestCreateTable(t *testing.T) {
	t.Logf("testTableName: %s \n", testTableName)

	err := u.CreateTable(testTableName, []string{"info", "config"})
	if err != nil {
		t.Fatalf("fail")
	}
	cells := u.FindInMeta(testTableName)
	printCells(t, cells)
}

func TestFindInMeta(t *testing.T) {
	cells := u.FindInMeta(testTableName)
	printCells(t, cells)
}

func TestInsertCell(t *testing.T) {
	table, cf, qualifier := testTableName, testCf, testQualifier
	var key, value string
	for i := 0; i < 10; i++ {
		key = fmt.Sprintf("00%d", i)
		value = fmt.Sprintf("testvalue:%d", i)
		_, err := u.InsertCell(table, cf, qualifier, key, value)
		if err != nil {
			t.Fail()
		}

		cell, err := u.GetCell(table, cf, qualifier, key)
		if err != nil {
			t.Fail()
		}
		t.Logf("cell: %s \n", cell)
	}
}

func TestDelCell(t *testing.T) {
	table, cf, qualifier := testTableName, testCf, testQualifier
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("00%d", i)
		_, err := u.DelCell(table, cf, qualifier, key)
		if err != nil {
			t.Fail()
		}

		cell, err := u.GetCell(table, cf, qualifier, key)
		if err != nil {
			t.Fail()
		}
		t.Logf("cell: %s \n", cell)
	}
}

func getTimestampString() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func TestListTable(t *testing.T) {
	names := u.FindTables("test1_.*")

	for _, n := range names {
		t.Logf("name: %s \n", string(n.Qualifier))
	}
}

func TestDeleteTable(t *testing.T) {
	names := u.FindTables("test1_.*")

	for _, n := range names {
		tableName := string(n.Qualifier)
		t.Logf("delete table: %s \n", tableName)
		u.DisableTable(tableName)
		u.DeleteTable(tableName)
	}
}

func TestScanTable(t *testing.T) {
	cells := u.ScanTable(testTableName, 2)
	printCells(t, cells)
}

func printCells(t *testing.T, cells []*Cell) {
	for _, cell := range cells {
		t.Logf("Row: [%s], Column: %s, Value: %s \n", cell.Row, cell.Column, cell.Value)
	}
}
