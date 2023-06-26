// Package main provides various examples of Fyne API capabilities.
package main

import (
	"FallGuys66/data"
	_ "FallGuys66/db"
	"FallGuys66/handler"
	"FallGuys66/live/douyu/DMconfig/config"
	"FallGuys66/live/douyu/DYtype"
	"FallGuys66/live/douyu/client"
	"FallGuys66/live/douyu/lib/logger"
	"FallGuys66/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/tutorials"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const preferenceCurrentTutorial = "currentTutorial"
const PDefaultLiveHostNo = "DefaultLiveHostNo"
const PLiveHosts = "LiveHosts"

var topWindow fyne.Window
var tabs *container.AppTabs

func main() {
	appSize := fyne.NewSize(995, 850)
	application := app.NewWithID("pro.lowking.fallguys66")
	application.Settings().SetTheme(&data.MyTheme{
		Regular:    data.FontSmileySansOblique,
		Bold:       data.FontSmileySansOblique,
		Italic:     data.FontSmileySansOblique,
		BoldItalic: data.FontSmileySansOblique,
		Monospace:  data.FontSmileySansOblique,
	})
	application.SetIcon(data.LogoWhite)
	// 托盘图标
	makeTray(application)
	logLifecycle(application)
	window := application.NewWindow("糖豆人直播助手：大周定制版")
	topWindow = window

	// 菜单
	window.SetMainMenu(makeMenu(application, window))
	window.SetMaster()

	logo := canvas.NewImageFromResource(data.LogoWhite)
	logoSize := 90
	logo.SetMinSize(fyne.NewSize(float32(logoSize), float32(logoSize)))
	cLogo := container.NewCenter(logo)
	cLogo.Resize(fyne.NewSize(100, 100))
	logoBlack := canvas.NewImageFromResource(data.LogoBlack)
	logoBlack.SetMinSize(fyne.NewSize(float32(logoSize), float32(logoSize)))
	cLogoBlack := container.NewCenter(logoBlack)
	cLogoBlack.Resize(fyne.NewSize(100, 100))
	padding := 10
	cLogo.Move(fyne.NewPos(float32(padding), float32(padding+10)))
	cLogoBlack.Move(fyne.NewPos(float32(padding), float32(padding+10)))
	cLogoBlack.Hide()
	go func() {
		for {
			flashEle(rand.Intn(5), cLogo, cLogoBlack)
			randTime := rand.Intn(900) + 100
			time.Sleep(time.Duration(randTime) * time.Millisecond)
		}
	}()

	// 需要渲染的元素
	var elements []fyne.CanvasObject
	elements = append(elements, cLogo)
	elements = append(elements, cLogoBlack)

	// 直播间label
	accentColor := color.RGBA{
		R: 46,
		G: 108,
		B: 246,
		A: 255,
	}
	lLiveHostLabel := canvas.NewText("直播间：", accentColor)
	lLiveHostLabel.TextSize = 35
	toolbarPaddingTop := padding + 10
	toolbarPaddingLeft := 120
	lLiveHostLabel.Move(fyne.NewPos(float32(toolbarPaddingLeft), float32(toolbarPaddingTop)))
	elements = append(elements, lLiveHostLabel)

	optionSize := fyne.NewSize(140, 37)
	// 直播间输入框
	defaultLiveHostNo := application.Preferences().StringWithFallback(PDefaultLiveHostNo, "")
	tLiveHostNo := canvas.NewText(defaultLiveHostNo, color.RGBA{
		R: 106,
		G: 135,
		B: 89,
		A: 255,
	})
	tLiveHostNo.TextSize = 35
	tLiveHostNo.Move(fyne.NewPos(cLogo.Size().Width+40+optionSize.Width+float32(padding*2), float32(toolbarPaddingTop)))
	elements = append(elements, tLiveHostNo)
	liveHosts := []string{
		"156277",
	}
	pLiveHosts := application.Preferences().StringWithFallback(PLiveHosts, "")
	if pLiveHosts != "" {
		liveHosts = strings.Split(pLiveHosts, "|")
	}
	// bindLiveHosts := binding.BindStringList(&liveHosts)
	liveHostOption := widget.NewSelectEntry(liveHosts)
	liveHostOption.TextStyle.Bold = true
	liveHostOption.OnCursorChanged = func() {
	}
	liveHostOption.SetPlaceHolder("请输入/选择直播间号")
	liveHost := ""
	liveHostOption.OnChanged = func(s string) {
		liveHost = s
	}
	liveHostOption.SetText(defaultLiveHostNo)
	cLiveHostOption := container.NewMax(liveHostOption)
	cLiveHostOption.Resize(optionSize)
	cLiveHostOption.Move(fyne.NewPos(cLogo.Size().Width+40+optionSize.Width+float32(padding*2), float32(toolbarPaddingTop)))
	elements = append(elements, cLiveHostOption)

	// 初始化copyright
	lCopyrightL := canvas.NewText("斗鱼ID：石疯悦耳", color.RGBA{
		R: 32,
		G: 32,
		B: 32,
		A: 100,
	})
	lCopyrightL.TextSize = 14
	// lCopyrightL.Move(fyne.NewPos(0, appSize.Height-25))
	lCopyrightL.Alignment = fyne.TextAlignTrailing
	lCopyrightL.Move(fyne.NewPos(appSize.Width-10, 140))
	lCopyrightR := canvas.NewText("© lowking 2023. All Rights Reserved.", color.RGBA{
		R: 32,
		G: 32,
		B: 32,
		A: 100,
	})
	lCopyrightR.TextSize = 14
	lCopyrightR.Alignment = fyne.TextAlignTrailing
	// lCopyrightR.Move(fyne.NewPos(appSize.Width-10, appSize.Height-25))
	lCopyrightR.Move(fyne.NewPos(appSize.Width-110, 140))
	elements = append(elements, lCopyrightL)
	elements = append(elements, lCopyrightR)

	// 给列表行添加背景
	// rowHeight := 39
	// for i := 0; i < int(appSize.Height/float32(rowHeight)); i+=2 {
	// 	tableRowLine := canvas.NewRectangle(color.RGBA{
	// 		R: 234,
	// 		G: 234,
	// 		B: 234,
	// 		A: 255,
	// 	})
	// 	tableRowLine.Resize(fyne.NewSize(appSize.Width-10, float32(rowHeight)))
	// 	tableRowLine.Move(fyne.NewPos(0, cLogo.Size().Height+float32(padding*3)+114+float32((rowHeight)*i)-float32(i)*1.4))
	// 	elements = append(elements, tableRowLine)
	// }

	// 初始化tab列表
	tabs = container.NewAppTabs(
		container.NewTabItem("未玩列表", makeEmptyList(accentColor)),
		container.NewTabItem("已玩列表", makeEmptyList(accentColor)),
		container.NewTabItem("收藏列表", makeEmptyList(accentColor)),
	)
	cTabList := container.NewBorder(nil, nil, nil, nil, tabs)
	cTabList.Move(fyne.NewPos(0, cLogo.Size().Height+float32(padding*3)))
	cTabList.Resize(fyne.NewSize(appSize.Width, appSize.Height-cLogo.Size().Height))
	elements = append(elements, cTabList)

	// 直线分割表头
	// tableHeaderLine := canvas.NewRectangle(color.RGBA{
	// 	R: 204,
	// 	G: 204,
	// 	B: 204,
	// 	A: 255,
	// })
	// tableHeaderLine.Resize(fyne.NewSize(appSize.Width, 5))
	// tableHeaderLine.Move(fyne.NewPos(0, cLogo.Size().Height+float32(padding*3)+75))
	// elements = append(elements, tableHeaderLine)

	// 连接直播间按钮
	var webSocketClient client.DyBarrageWebSocketClient
	var btnCon *widget.Button
	btnCon = widget.NewButtonWithIcon("连接", theme.NavigateNextIcon(), func() {
		if liveHost == "" {
			btnConSetDefault(btnCon, liveHostOption)
			dialog.ShowInformation("提示", "请填写斗鱼直播间号", window)
			return
		} else if _, err := strconv.Atoi(liveHost); err != nil && liveHost != "dev" {
			btnConSetDefault(btnCon, liveHostOption)
			dialog.ShowInformation("提示", "请输入纯数字直播间号", window)
			liveHostOption.SetText("")
			window.Canvas().Focus(liveHostOption)
			return
		}
		// 根据当前状态执行对应操作
		switch btnCon.Importance {
		case widget.MediumImportance, widget.DangerImportance:
			// 未连接 连接错误，再次点击连接弹幕，状态设置成连接中
			btnCon.Importance = widget.HighImportance
			btnCon.Icon = theme.ViewRefreshIcon()
			btnCon.Text = "..."
			btnCon.Refresh()
			connectSuc := false
			time.Sleep(1 * time.Second)
			if liveHost != "dev" {
				// todo 连接弹幕
				logger.Infof("connecting douyu: %s", liveHost)
				spiderConfig := &config.DMconfig{
					Rid:            liveHost,
					LoginMsg:       "type@=loginreq/room_id@=%s/dfl@=sn@A=105@Sss@A=1/username@=%s/uid@=%s/ver@=20190610/aver@=218101901/ct@=0/",
					LoginJoinGroup: "type@=joingroup/rid@=%s/gid@=-9999/",
					Url:            "wss://danmuproxy.douyu.com:8506/",
				}
				out := make(chan client.Item)
				webSocketClient = client.DyBarrageWebSocketClient{
					ItemIn: out,
					Config: spiderConfig,
					MsgBreakers: DYtype.CodeBreakershandler{
						IsLive: false,
					},
				}
				go func() {
					webSocketClient.Init()
					webSocketClient.Start()
					logger.Infof("disconnect %s", liveHost)
				}()
				go func() {
					// 获取弹幕，处理消息
					for {
						msg := <-out
						switch msg.Type {
						case "chatmsg":
							handler.FilterMap(msg)
						default:
							logger.ShowJson("[%s]not handle msg: %s", msg.Type, msg)
						}
					}
				}()
				connectSuc = true
			} else {
				// 开发模式模拟获取id
				connectSuc = rand.Intn(10)%2 == 1
				// 模拟数据，写入数据库
				go func() {
					count := 0
					for {
						if btnCon.Importance != widget.HighImportance {
							log.Printf("exit filter loop")
							break
						}
						// mock map id
						mapId := fmt.Sprintf("%d-%d-%d ceshi sdfkjij", rand.Intn(9000)+1000, rand.Intn(9000)+1000, rand.Intn(9000)+1000)
						log.Printf("loop %s", mapId)
						count++
						handler.FilterMap(client.Item{
							Rid:     "",
							Cid:     "",
							Uid:     "",
							Type:    "",
							Txt:     mapId,
							Nn:      "",
							Level:   "",
							Payload: nil,
						})
					}
					logger.Infof("mock %d map", count)
				}()
			}
			if connectSuc {
				// 连接成功
				btnCon.Importance = widget.HighImportance
				btnCon.Text = "成功"
				btnCon.Icon = theme.ConfirmIcon()
				liveHostOption.Hide()
				tLiveHostNo.Text = liveHost
				tLiveHostNo.Refresh()
				// 更新下拉框
				if liveHost != "dev" && !utils.In(liveHosts, liveHost) {
					liveHosts = append(liveHosts, liveHost)
					liveHostOption.SetOptions(liveHosts)
					liveHostOption.Refresh()
					application.Preferences().SetString(PLiveHosts, strings.Join(liveHosts, "|"))
				}
				// application.Preferences().SetString(PLiveHosts, "")
				// 保存默认值
				application.Preferences().SetString(PDefaultLiveHostNo, liveHost)

				// 定时查询获取最新地图
				go func() {
					// playedMapList := tabs.Items[1]
					// starMapList := tabs.Items[2]
					for {
						time.Sleep(2 * time.Second)
						if btnCon.Importance != widget.HighImportance {
							log.Printf("exit query loop")
							return
						}
						go handler.RefreshMapList(tabs, 0, `and state="0"`, `order by created asc, mapId`)
						go handler.RefreshMapList(tabs, 1, `and state="1"`, `order by created desc, mapId`)
						go handler.RefreshMapList(tabs, 2, `and star="1"`, `order by created desc, mapId`)
					}
				}()
			} else {
				// 连接失败
				btnCon.Importance = widget.DangerImportance
				btnCon.Text = "错误"
				btnCon.Icon = theme.CancelIcon()
			}
		case widget.HighImportance:
			// 连接中 已连接，再次点击断开连接，状态设置成未连接
			btnConSetDefault(btnCon, liveHostOption)
			if liveHost != "dev" {
				webSocketClient.Stop()
			}
		}
		btnCon.Refresh()
	})
	btnCon.Importance = widget.MediumImportance
	cBtnCon := container.NewMax(btnCon)
	cBtnCon.Resize(fyne.NewSize(65, 36))
	cBtnCon.Move(fyne.NewPos(cLogo.Size().Width+120+float32(padding), float32(toolbarPaddingTop)))
	elements = append(elements, cBtnCon)

	gridLayout := container.NewWithoutLayout(elements...)

	window.SetContent(gridLayout)
	window.Resize(appSize)
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.ShowAndRun()
}

func btnConSetDefault(btnCon *widget.Button, liveHostOption *widget.SelectEntry) {
	btnCon.Importance = widget.MediumImportance
	btnCon.Text = "连接"
	btnCon.Icon = theme.NavigateNextIcon()
	liveHostOption.Show()
}

func makeEmptyList(accentColor color.RGBA) *fyne.Container {
	text := canvas.NewText("无数据 ...", accentColor)
	text.TextSize = 20
	cEmpty := container.NewCenter(text)
	return cEmpty
}

func flashEle(times int, elements ...*fyne.Container) {
	for i := 0; i < times; i++ {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		for _, element := range elements {
			if element.Visible() {
				element.Hide()
			} else {
				element.Show()
			}
		}
	}
}

func allWidgets(application fyne.App, window fyne.Window) {
	topWindow = window
	window.CenterOnScreen()
	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as application")
	intro.Wrapping = fyne.TextWrapWord

	setTutorial := func(t tutorials.Tutorial) {
		if fyne.CurrentDevice().IsMobile() {
			child := application.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = window
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(window)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		window.SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		window.SetContent(split)
	}
	window.Resize(fyne.NewSize(800, 600))
	window.Show()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
		// 第一次加载列表
		if a.Preferences().StringWithFallback(PDefaultLiveHostNo, "") != "" {
			go handler.RefreshMapList(tabs, 0, `and state="0"`, `order by created asc, mapId`)
		}
		go handler.RefreshMapList(tabs, 1, `and state="1"`, `order by created desc, mapId`)
		go handler.RefreshMapList(tabs, 2, `and star="1"`, `order by created desc, mapId`)
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)
	checkedItem := fyne.NewMenuItem("Checked", nil)
	checkedItem.Checked = true
	disabledItem := fyne.NewMenuItem("Disabled", nil)
	disabledItem.Disabled = true
	otherItem := fyne.NewMenuItem("Other", nil)
	mailItem := fyne.NewMenuItem("Mail", func() { logger.Debug("Menu New->Other->Mail") })
	mailItem.Icon = theme.MailComposeIcon()
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { logger.Debug("Menu New->Other->Project") }),
		mailItem,
	)
	fileItem := fyne.NewMenuItem("File", func() { logger.Debug("Menu New->File") })
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("Directory", func() { logger.Debug("Menu New->Directory") })
	dirItem.Icon = theme.FolderIcon()
	newItem.ChildMenu = fyne.NewMenu("",
		fileItem,
		dirItem,
		otherItem,
	)

	openSettings := func() {
		w := a.NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	cutShortcut := &fyne.ShortcutCut{Clipboard: w.Clipboard()}
	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(cutShortcut, w)
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: w.Clipboard()}
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(copyShortcut, w)
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: w.Clipboard()}
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(pasteShortcut, w)
	})
	pasteItem.Shortcut = pasteShortcut
	performFind := func() { logger.Debug("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = a.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", newItem, checkedItem, disabledItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	main := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
	checkedItem.Action = func() {
		checkedItem.Checked = !checkedItem.Checked
		main.Refresh()
	}
	return main
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("显示所有组件（开发用）", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			allWidgets(a, a.NewWindow("Widgets"))
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func unsupportedTutorial(t tutorials.Tutorial) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

func makeNav(setTutorial func(tutorial tutorials.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tutorials.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tutorials.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := tutorials.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := tutorials.Tutorials[uid]; ok {
				if unsupportedTutorial(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutCut:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutPaste:
		sh.Clipboard = w.Clipboard()
	}
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
