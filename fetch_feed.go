package main

import (
	"strings"

	"github.com/mmcdole/gofeed"
)

func fetchEachFeedURL(url string, fp *gofeed.Parser) []ViewerResult {
	fetchedFeeds, _ := fp.ParseURL(url)
	res := make([]ViewerResult, len(fetchedFeeds.Items))
	for i, f := range fetchedFeeds.Items {
		res[i] = ViewerResult{Title: f.Title, URL: f.Link, Date: *f.PublishedParsed}
	}
	return res
}
func fetchFeed(srcs RSSFeed) []ViewerResult {
	fp := gofeed.NewParser()
	var feedResults []ViewerResult
	for _, src := range srcs.Src {
		if src.Main != nil {
			feedResults = append(feedResults, fetchEachFeedURL(*src.Main, fp)...)
		}
		if src.Topic != nil {
			feedResults = append(feedResults, fetchTopicFeed(*src.Topic, fp)...)
		}
		if src.User != nil {
			feedResults = append(feedResults, fetchUserFeed(*src.User, fp)...)
		}
	}
	return feedResults
}

func fetchTopicFeed(t Topic, fp *gofeed.Parser) []ViewerResult {
	var feedResults []ViewerResult
	for _, fol := range t.Following {
		feedResults = append(feedResults, fetchEachFeedURL(strings.ReplaceAll(t.URL, "$topic", fol), fp)...)
	}
	return feedResults
}

func fetchUserFeed(t Topic, fp *gofeed.Parser) []ViewerResult {
	var feedResults []ViewerResult
	for _, fol := range t.Following {
		feedResults = append(feedResults, fetchEachFeedURL(strings.ReplaceAll(t.URL, "$topic", fol), fp)...)
	}
	return feedResults
}
