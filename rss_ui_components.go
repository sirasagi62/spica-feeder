package main

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rivo/tview"
)

type RSSListUIComponents struct {
	App            *tview.Application
	List           *tview.List
	MainTextView   *tview.TextView
	StatusTextView *tview.TextView
	Pages          *tview.Pages
	SearchInput    *tview.InputField
}

func (ui *RSSListUIComponents) drawRSSList(RSS []ViewerResult, af ArticleFetcher) {
	ui.List.Clear()
	text := ui.SearchInput.GetText()
	results := filterViewerResultByName(text, &RSS)
	for _, item := range results {
		ui.List.AddItem(item.Title+" - "+item.Date.Local().UTC().Format("2006/1/2 15:04"), "#"+strings.Join(item.Categories, "#")+":"+item.Description, 0, nil)
		ui.List.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
			ui.MainTextView.Clear()
			ui.MainTextView.ScrollToBeginning()
			ui.MainTextView.SetText(af.DrawArticle(results[i]))
			ui.MainTextView.SetTitle(results[i].Title)
			ui.App.SetFocus(ui.MainTextView)
			ui.Pages.SwitchToPage("main")
		})
	}
}

func (ui *RSSListUIComponents) redrawRSSListUntilComplete(svr *SafeViewerResults, af ArticleFetcher) {
	for {
		if svr.Done {
			log.Println("Completed to fetch RSS Feeds.")
			ui.App.QueueUpdateDraw(func() {
				ui.StatusTextView.SetText("Completed to fetch :)")
				ui.drawRSSList(svr.ViewerResults, af)
			})

			return
		}
		ui.App.QueueUpdateDraw(func() {
			ui.StatusTextView.SetText("Fetching :" + svr.FetchingURL)
			ui.drawRSSList(svr.ViewerResults, af)
		})
		time.Sleep(1 * time.Second)
	}
}
