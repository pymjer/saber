package main

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
)

func TestSet(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		log.Fatal(err)
	}
	var tests = []struct {
		key, val string
	}{
		{"t1", "v1"},
		{"t2", "v2"},
		{"t3", "v3"},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("put %s,%s", tt.key, tt.val)
		t.Run(testname, func(t *testing.T) {
			Set(db, tt.key, tt.val)
			ans := Query(db, tt.key)
			if ans != tt.val {
				t.Errorf("return:%s, want %s", ans, tt.val)
			}
		})
	}
	defer db.Close()
	//Delete(db, k)
	//Query(db, k)
	//Seq(db)

}

func TestSetWithTTL(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	key := "aa"
	SetWithTTL(db, key, "aavalue", 3)
	for i := 0; i < 5; i++ {
		fmt.Printf("after %v second...\n", i)
		time.Sleep(time.Second)
		view(db, []byte(key))
	}
}

func view(db *badger.DB, key []byte) error {
	return db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			log.Fatal(err)
		}
		meta := item.UserMeta()
		valueCopy, err := item.ValueCopy(nil)
		fmt.Printf("key[%s] meta[%v] value[%s]\n", key, meta, valueCopy)
		return err
	})
}

func TestMerge(t *testing.T) {

}
