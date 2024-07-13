package main

import (
	"log"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/syndtr/goleveldb/leveldb"
)

type RSSFetcher struct {
	Now                  time.Time
	CacheLifeTimeSeconds float64
	Fp                   *gofeed.Parser
	DB                   *leveldb.DB
}

func getFeedResults(url string, fp *gofeed.Parser, db *leveldb.DB) []ViewerResult {
	encodedCVR, err := db.Get([]byte(url), nil)
	// そもそもdbを取得できなかった。
	if leveldb.ErrNotFound == err {
		return fetchEachFeedURL(url, fp, db)
	} else if err != nil {
		log.Fatal("Failed to read db")
		return []ViewerResult{}
	}
	cvr, err := DecodeCachedViewerResults(encodedCVR)
	if err != nil {
		log.Fatal("Failed to decode CachedViewerResults")
		return []ViewerResult{}
	}
	duration := time.Now().Sub(cvr.CachedDate).Seconds()
	// キャッシュ切れ
	if duration > 3600.0 {
		return fetchEachFeedURL(url, fp, db)
	}
	return cvr.Value
}
func fetchEachFeedURL(url string, fp *gofeed.Parser, db *leveldb.DB) []ViewerResult {
	fetchedFeeds, _ := fp.ParseURL(url)
	res := make([]ViewerResult, len(fetchedFeeds.Items))
	for i, f := range fetchedFeeds.Items {
		res[i] = ViewerResult{Title: f.Title, URL: f.Link, Date: *f.PublishedParsed}
	}
	cacheData, err := EncodeCachedViewerResults(CachedViewerResults{CachedDate: time.Now(), Value: res})
	if err != nil {
		log.Fatal("Failed to encode CVR")
		return []ViewerResult{}
	}
	err = db.Put([]byte(url), cacheData, nil)
	if err != nil {
		log.Fatal("Failed to save cache data.")
		return []ViewerResult{}
	}
	return res
}
func fetchFeed(srcs RSSFeed, db *leveldb.DB) []ViewerResult {
	fp := gofeed.NewParser()
	var feedResults []ViewerResult
	for _, src := range srcs.Src {
		if src.Main != nil {
			println(*src.Main)
			feedResults = append(feedResults, fetchEachFeedURL(*src.Main, fp, db)...)
		}
		// if src.Topic != nil {
		// 	feedResults = append(feedResults, fetchTopicFeed(*src.Topic, fp)...)
		// }
		// if src.User != nil {
		// 	feedResults = append(feedResults, fetchUserFeed(*src.User, fp)...)
		// }
	}
	return feedResults
}

func fetchTopicFeed(t Topic, fp *gofeed.Parser, db *leveldb.DB) []ViewerResult {
	var feedResults []ViewerResult
	for _, fol := range t.Following {
		feedResults = append(feedResults, fetchEachFeedURL(strings.ReplaceAll(t.URL, "$topic", fol), fp, db)...)
	}
	return feedResults
}

func fetchUserFeed(t Topic, fp *gofeed.Parser, db *leveldb.DB) []ViewerResult {
	var feedResults []ViewerResult
	for _, fol := range t.Following {
		feedResults = append(feedResults, fetchEachFeedURL(strings.ReplaceAll(t.URL, "$topic", fol), fp, db)...)
	}
	return feedResults
}
