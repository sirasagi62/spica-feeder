package main

import (
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/syndtr/goleveldb/leveldb"
)

type RSSListUIComponents struct {
	App            *tview.Application
	Pages          *tview.Pages
	List           *tview.List
	MainTextView   *tview.TextView
	StatusTextView *tview.TextView
	SearchInput    *tview.InputField
}

func drawRSSList(ui *RSSListUIComponents, RSS []ViewerResult) {
	ui.List.Clear()
	text := ui.SearchInput.GetText()
	results := filterViewerResultByName(text, &RSS)
	for _, item := range results {
		ui.List.AddItem(item.Title+" - "+item.Date.Local().UTC().Format("2006/1/2 15:04"), item.URL, 0, nil)
		ui.List.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
			ui.MainTextView.Clear()
			ui.MainTextView.ScrollToBeginning()
			ui.MainTextView.SetText(drawArticle(results[i].URL))
			ui.MainTextView.SetTitle(results[i].Title)
			ui.App.SetFocus(ui.MainTextView)
			ui.Pages.SwitchToPage("main")
		})
	}
}

func redrawRSSListUntilComplete(ui *RSSListUIComponents, svr *SafeViewerResults) {
	for {
		if svr.Done {
			log.Println("Completed to fetch RSS Feeds.")
			ui.App.QueueUpdateDraw(func() {
				ui.StatusTextView.SetText("Completed to fetch :)")
				drawRSSList(ui, svr.ViewerResults)
			})

			return
		}
		ui.App.QueueUpdateDraw(func() {
			ui.StatusTextView.SetText("Fetching :" + svr.FetchingURL)
			drawRSSList(ui, svr.ViewerResults)
		})
		time.Sleep(1 * time.Second)
	}
}

type SafeViewerResults struct {
	Done          bool
	FetchingURL   string
	Mu            sync.Mutex
	WG            sync.WaitGroup
	ViewerResults []ViewerResult
}

func main() {
	// Logging
	// ログを書き込むファイルを開く（なければ作成）
	file, err := os.OpenFile("spica.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// エラーハンドリング
		log.Fatal(err)
	}

	// 関数が終了する際にファイルを閉じる
	defer file.Close()

	// ログの出力先をファイルに設定
	log.SetOutput(file)

	// ログのフォーマットを設定（時間を含める）
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	// Logging END

	start_text := `
┏┓  •     ┏┓     ┓
┗┓┏┓┓┏┏┓  ┣ ┏┓┏┓┏┫
┗┛┣┛┗┗┗┻  ┻ ┗ ┗ ┗┻
  ┛                                          

 - Press '/' to search
	`
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := leveldb.OpenFile("data.db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Print("Init App.")

	app := tview.NewApplication()
	safeViewerResults := SafeViewerResults{ViewerResults: []ViewerResult{}, Done: false}
	initFeeder(db, &safeViewerResults)
	// メイン画面のテキストビュー
	mainTextView := tview.NewTextView().
		SetText(start_text).
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	mainTextView.SetBorder(true)
	mainTextView.SetTitle("🚀 ZennView")

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
	list.SetTitle("🔍 Search Results")
	list.ShowSecondaryText(true)

	// 検索画面のレイアウト
	searchFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 1, 1, true).
		AddItem(list, 0, 1, false)

	// Pagesを作成
	pages := tview.NewPages().
		AddPage("main", mainTextView, true, false).
		AddPage("search", searchFlex, true, true)

	ui := &RSSListUIComponents{}
	ui.App = app
	ui.List = list
	ui.MainTextView = mainTextView
	ui.StatusTextView = keybindings
	ui.Pages = pages
	ui.SearchInput = inputField
	drawRSSList(ui, safeViewerResults.ViewerResults)

	// 読込み終了まで再描画
	go redrawRSSListUntilComplete(ui, &safeViewerResults)

	// テキストボックスの入力が変更されたときのハンドラ
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			drawRSSList(ui, safeViewerResults.ViewerResults)
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
	log.Printf("Bye.")
}
