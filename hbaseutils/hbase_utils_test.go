package tools

import (
	"fmt"
	"testing"
	"time"
)

var host string
var table string

func init() {
	host = "cdh1.iwellmass.com"
	table = "janusgraph"
}

func TestCreateTable(t *testing.T) {
	testTableName := "test1_" + getTimestampString()
	err := CreateTable(host, testTableName, []string{"info", "config"})
	if err != nil {
		t.Fatalf("fail")
	}
}

func getTimestampString() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func TestListTable(t *testing.T) {
	names := FindTables(host, ".*janus.*")

	for _, n := range names {
		t.Logf("name: %s \n", string(n.Qualifier))
	}
}

func TestScanTable(t *testing.T) {
	cells := ScanTable(host, table, 2)

	for _, cell := range cells {
		t.Logf("Row: [%s], Column: %s, Value: %s \n", cell.Row, cell.Column, cell.Value)
	}
}
