package main

import "sync"

type SafeViewerResults struct {
	Done          bool
	FetchingURL   string
	Mu            sync.Mutex
	WG            sync.WaitGroup
	ViewerResults []ViewerResult
}

func (svr *SafeViewerResults) append(extend []ViewerResult, fetchingURL string) {
	svr.Mu.Lock()
	svr.FetchingURL = fetchingURL
	svr.ViewerResults = append(svr.ViewerResults, extend...)
	svr.Mu.Unlock()
}
