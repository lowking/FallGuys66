// Package main provides various examples of Fyne API capabilities.
package main

import (
	"FallGuys66/data"
	"FallGuys66/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/tutorials"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"math/rand"
	"net/url"
	"time"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window
var dev = true

func main() {
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
	lLiveHostLabel := canvas.NewText("直播间：", color.RGBA{
		R: 100,
		G: 100,
		B: 100,
		A: 255,
	})
	lLiveHostLabel.TextSize = 35
	toolbarPaddingTop := padding + 10
	toolbarPaddingLeft := 120
	lLiveHostLabel.Move(fyne.NewPos(float32(toolbarPaddingLeft), float32(toolbarPaddingTop)))
	elements = append(elements, lLiveHostLabel)

	optionSize := fyne.NewSize(140, 37)
	// 直播间输入框
	liveHostNo := canvas.NewText("156277", color.RGBA{
		R: 106,
		G: 135,
		B: 89,
		A: 255,
	})
	liveHostNo.TextSize = 35
	liveHostNo.Move(fyne.NewPos(cLogo.Size().Width+40+optionSize.Width+float32(padding*2), float32(toolbarPaddingTop)))
	elements = append(elements, liveHostNo)
	liveHosts := []string{
		"156277",
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
	cLiveHostOption := container.NewMax(liveHostOption)
	cLiveHostOption.Resize(optionSize)
	cLiveHostOption.Move(fyne.NewPos(cLogo.Size().Width+40+optionSize.Width+float32(padding*2), float32(toolbarPaddingTop)))
	elements = append(elements, cLiveHostOption)

	// 连接直播间按钮
	var btnCon *widget.Button
	btnCon = widget.NewButtonWithIcon("连接", theme.NavigateNextIcon(), func() {
		// 根据当前状态执行对应操作
		switch btnCon.Importance {
		case widget.MediumImportance, widget.DangerImportance:
			// 未连接 连接错误，再次点击连接弹幕，状态设置成连接中
			btnCon.Importance = widget.HighImportance
			btnCon.Icon = theme.ViewRefreshIcon()
			btnCon.Text = "..."
			btnCon.Refresh()
			// todo 连接弹幕
			time.Sleep(1 * time.Second)
			if rand.Intn(10)%2 == 1 {
				// 连接成功
				btnCon.Importance = widget.HighImportance
				btnCon.Text = "成功"
				btnCon.Icon = theme.ConfirmIcon()
				liveHostOption.Hide()
			} else {
				// 连接失败
				btnCon.Importance = widget.DangerImportance
				btnCon.Text = "错误"
				btnCon.Icon = theme.CancelIcon()
			}
		case widget.HighImportance, widget.WarningImportance:
			// 连接中 已连接，再次点击断开连接，状态设置成未连接
			btnCon.Importance = widget.MediumImportance
			btnCon.Text = "连接"
			btnCon.Icon = theme.NavigateNextIcon()
			liveHostOption.Show()
		}
		if liveHost != "" && !utils.In(liveHosts, liveHost) {
			liveHosts = append(liveHosts, liveHost)
			liveHostOption.SetOptions(liveHosts)
			liveHostOption.Refresh()
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
	window.Resize(fyne.NewSize(640, 460))
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.ShowAndRun()
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
	mailItem := fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") })
	mailItem.Icon = theme.MailComposeIcon()
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		mailItem,
	)
	fileItem := fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") })
	fileItem.Icon = theme.FileIcon()
	dirItem := fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") })
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
	performFind := func() { fmt.Println("Menu Find") }
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
