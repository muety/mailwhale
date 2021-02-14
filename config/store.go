package config

import (
	"github.com/emvi/logbuch"
	"github.com/timshannon/bolthold"
)

var store *bolthold.Store

func InitStore(path string) *bolthold.Store {
	if s, err := bolthold.Open(path, 0664, nil); err != nil {
		logbuch.Fatal("failed to open store: %v", err)
	} else {
		store = s
	}
	return store
}

func GetStore() *bolthold.Store {
	return store
}

func CloseStore() {
	if store != nil {
		store.Close()
	}
}
