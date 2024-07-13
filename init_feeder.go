package main

import (
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

func initFeeder(db *leveldb.DB) []ViewerResult {
	f, err := os.ReadFile("./default2.toml")
	if err != nil {
		log.Fatal("cannot load config file.")
	}
	rf, _ := UnmarshalRSSFeed(f)
	return fetchFeed(rf, db)
}
