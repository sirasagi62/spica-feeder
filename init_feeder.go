package main

import (
	"log"
	"os"
)

func initFeeder() []ViewerResult {
	f, err := os.ReadFile("./default2.toml")
	if err != nil {
		log.Fatal("cannot load config file.")
	}
	rf, _ := UnmarshalRSSFeed(f)
	return fetchFeed(rf)
}
