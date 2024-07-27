package main

import (
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mmcdole/gofeed"
	"github.com/syndtr/goleveldb/leveldb"
)

type RSSFetcher struct {
	Now                  time.Time
	CacheLifeTimeSeconds float64
	Fp                   *gofeed.Parser
	DB                   *leveldb.DB
}

func (rf *RSSFetcher) getFeedResults(url string, extraTagInfo []string) []ViewerResult {
	encodedCVR, err := rf.DB.Get([]byte(url), nil)
	// そもそもdbを取得できなかった。
	if leveldb.ErrNotFound == err {
		return rf.fetchEachFeedURLOverNetwork(url, extraTagInfo)
	} else if err != nil {
		log.Fatal("Failed to read db")
		return []ViewerResult{}
	}
	cvr, err := DecodeCachedViewerResults(encodedCVR)
	if err != nil {
		log.Fatal("Failed to decode CachedViewerResults")
		return []ViewerResult{}
	}
	// duration := rf.Now.Sub(cvr.CachedDate).Seconds()
	// // キャッシュ切れ
	// if duration > 3600.0 {
	// 	return rf.fetchEachFeedURLOverNetwork(url, extraTagInfo)
	// }
	log.Printf("Use cache for %s", url)
	return cvr.Value
}
func (rf *RSSFetcher) fetchEachFeedURLOverNetwork(url string, extraTagInfo []string) []ViewerResult {
	time.Sleep(2 * time.Second)
	log.Printf("Fetch data : %s", url)
	// Fetch
	fetchedFeeds, _ := rf.Fp.ParseURL(url)
	res := make([]ViewerResult, len(fetchedFeeds.Items))
	for i, f := range fetchedFeeds.Items {
		res[i] = ViewerResult{Title: f.Title, URL: f.Link, Date: *f.PublishedParsed, Categories: append(f.Categories, extraTagInfo...), Description: f.Description}
	}

	// DBにキャッシュを保存
	cacheData, err := EncodeCachedViewerResults(CachedViewerResults{CachedDate: rf.Now, Value: res})
	if err != nil {
		log.Fatal("Failed to encode CVR")
		return []ViewerResult{}
	}
	err = rf.DB.Put([]byte(url), cacheData, nil)
	if err != nil {
		log.Fatal("Failed to save cache data.")
		return []ViewerResult{}
	}
	return res
}
func (rf *RSSFetcher) GetFeed(srcs RSSFeed, svr *SafeViewerResults) {
	for _, src := range srcs.Src {
		if src.Main != nil {
			log.Printf("Processing ....:%s", *src.Main)
			svr.append(rf.getFeedResults(*src.Main, []string{}), *src.Main)
		}
		if src.Topic != nil {
			svr.append(rf.fetchTopicFeed(*src.Topic), src.Topic.URL)
		}
		if src.User != nil {
			svr.append(rf.fetchUserFeed(*src.User), src.User.URL)
		}
	}
	svr.Done = true
	sort.Slice(svr.ViewerResults, func(i, j int) bool {
		return svr.ViewerResults[i].Date.After(svr.ViewerResults[j].Date)
	})
}

func (rf *RSSFetcher) fetchTopicFeed(t Topic) []ViewerResult {
	var feedResults []ViewerResult
	for _, fol := range t.Following {
		feedResults = append(feedResults, rf.getFeedResults(strings.ReplaceAll(t.URL, "$topic", fol), []string{fol})...)
	}
	return feedResults
}

func (rf *RSSFetcher) fetchUserFeed(t Topic) []ViewerResult {
	var feedResults []ViewerResult
	for _, fol := range t.Following {
		feedResults = append(feedResults, rf.getFeedResults(strings.ReplaceAll(t.URL, "$topic", fol), []string{})...)
	}
	return feedResults
}
