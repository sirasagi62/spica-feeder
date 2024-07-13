package main

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"

	"github.com/samber/lo"
)

type ViewerResult struct {
	Title string
	URL   string
	Date  time.Time
}

func filterViewerResultByName(searchStr string, vr *[]ViewerResult) []ViewerResult {
	if searchStr == "" {
		return *vr
	}
	return lo.Filter(*vr, func(item ViewerResult, index int) bool {
		return strings.Contains(item.Title, searchStr)
	})
}

type CachedViewerResults struct {
	CachedDate time.Time
	Value      []ViewerResult
}

// CachedViewerResultsをgobにエンコードする関数
func EncodeCachedViewerResults(cvr CachedViewerResults) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(cvr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// gobからCachedViewerResultsをデコードする関数
func DecodeCachedViewerResults(data []byte) (CachedViewerResults, error) {
	var cvr CachedViewerResults
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&cvr)
	if err != nil {
		return CachedViewerResults{}, err
	}
	return cvr, nil
}
