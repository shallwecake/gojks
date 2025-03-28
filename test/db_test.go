package test

import (
	"github.com/dgraph-io/badger/v3"
	"log"
	"testing"
)

func TestDb(t *testing.T) {
	// 创建或打开数据库
	opts := badger.DefaultOptions("./db")
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 插入数据
	//err = db.Update(func(txn *badger.Txn) error {
	//	return txn.Set([]byte("answer"), []byte("42"))
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 读取数据
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		log.Printf("The answer is: %s\n", string(val))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
