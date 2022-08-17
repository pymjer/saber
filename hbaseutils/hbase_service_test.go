package hbaseutils

import (
	"fmt"
	"testing"
	"time"
)

var testTableName string
var testCf string
var testQualifier string
var u *HBaseUtils

func init() {
	var host = "zj01.iwellmass.com"
	// testTableName = "test1_" + getTimestampString()
	testTableName = "test0817"
	testCf = "info"
	testQualifier = "name"
	u = NewHBaseUtils(host)
}

func TestListTable(t *testing.T) {
	names := u.ListTables()

	for _, n := range names {
		t.Logf("name: %s \n", string(n.Qualifier))
	}
}

func TestFindTables(t *testing.T) {
	names := u.FindTables("zlx.*")

	for _, n := range names {
		t.Logf("name: %s \n", string(n.Qualifier))
	}
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
		key = fmt.Sprintf("02%d", i)
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
		key := fmt.Sprintf("01%d", i)
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

func TestScanWithPrefixFilter(t *testing.T) {
	cells, err := u.ScanWithPrefixFilter(testTableName, "01", 5)
	if err != nil {
		t.Fail()
	}
	printCells(t, cells)
}

func printCells(t *testing.T, cells []*Cell) {
	for _, cell := range cells {
		t.Logf("Row: [%s], Column: %s, Value: %s \n", cell.Row, cell.Column, cell.Value)
	}
}
