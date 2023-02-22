package simulate

import "testing"

func TestSimulateData(t *testing.T) {
	SimulateData("id@int,name,amount@float", "test.json", "json", 10)
	SimulateData("id@int,name,age@int,amount@float", "test.csv", "csv", 50)
}
