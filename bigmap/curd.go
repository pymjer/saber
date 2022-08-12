package bigmap

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger"
)

// Create
func Set(db *badger.DB, k string, v string) error {
	txn := db.NewTransaction(true)
	defer txn.Discard()
	err := txn.Set([]byte(k), []byte(v))
	txn.Commit()
	return err
}

func SetWithTTL(db *badger.DB, key string, val string, second int) {
	db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(val)).
			WithTTL(time.Second * time.Duration(second)).
			WithMeta(byte(3))
		err := txn.SetEntry(e)
		return err
	})
}

// Query
func Query(db *badger.DB, k string) (string, error) {
	txn := db.NewTransaction(false)
	item, err := txn.Get([]byte(k))
	if err != nil {
		return "", err
	}
	valCopy, err := item.ValueCopy(nil)
	if err != nil {
		log.Fatal(err)
	}
	return string(valCopy), nil
}

// Delete
func Delete(db *badger.DB, k string) error {
	txn := db.NewTransaction(true)
	defer txn.Discard()
	err := txn.Delete([]byte(k))
	txn.Commit()
	return err
}

func Seq(db *badger.DB) {
	seq, _ := db.GetSequence([]byte("k"), 1000)
	defer seq.Release()
	for i := 0; i < 10; i++ {
		next, err := seq.Next()
		fmt.Printf("next: %v, err: %v\n", next, err)
	}
}

// Merge
func add(originalValue, newValue []byte) []byte {
	return append(originalValue, newValue...)
}

func Merge(db *badger.DB) {
	key := []byte("merge")

	m := db.GetMergeOperator(key, add, 200*time.Millisecond)
	defer m.Stop()

	m.Add([]byte("A"))
	m.Add([]byte("B"))
	m.Add([]byte("C"))

	res, _ := m.Get()
	fmt.Printf("Merge result:%s\n", res)
}

func uint64ToBytes(i uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)
	return buf[:]
}

func bytesToUint64(b []byte) uint64 {
	//fmt.Println(b)
	return binary.BigEndian.Uint64(b)
}
