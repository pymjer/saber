package main

import (
	"testing"

	"github.com/dgraph-io/badger"
)

func TestAllKey(t *testing.T) {
	db, _ := badger.Open(badger.DefaultOptions("./data"))
	defer db.Close()

	//All(db)
	//Seek(db, "k-888") // 指定前缀查询
	AllKey(db)
}
