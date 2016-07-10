package dashboard

import (
	"strings"

	"github.com/lxn/walk"
	ui "github.com/lxn/walk/declarative"
	"github.com/simpleelegant/devtools/plugins/network"
	"github.com/skratchdot/open-golang/open"
)

// Run 启用GUI控制面板，只能用于 Windows
func Run(url string) {
	var ip *walk.TextEdit

	ui.MainWindow{
		Title:   "DevTools 服务端",
		MinSize: ui.Size{Width: 600, Height: 400},
		Layout:  ui.VBox{},
		Children: []ui.Widget{
			ui.VSpacer{
				Size: 20,
			},
			ui.HSplitter{
				Children: []ui.Widget{
					ui.Label{
						Text: "Server IP:",
					},
				},
			},
			ui.HSplitter{
				Children: []ui.Widget{
					ui.TextEdit{AssignTo: &ip},
					ui.VSplitter{
						Children: []ui.Widget{
							ui.PushButton{
								Text: "刷新",
								OnClicked: func() {
									ips, err := network.GetLocalIPs()
									if err != nil {
										ip.SetText("error on refresh!")
									}

									ip.SetText(strings.Join(ips, "\r\n"))
								},
							},
							ui.PushButton{
								Text: "复制",
								OnClicked: func() {
									if err := walk.Clipboard().SetText(ip.Text()); err != nil {
										ip.SetText("error on copy!")
									}
								},
							},
						},
					},
				},
			},
			ui.VSpacer{
				Size: 20,
			},
			ui.PushButton{
				Text:    "在浏览器中打开 DevTools",
				MinSize: ui.Size{Height: 40},
				OnClicked: func() {
					open.Start(url)
				},
			},
			ui.VSpacer{
				Size: 100,
			},
			ui.Label{
				Text: "Author: Wang Yujian <simpleelegant@163.com>",
			},
		},
	}.Run()
}
