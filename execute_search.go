package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func executeSearch(q string) ZennSearchResult {
	// URLを指定
	url := "https://zenn.dev/api/search?q=" + url.QueryEscape(q) + "&order=daily&source=articles&page=1"

	// GETリクエストを送信
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error:", err)
		return ZennSearchResult{}
	}
	defer response.Body.Close()

	// 結果をバイトスライスとして読み取る
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error:", err)
		return ZennSearchResult{}
	}

	r, err := UnmarshalZennSearchResult(body)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error:", err)
		return ZennSearchResult{}
	}
	return r
}

func convertResult(z ZennSearchResult) []ViewerResult {
	res := make([]ViewerResult, len(z.Articles))
	for i, item := range z.Articles {
		res[i] = ViewerResult{
			Title: item.Emoji + "　" + item.Title,
			URL:   "https://zenn.dev" + item.Path,
			Date:  item.BodyUpdatedAt,
		}
	}
	return res
}
