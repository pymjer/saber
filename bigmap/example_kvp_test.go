// 这个示例演示如何使用Bigmap
package bigmap_test

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
	"prolion.top/saber/bigmap"
)

type KVP struct {
	key, val string
}

// 这个示例打开一个默认的kv数据库，添加一些值，然后做查询操作
func Example_kVP() {
	path := "./data"
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}

	var tests = []KVP{
		{"t1", "v1"},
		{"t2", "v2"},
		{"t3", "v3"},
	}

	for _, tt := range tests {
		bigmap.Set(db, tt.key, tt.val)
		ans, err := bigmap.Query(db, tt.key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s:%s ", tt.key, ans)
	}
	defer db.Close()
	// Output:
	// t1:v1 t2:v2 t3:v3
}
