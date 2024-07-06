package main

import (
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/charmbracelet/glamour"
	"github.com/go-shiori/go-readability"
	"github.com/rivo/tview"
)

// 空白と改行のみの行を削除する関数
func RemoveEmptyLines(input string) string {
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

func fetchArticle(url string) string {
	converter := md.NewConverter("", true, nil)
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return "Not Article"
	}
	m, err := converter.ConvertString(article.Content)
	if err != nil {
		return "Failed to convert to markdown"
	}
	return RemoveEmptyLines(m)
}

func drawArticle(url string) string {
	out, err := glamour.Render(fetchArticle(url), "dark")
	if err != nil {
		return "Failed to drawing the content."
	}
	return tview.TranslateANSI(out)
}
