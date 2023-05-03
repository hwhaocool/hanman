package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var outTE *walk.TextEdit

	var urlInput *walk.LineEdit
	var pathInput *walk.LineEdit
	var proxyInput *walk.LineEdit

	var startBtn, stopBtn *walk.PushButton
	var down *Download
	// stopBtn

	var ch = make(chan bool)

	var mw *MainWindow

	mw = &MainWindow{
		Title:   "漫画下载器",
		MinSize: Size{Width: 200, Height: 200},
		MaxSize: Size{Width: 400, Height: 400},
		Size:    Size{Width: 400, Height: 400},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "网址", TextColor: walk.RGB(255, 0, 127)},
					LineEdit{
						AssignTo: &urlInput,
					},
					Label{Text: "保存路径"},
					LineEdit{
						AssignTo: &pathInput,
					},
					Label{Text: "代理"},
					LineEdit{
						AssignTo: &proxyInput, ToolTipText: "比如 socks5://127.0.0.1:1080",
					},
				},
			},
			PushButton{
				Text:       "下载",
				AssignTo:   &startBtn,
				MaxSize:    Size{Width: 400, Height: 100},
				MinSize:    Size{Width: 300, Height: 30},
				Background: SolidColorBrush{Color: walk.RGB(0x5F, 0x69, 0x8E)},
				OnClicked: func() {
					url := urlInput.Text()
					path := pathInput.Text()
					proxy := proxyInput.Text()

					if url == "" {
						walk.MsgBox(nil, "错误", "网址为空", walk.MsgBoxIconError)
						return
					}
					if path == "" {
						walk.MsgBox(nil, "错误", "网址为空", walk.MsgBoxIconError)
						return
					}

					outTE.SetText("开始下载\r\n")

					startBtn.SetEnabled(false)
					stopBtn.SetEnabled(true)

					down = NewDownload(url, path, proxy, outTE)

					// 在协程里下载，避免界面卡死
					go down.Start(ch)

					// 在协程里等待任务执行完成
					go func() {
						<-ch
						startBtn.SetEnabled(true)
						walk.MsgBox(nil, "OK", "任务已结束", walk.MsgBoxOK)
					}()

				},
			},
			PushButton{
				Text:       "停止",
				AssignTo:   &stopBtn,
				MaxSize:    Size{Width: 400, Height: 100},
				MinSize:    Size{Width: 300, Height: 30},
				Background: SolidColorBrush{Color: walk.RGB(0x5F, 0x69, 0x8E)},
				Enabled:    false,
				OnClicked: func() {

					down.Stop()

				},
			},

			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "log", TextColor: walk.RGB(255, 0, 127)},
					TextEdit{AssignTo: &outTE, ReadOnly: true, VScroll: true},
				},
			},
		},
	}

	mw.Run()
}
