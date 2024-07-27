package main

import (
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	initLogging()
	db := initDatabase()
	defer db.Close()
	af := InitArticleFetcher()
	defer af.Close()

	safeViewerResults := SafeViewerResults{ViewerResults: []ViewerResult{}, Done: false}
	initFeeder(db, &safeViewerResults)

	ui := initUI(db, af, &safeViewerResults)
	ui.run()
}

// 初期化部分
func initLogging() {
	file, err := os.OpenFile("spica.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// エラーハンドリング
		log.Fatal(err)
	}

	// 関数が終了する際にファイルを閉じる
	defer file.Close()

	// ログの出力先をファイルに設定
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Print("Init App.")
}

func initDatabase() *leveldb.DB {
	db, err := leveldb.OpenFile("data.db", nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
