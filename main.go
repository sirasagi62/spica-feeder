package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// メイン画面のテキストビュー
	mainTextView := tview.NewTextView().
		SetText("Press '/' to search").
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

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

	// 初期データをリストに追加
	items := []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"}
	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}

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
			text := inputField.GetText()
			list.Clear()
			for _, item := range items {
				if text == "" || contains(item, text) {
					item := item // クロージャで変数のコピーを作成
					list.AddItem(item, "", 0, func() {
						mainTextView.SetText("Selected: " + item)
						app.SetFocus(mainTextView)
						pages.SwitchToPage("main")
					})
				}
			}
		}
	})
	// "/"キーを押したときのハンドラ
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			inputField.SetText("")
			pages.SwitchToPage("search")
			app.SetFocus(inputField)
		}
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

// 部分一致をチェックする関数
func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}
