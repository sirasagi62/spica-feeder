package main

import (
	"log"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
	"github.com/go-shiori/go-readability"
	"github.com/rivo/tview"
	"github.com/syndtr/goleveldb/leveldb"
)

// 空白と改行のみの行を削除する関数
func removeEmptyLines(input string) string {
	var output []string
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		// 空白と改行のみの行を削除
		if strings.TrimSpace(line) != "" {
			output = append(output, line)
		}
	}

	return strings.Join(output, "\n")
}

type ArticleFetcher struct {
	DB        *leveldb.DB
	converter *md.Converter
}

func InitArticleFetcher() ArticleFetcher {
	// Open the article.db file for saving articles. It will be created if it doesn't exist.
	articledb, err := leveldb.OpenFile("article.db", nil)
	if err != nil {
		log.Fatal(err)
	}
	return ArticleFetcher{converter: md.NewConverter("", true, nil), DB: articledb}
}
func (af ArticleFetcher) Close() {
	af.DB.Close()
}
func (af ArticleFetcher) fetchArticle(vr ViewerResult) string {
	article, err := readability.FromURL(vr.URL, 30*time.Second)
	if err != nil {
		return "Not Article"
	}
	m, err := af.converter.ConvertString(article.Content)
	if err != nil {
		return "Failed to convert to markdown"
	}
	content := removeEmptyLines(m)
	eca, err := EncodeCachedArticle(CachedArticle{Content: content, URL: vr.URL, Description: vr.Description, Categories: vr.Categories})
	if err != nil {
		return content
	}
	err = af.DB.Put([]byte(vr.URL), eca, nil)
	if err != nil {
		return content
	}
	return content
}

func (af ArticleFetcher) GetArticle(vr ViewerResult) string {
	encodedArticle, err := af.DB.Get([]byte(vr.URL), nil)
	if err == leveldb.ErrNotFound {
		return af.fetchArticle(vr)
	} else if err != nil {
		return ""
	}
	ca, err := DecodeCachedArticle(encodedArticle)
	if err != nil {
		log.Fatal("Failed to decode article")
	}
	return ca.Content
}

func (af ArticleFetcher) DrawArticle(vr ViewerResult) string {
	out, err := glamour.Render(af.GetArticle(vr), "dark")
	if err != nil {
		return "Failed to drawing the content."
	}
	return tview.TranslateANSI(out)
}
