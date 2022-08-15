package hbaseutils

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

// Name of the meta region.
const metaTableName = "hbase:meta"

type HBaseUtils struct {
	host   string
	client gohbase.Client
	ac     gohbase.AdminClient
}

// 一个Cell包含了HBase一个KV键值对信息
type Cell struct {
	Row    string
	Column string
	Value  string
}

func (c Cell) String() string {
	return fmt.Sprintf("%s:%s/%s", c.Row, c.Column, c.Value)
}

func NewHBaseUtils(host string) *HBaseUtils {
	u := HBaseUtils{host: host}
	u.client = gohbase.NewClient(host)
	u.ac = gohbase.NewAdminClient(host)
	return &u
}

func (u *HBaseUtils) CreateTable(tableName string, family []string) error {
	if family == nil || len(family) < 1 {
		return fmt.Errorf("family can't be nil")
	}
	var cFamilies = make(map[string]map[string]string)
	for _, f := range family {
		cFamilies[f] = map[string]string{}
	}

	crt := hrpc.NewCreateTable(context.Background(), []byte(tableName), cFamilies)
	if err := u.ac.CreateTable(crt); err != nil {
		return fmt.Errorf("CreateTable returned an error: %v", err)
	}
	return nil
}

func (u *HBaseUtils) DisableTable(tableName string) error {
	dt := hrpc.NewDisableTable(context.Background(), []byte(tableName))
	err := u.ac.DisableTable(dt)
	return err
}

func (u *HBaseUtils) DeleteTable(tableName string) error {
	dt := hrpc.NewDeleteTable(context.Background(), []byte(tableName))
	err := u.ac.DeleteTable(dt)
	return err
}

func (u *HBaseUtils) FindTables(reg string) []*pb.TableName {
	tn, err := hrpc.NewListTableNames(
		context.Background(),
		hrpc.ListRegex(reg),
	)
	if err != nil {
		log.Fatal(err)
	}

	names, err := u.ac.ListTableNames(tn)
	if err != nil {
		log.Fatal(err)
	}

	return names
}

func (u *HBaseUtils) FindInMeta(tableName string) []*Cell {
	metaKey := tableName + ","
	keyFilter := filter.NewPrefixFilter([]byte(metaKey))
	scan, err := hrpc.NewScanStr(context.Background(), metaTableName, hrpc.Filters(keyFilter))
	if err != nil {
		log.Fatalf("Failed to create Scan request: %s", err)
	}

	return u.GetCells(scan)
}

func (u *HBaseUtils) InsertCell(table, cf, qualifier, key, value string) (*hrpc.Result, error) {
	values := map[string]map[string][]byte{cf: {qualifier: []byte(value)}}
	putRequest, err := hrpc.NewPutStr(context.Background(), table, key, values)
	if err != nil {
		return nil, err
	}
	rsp, err := u.client.Put(putRequest)
	return rsp, err
}

func (u *HBaseUtils) GetCell(table, cf, qualifier, key string) (*Cell, error) {
	family := map[string][]string{cf: {qualifier}}
	req, err := hrpc.NewGetStr(context.Background(), table, key,
		hrpc.Families(family))
	if err != nil {
		return nil, err
	}
	rsp, err := u.client.Get(req)
	cells := rsp.Cells
	if len(cells) == 0 {
		return nil, nil
	}
	return toCell(cells[0]), err
}

func (u *HBaseUtils) DelCell(table, cf, qualifier, key string) (*hrpc.Result, error) {
	values := map[string]map[string][]byte{cf: {qualifier: nil}}
	delReq, err := hrpc.NewDelStr(context.Background(), table, key, values)
	if err != nil {
		return nil, err
	}
	res, err := u.client.Delete(delReq)
	return res, err
}

func (u *HBaseUtils) GetRow(table, key string) (*hrpc.Result, error) {
	res, err := hrpc.NewGetStr(context.Background(), table, key)
	if err != nil {
		return nil, err
	}
	rsp, err := u.client.Get(res)
	return rsp, err
}

func (u *HBaseUtils) ScanTable(tableName string, numberOfRows int) []*Cell {
	table := []byte(tableName)
	scan, err := hrpc.NewScan(context.Background(), table, hrpc.NumberOfRows(uint32(numberOfRows)))
	if err != nil {
		log.Fatal(err)
	}
	return u.GetCells(scan)
}

func (u *HBaseUtils) ScanWithPrefixFilter(tableName string, prefix string) ([]*Cell, error) {
	pFilter := filter.NewPrefixFilter([]byte(prefix))
	scanReq, err := hrpc.NewScanStr(context.Background(), tableName, hrpc.Filters(pFilter))
	if err != nil {
		return nil, err
	}
	return u.GetCells(scanReq), nil
}

func (u *HBaseUtils) GetCells(scan *hrpc.Scan) []*Cell {
	var rsp []*hrpc.Result
	scanner := u.client.Scan(scan)
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
			cells = append(cells, toCell(cell))
		}
	}
	return cells
}

func toCell(cell *hrpc.Cell) *Cell {
	return &Cell{
		Row:    string(cell.Row),
		Column: string(cell.Family) + ":" + string(cell.Qualifier),
		Value:  string(cell.Value),
	}
}

func uint64ToBytes(i uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return buf[:]
}
