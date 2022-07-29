package tools

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
)

type Cell struct {
	Row    string
	Column string
	Value  []byte
}

func CreateTable(host string, tableName string, family []string) error {
	if family == nil || len(family) < 1 {
		return fmt.Errorf("family can't be nil")
	}
	var cFamilies = make(map[string]map[string]string)
	for _, f := range family {
		cFamilies[f] = nil
	}

	ac := gohbase.NewAdminClient(host)
	crt := hrpc.NewCreateTable(context.Background(), []byte(tableName), cFamilies)
	if err := ac.CreateTable(crt); err != nil {
		return fmt.Errorf("CreateTable returned an error: %v", err)
	}

	return nil
}

func FindTables(host string, reg string) []*pb.TableName {
	ac := gohbase.NewAdminClient(host)

	tn, err := hrpc.NewListTableNames(
		context.Background(),
		hrpc.ListRegex(reg),
	)
	if err != nil {
		log.Fatal(err)
	}

	names, err := ac.ListTableNames(tn)
	if err != nil {
		log.Fatal(err)
	}

	return names
}

func ScanTable(host string, tableName string, numberOfRows int) []*Cell {
	table := []byte(tableName)
	client := gohbase.NewClient(host)
	scan, err := hrpc.NewScan(context.Background(), table, hrpc.NumberOfRows(uint32(numberOfRows)))
	if err != nil {
		log.Fatal(err)
	}

	var rsp []*hrpc.Result
	scanner := client.Scan(scan)
	for {
		res, err := scanner.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		rsp = append(rsp, res)
	}

	var cells []*Cell
	for _, row := range rsp {
		for _, cell := range row.Cells {
			cells = append(cells, &Cell{
				Row:    string(cell.Row),
				Column: string(cell.Family) + ":" + string(cell.Qualifier),
				Value:  cell.Value,
			})
		}
	}
	return cells
}

func uint64ToBytes(i uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return buf[:]
}

func delCell(client gohbase.Client) (error, *hrpc.Result) {
	values := map[string]map[string][]byte{"cf1": {"age": uint64ToBytes(18)}}
	delReq, err := hrpc.NewDelStr(context.Background(), "test1", "000", values)
	res, err := client.Delete(delReq)
	return err, res
}

func scanWithFilter(client gohbase.Client) error {
	pFilter := filter.NewPrefixFilter([]byte("00"))
	scanReq, err := hrpc.NewScanStr(context.Background(), "test1", hrpc.Filters(pFilter))
	scanner := client.Scan(scanReq)
	for rec, err := scanner.Next(); err == nil; rec, err = scanner.Next() {
		fmt.Println(rec)
	}
	return err
}

func getSpecificCell(client gohbase.Client) (error, *hrpc.Result) {
	// Perform a get for the cell with key "15", column family "cf" and qualifier "a"
	family := map[string][]string{"cf1": {"age"}}
	getRequest, err := hrpc.NewGetStr(context.Background(), "test1", "000",
		hrpc.Families(family))
	rsp, err := client.Get(getRequest)

	return err, rsp
}

func getCell(client gohbase.Client) (error, *hrpc.Result) {
	getReq, err := hrpc.NewGetStr(context.Background(), "test1", "000")
	rsp, err := client.Get(getReq)
	return err, rsp
}

func insertCell(client gohbase.Client) (error, *hrpc.Result) {
	values := map[string]map[string][]byte{"cf1": {"age": uint64ToBytes(18)}}
	putRequest, err := hrpc.NewPutStr(context.Background(), "test1", "000", values)
	rsp, err := client.Put(putRequest)
	return err, rsp
}
