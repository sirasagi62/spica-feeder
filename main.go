package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	file := initLogging()
	// 関数が終了する際にファイルを閉じる
	defer file.Close()
	db := initDatabase()
	defer db.Close()
	af := InitArticleFetcher()
	defer af.Close()

	safeViewerResults := SafeViewerResults{ViewerResults: []ViewerResult{}, Done: false}
	initFeeder(db, &safeViewerResults)

	//bulk_fetch_article(&safeViewerResults, af)

	ui := initUI(af, &safeViewerResults)
	ui.run()
}

// 初期化部分
func initLogging() *os.File {
	file, err := os.OpenFile("spica.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// エラーハンドリング
		log.Fatal(err)
	}

	// ログの出力先をファイルに設定
	log.SetOutput(file)
	log.Print("Init App.")

	return file
}

func initDatabase() *leveldb.DB {
	db, err := leveldb.OpenFile("data.db", nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
