package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	start_text := `
  _____                 __     ___               
 |__  /___ _ __  _ __   \ \   / (_) _____      __
   / // _ \ '_ \| '_ \   \ \ / /| |/ _ \ \ /\ / /
  / /|  __/ | | | | | |   \ V / | |  __/\ V  V / 
 /____\___|_| |_|_| |_|    \_/  |_|\___| \_/\_/  

 - Press '/' to search
	`
	app := tview.NewApplication()

	// メイン画面のテキストビュー
	mainTextView := tview.NewTextView().
		SetText(start_text).
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	mainTextView.SetBorder(true)
	mainTextView.SetTitle("🚀ZennView")

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

	// 検索画面のリスト
	list := tview.NewList().
		ShowSecondaryText(false)

	list.SetBorder(true)
	list.SetTitle("Search Results")
	list.ShowSecondaryText(true)

	// 検索画面のレイアウト
	searchFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 1, true).
		AddItem(list, 0, 1, false)

	// Pagesを作成
	pages := tview.NewPages().
		AddPage("main", mainTextView, true, true).
		AddPage("search", searchFlex, true, false)

	// テキストボックスの入力が変更されたときのハンドラ
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			list.Clear()
			text := inputField.GetText()
			if text == "" {
				return
			}
			results := convertResult(executeSearch(text))
			for _, item := range results {
				// item := item // クロージャで変数のコピーを作成
				list.AddItem(item.Title, "- Update at "+item.Date.Local().UTC().Format("2006/1/2"), 0, nil)
				list.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
					mainTextView.Clear()
					mainTextView.ScrollToBeginning()
					mainTextView.SetText(drawArticle(results[i].URL))
					mainTextView.SetTitle(results[i].Title)
					app.SetFocus(mainTextView)
					pages.SwitchToPage("main")
				})
			}
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
		if event.Key() == tcell.KeyCtrlQ {
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

	// メインレイアウトを作成
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(keybindings, 1, 1, false)
	// アプリケーションの起動
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
