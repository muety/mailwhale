package config

import (
	"encoding/json"
	"github.com/emvi/logbuch"
	"github.com/timshannon/bolthold"
)

var store *bolthold.Store

func LoadStore(path string) *bolthold.Store {
	if s, err := bolthold.Open(path, 0664, &bolthold.Options{
		Encoder: json.Marshal,
		Decoder: json.Unmarshal,
	}); err != nil {
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
