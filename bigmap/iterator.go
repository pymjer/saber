package main

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

func AllKey(db *badger.DB) []string {
	res := []string{}
	db.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchValues = false // 是否取值，如果不取可以快很多，然后在循环的时候取
		it := txn.NewIterator(options)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			res = append(res, string(item.Key()))
		}
		return nil
	})
	return res
}

type KVPair struct {
	key   string
	value string
}

func (p KVPair) String() string {
	return fmt.Sprintf("%s:%s", p.key, p.value)
}

func Seek(db *badger.DB, prefix string) []KVPair {
	res := []KVPair{}
	db.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchSize = 10
		it := txn.NewIterator(options)
		defer it.Close()
		prefixByte := []byte(prefix)
		for it.Seek(prefixByte); it.ValidForPrefix(prefixByte); it.Next() {
			item := it.Item()
			valCopy, _ := item.ValueCopy(nil)
			res = append(res, KVPair{string(item.Key()), string(valCopy)})
		}
		return nil
	})
	return res
}

func All(db *badger.DB) error {
	return db.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchSize = 10

		it := txn.NewIterator(options)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			valCopy, _ := item.ValueCopy(nil)
			fmt.Printf("key[%s], value[%s]\n", key, valCopy)
		}

		return nil
	})
}
