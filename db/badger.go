package db

import (
	"bytes"
	"sync"

	"log"

	"github.com/dgraph-io/badger/v3"
)

type BadgerDB struct {
	Badger *badger.DB
	sync.Mutex
}

// var Badger *badger.DB
var b BadgerDB

func NewDB() {
	dbPointer, err := badger.Open(badger.DefaultOptions("./badger"))

	if err != nil {
		log.Fatal(err)
	}
	b.Badger = dbPointer

	// badger.DefaultOptions("").WithInMemory(true)
}

func BadgerClose() {
	b.Badger.Close()
}

func Add(key string, value []byte) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.Badger.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(key), []byte(value)); err == badger.ErrTxnTooBig {
			_ = txn.Commit()
			txn = b.Badger.NewTransaction(true)
			_ = txn.Set([]byte(key), []byte(value))
			_ = txn.Commit()
		}
		return nil
	})
}

func Get(key string) ([]byte, bool) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	var buf bytes.Buffer
	ok := false
	err := b.Badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == nil {
			item.Value(func(val []byte) error {
				buf.Write(val)
				ok = true
				return nil
			})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes(), ok
}

func Remove(key string) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()
	b.Badger.Update(func(txn *badger.Txn) error {
		txn.Delete([]byte(key))
		return nil
	})
}
