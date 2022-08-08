package main

import (
	"bytes"
	"context"

	"github.com/dgraph-io/badger"
)

func Stream() {
	db, _ := badger.Open(badger.DefaultOptions("./data"))
	defer db.Close()

	stream := db.NewStream()

	// -- Optional setting
	stream.NumGo = 16
	stream.Prefix = []byte("k-9")
	stream.LogPrefix = "Badger.Streaming"

	stream.ChooseKey = func(item *badger.Item) bool {
		return bytes.HasSuffix(item.Key(), []byte("2"))
	}

	stream.KeyToList = nil // convert badger data into custom key-values

	// -- End of optional settings

	// stream.Send = func(b *z.Buffer) error {
	// 	fmt.Printf(string(b.Bytes()))
	// 	return nil
	// }

	if err := stream.Orchestrate(context.Background()); err != nil {

	}
}
