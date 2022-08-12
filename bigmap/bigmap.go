package bigmap

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

func BigMapMain() {
	path := "./data"
	fmt.Printf("请输入数据目录(默认为%s)：\n", path)
	fmt.Scanln(&path)
	fmt.Printf("当前数据目录为: %s \n", path)

	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for {
		fmt.Print("> ")
		var cmd, key, val string
		fmt.Scanln(&cmd, &key, &val)
		switch cmd {
		case "quit", "q":
			return
		case "help", "h":
			printUsage()
		case "list", "l":
			keys := AllKey(db)
			for _, v := range keys {
				fmt.Printf("%v\n", v)
			}
		case "seek", "s":
			pairs := Seek(db, key)
			printKVPairs(pairs)
		case "get":
			value, err := Query(db, key)
			if err == badger.ErrKeyNotFound {
				log.Println(err)
			} else if err != nil {
				log.Printf("未知错误:%v\n", err)
			} else {
				fmt.Printf("%v\n", value)
			}
		case "put":
			Set(db, key, val)
		case "del":
			Delete(db, key)
		case "":
		default:
			log.Fatalf("未知命令: %s", cmd)
		}
	}
}

func printKVPairs(pairs []KVPair) {
	for _, v := range pairs {
		fmt.Println(v)
	}
}

func printUsage() {
	help := `
quit,q 退出
help,h 查看帮助
list,l 查看所有key
seek,s 根据指定前缀查找所有kv对
get key 获取指定key的值
put key value 添加/跟新指定key的值
del key 删除指定key的值
	`
	fmt.Println(help)
}
