package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mmcdole/gofeed"
	"github.com/syndtr/goleveldb/leveldb"
)

func initFeeder(db *leveldb.DB, svr *SafeViewerResults) {
	f, err := os.ReadFile("./default2.toml")
	if err != nil {
		log.Fatal("cannot load config file.")
	}
	rf, _ := UnmarshalRSSFeed(f)
	fetcher := FeedFetcher{Now: time.Now(), CacheLifeTimeSeconds: 3600, Fp: gofeed.NewParser(), DB: db}
	go fetcher.GetFeed(rf, svr)
}
