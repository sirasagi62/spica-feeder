package main

import (
	"strings"
	"time"

	"github.com/samber/lo"
)

type ViewerResult struct {
	Title string
	URL   string
	Date  time.Time
}

func fileterViewerResultByName(searchStr string, vr *[]ViewerResult) []ViewerResult {
	return lo.Filter(*vr, func(item ViewerResult, index int) bool {
		return strings.Contains(item.Title, searchStr)
	})
}
