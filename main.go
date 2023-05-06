package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var Version = "dev"

func main() {
	var outTE *walk.TextEdit

	var urlInput *walk.LineEdit
	var pathInput *walk.LineEdit
	var proxyInput *walk.LineEdit

	// walk.TabPage

	var startBtn, stopBtn *walk.PushButton
	var down *Download
	// stopBtn

	var ch = make(chan bool)

	var mw *MainWindow

	mw = &MainWindow{
		Title:  "漫画下载器" + Version,
		Size:   Size{Width: 400, Height: 600},
		Layout: VBox{},
		MenuItems: []MenuItem{
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text:        "About",
						OnTriggered: func() { aboutAction_Triggered() },
					},
					Action{
						Text:        "页面空白",
						OnTriggered: func() { tips() },
					},
				},
			},
		},
		Children: []Widget{
			Label{Text: "页面空白的话，切换一下tab", TextColor: walk.RGB(102, 178, 255)},
			TabWidget{
				Pages: []TabPage{

					{
						Title: "下载",
						// Layout: Grid{Columns: 2},
						Visible: true,
						Layout:  VBox{},
						Children: []Widget{
							Composite{
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{Text: "网址", TextColor: walk.RGB(255, 0, 127)},
									LineEdit{
										Text:     "",
										AssignTo: &urlInput,
									},
									Label{Text: "保存路径"},
									LineEdit{
										Text:     ``,
										AssignTo: &pathInput,
									},
									Label{Text: "代理"},
									LineEdit{
										Text:     "",
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
										stopBtn.SetEnabled(false)
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
					},
					{
						Title: "About",
						// Layout: Grid{Columns: 2},
						Layout:  VBox{},
						Visible: true,
						Children: []Widget{
							Label{Text: "漫画下载器" + Version},
							Label{Text: "下载：有的网站会限流，会导致失败；失败的图片不会保存；重新下载即可（多次下载的时候遇到磁盘已有文件就跳过）"},
							TextEdit{
								Text: "项目地址: https://github.com/hwhaocool/hanman\r\n" +
									"支持解析 www.tuwenhanman.com\r\n" +
									"支持解析 bingmh.com\r\n",

								ReadOnly: true},
						},
					},
				},
			},
		},
	}

	mw.Run()
}

func aboutAction_Triggered() {
	walk.MsgBox(nil,
		"提示",
		"漫画下载器, enjoy it~",
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}
func tips() {
	walk.MsgBox(nil,
		"提示",
		"页面空白的话，切换一下tab",
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}
