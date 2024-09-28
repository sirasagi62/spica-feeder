package main

import "time"

func bulk_fetch_article(svr *SafeViewerResults, af ArticleFetcher) {
	// RSSの内容を読み込み終わるまで待機
	for !svr.Done {
		time.Sleep(1 * time.Second)
	}

	for _, article := range svr.ViewerResults {
		println(af.GetArticle(article))
		println("----Get", article.URL)
		time.Sleep(1 * time.Second)
	}
}
