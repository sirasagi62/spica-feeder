package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/syndtr/goleveldb/leveldb"
)

// UIに関わる部分

func initUI(db *leveldb.DB, af ArticleFetcher, safeViewerResults *SafeViewerResults) *RSSListUIComponents {
	startText := `
┏┓  •     ┏┓     ┓
┗┓┏┓┓┏┏┓  ┣ ┏┓┏┓┏┫
┗┛┣┛┗┗┗┻  ┻ ┗ ┗ ┗┻
  ┛                                          
  ┛                                          

  ┛

  - Press '/' to search
    `
	app := tview.NewApplication()

	mainTextView := tview.NewTextView().
		SetText(startText).
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	mainTextView.SetBorder(true)
	mainTextView.SetTitle("🔷 Spica")

	// キーバインド情報の表示
	keybindings := tview.NewTextView().
		SetText("Press '/' to search, 'Ctrl+q' to quit").
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	// 検索画面のテキストボックス
	inputField := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30)

	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true)
	list.SetTitle("🔍 Search Results")
	list.ShowSecondaryText(true)

	searchFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 1, true).
		AddItem(list, 0, 1, false)

	pages := tview.NewPages().
		AddPage("main", mainTextView, true, false).
		AddPage("search", searchFlex, true, true)

	ui := &RSSListUIComponents{
		App:            app,
		List:           list,
		MainTextView:   mainTextView,
		StatusTextView: keybindings,
		Pages:          pages,
		SearchInput:    inputField,
	}

	ui.drawRSSList(safeViewerResults.ViewerResults, af)

	// 読込み終了まで再描画
	go ui.redrawRSSListUntilComplete(safeViewerResults, af)

	// テキストボックスの入力が変更されたときのハンドラ
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			ui.drawRSSList(safeViewerResults.ViewerResults, af)
			app.SetFocus(list)
		}
	})

	// "/"キーを押したときのハンドラ
	mainTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			pages.SwitchToPage("search")
			if inputField.GetText() == "" {
				app.SetFocus(inputField)
			} else {
				app.SetFocus(list)
			}
		}
		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	// 検索画面用のカスタム入力キャプチャを設定
	searchFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if app.GetFocus() == inputField {
				app.SetFocus(list)
			} else {
				app.SetFocus(inputField)
			}
			return nil
		}
		return event
	})

	return ui
}

func (ui *RSSListUIComponents) run() {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.Pages, 0, 1, true).
		AddItem(ui.StatusTextView, 1, 1, false)

	if err := ui.App.SetRoot(layout, true).Run(); err != nil {
		log.Fatal(err)
	}
	log.Print("Bye.")
}
