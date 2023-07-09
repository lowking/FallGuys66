package settings

import (
	"FallGuys66/common/cbm"
	"FallGuys66/config"
	"FallGuys66/data"
	"FallGuys66/db"
	"FallGuys66/db/model"
	"FallGuys66/hotkeys"
	"FallGuys66/live/douyu/lib/logger"
	"FallGuys66/widgets/searchentry"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"golang.design/x/hotkey"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	lineHeight                  = float32(50)
	cbIndentWidth               = float32(30)
	PAutoGetFgPid               = "PAutoGetFgPid"
	PAutoFillMapId              = "PAutoFillMapId"
	PAutoConnect                = "PAutoConnect"
	PSelectShowPos              = "PSelectShowPos"
	PEnterMapIdPos              = "PEnterMapIdPos"
	PCodeEntryPos               = "PCodeEntryPos"
	PConfirmBtnPos              = "PConfirmBtnPos"
	PHotKeyPlayNext             = "PHotKeyPlayNext"
	PHotKeyPlayNextOnSelectShow = "PHotKeyPlayNextOnSelectShow"
	PScFocusSearch              = "PScFocusSearch"
	FgName                      = "FallGuys_client"
	fgLabelWidth                = float32(100)
	commonShortcutLabelWidth    = float32(150)
)

type Settings struct {
	FgPid         string
	BtnGetFgPid   *widget.Button
	BtnCon        *widget.Button
	AutoGetFgPid  bool
	AutoFillMapId bool
	AutoConnect   bool
	Window        *fyne.Window

	PosSelectShow *string
	PosEnterMapId *string
	PosCodeEntry  *string
	PosConfirmBtn *string

	SearchShortcut *desktop.CustomShortcut
	OtherEntry     map[string]*searchentry.SearchEntry

	propertyLock       sync.RWMutex
	commonSettingItems []fyne.CanvasObject
	fgSettingItems     []fyne.CanvasObject
	hotKey             map[string]*hotkey.Hotkey
	altName            string
	isNotify           *bool
}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) Init(window *fyne.Window) *container.AppTabs {
	s.Window = window
	b := false
	s.isNotify = &b
	switch runtime.GOOS {
	case "darwin":
		s.altName = "option"
	default:
		s.altName = "alt"
	}
	settingTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("通用", theme.SettingsIcon(), s.GenCommonSettings()),
		container.NewTabItemWithIcon("　糖豆人　", data.FgLogo, s.GenFgSettings()),
	)
	settingTabs.SetTabLocation(container.TabLocationLeading)
	s.hotKey = make(map[string]*hotkey.Hotkey)
	s.OtherEntry = make(map[string]*searchentry.SearchEntry)

	return settingTabs
}

func (s *Settings) GenCommonSettings() *fyne.Container {
	app := fyne.CurrentApp()
	y := config.Padding
	startupLabel := widget.NewLabel("启动时")
	startupLabel.Alignment = fyne.TextAlignLeading
	startupLabel.Resize(fyne.NewSize(commonShortcutLabelWidth, lineHeight))
	startupLabel.Move(fyne.NewPos(config.Padding, y))
	s.commonSettingItems = append(s.commonSettingItems, startupLabel)

	y += lineHeight * 0.5
	s.AutoGetFgPid = app.Preferences().BoolWithFallback(PAutoGetFgPid, false)
	cbAutoGetFgPid := widget.NewCheckWithData("自动获取糖豆人进程ID", binding.BindBool(&s.AutoGetFgPid))
	cbAutoGetFgPid.OnChanged = func(b bool) {
		s.AutoGetFgPid = b
		app.Preferences().SetBool(PAutoGetFgPid, s.AutoGetFgPid)
	}
	cCbAutoGetFgPid := container.NewHBox(cbAutoGetFgPid)
	cCbAutoGetFgPid.Move(fyne.NewPos(config.Padding+cbIndentWidth, y))
	s.commonSettingItems = append(s.commonSettingItems, cCbAutoGetFgPid)

	y += lineHeight * 0.5
	s.AutoConnect = app.Preferences().BoolWithFallback(PAutoConnect, false)
	cbAutoConnect := widget.NewCheckWithData("自动连接直播间弹幕", binding.BindBool(&s.AutoConnect))
	cbAutoConnect.OnChanged = func(b bool) {
		s.AutoConnect = b
		app.Preferences().SetBool(PAutoConnect, s.AutoConnect)
	}
	cCbAutoConnect := container.NewHBox(cbAutoConnect)
	cCbAutoConnect.Move(fyne.NewPos(config.Padding+cbIndentWidth, y))
	s.commonSettingItems = append(s.commonSettingItems, cCbAutoConnect)

	y += lineHeight
	shortcutLabel := widget.NewLabel("快捷键")
	shortcutLabel.Alignment = fyne.TextAlignLeading
	shortcutLabel.Resize(fyne.NewSize(commonShortcutLabelWidth, lineHeight))
	shortcutLabel.Move(fyne.NewPos(config.Padding, y))
	s.commonSettingItems = append(s.commonSettingItems, shortcutLabel)

	// 定位搜索文本框
	y += lineHeight * 0.5
	scFocusSearchLabel := widget.NewLabel("定位搜索文本框：")
	scFocusSearchLabel.Alignment = fyne.TextAlignTrailing
	scFocusSearchLabel.Resize(fyne.NewSize(commonShortcutLabelWidth, lineHeight))
	scFocusSearchLabel.Move(fyne.NewPos(config.Padding, y))
	s.commonSettingItems = append(s.commonSettingItems, scFocusSearchLabel)
	scFocusSearchEntry := searchentry.NewSearchEntry("例：ctrl+f，回车保存")
	scFocusSearchEntry.Wrapping = fyne.TextTruncate
	scFocusSearchEntry.Resize(fyne.NewSize(170, 35))
	scFocusSearchEntry.Move(fyne.NewPos(config.Padding+scFocusSearchLabel.Size().Width, y))
	scFocusSearchEntry.OnSubmitted = func(str string) {
		if str == "" {
			return
		}
		keys := strings.Split(str, "+")
		if len(keys) != 2 {
			if *s.isNotify {
				dialog.ShowInformation("提示", "快捷键必须是2个按键", *s.Window)
			}
			return
		}
		scFocusSearch := &desktop.CustomShortcut{
			KeyName:  s.getKeyForFyne(keys[1]),
			Modifier: s.getModifier(keys[0]),
		}
		s.SearchShortcut = scFocusSearch
		(*s.Window).Canvas().RemoveShortcut(scFocusSearch)
		(*s.Window).Canvas().AddShortcut(scFocusSearch, func(shortcut fyne.Shortcut) {
			s.OtherEntry["keyWordEntry"].Focus()
			s.OtherEntry["keyWordEntry"].TypedShortcut(&fyne.ShortcutSelectAll{})
		})
		if *s.isNotify {
			app.Preferences().SetString(PScFocusSearch, strings.TrimSpace(str))
			dialog.ShowInformation("提示", fmt.Sprintf("快捷键[%s]设置完成", str), *s.Window)
		}
	}
	scFocusSearchStr := app.Preferences().StringWithFallback(PScFocusSearch, "")
	scFocusSearchEntry.OnSubmitted(scFocusSearchStr)
	scFocusSearchEntry.SetText(scFocusSearchStr)
	s.commonSettingItems = append(s.commonSettingItems, scFocusSearchEntry)

	return container.NewWithoutLayout(s.commonSettingItems...)
}

func (s *Settings) GenFgSettings() *fyne.Container {
	// 糖豆人进程获取设置
	s.genGetFgPidSettingsRow()
	// 糖豆人自动点击坐标设置
	s.genFgAutoClickSettingsRow()

	return container.NewWithoutLayout(s.fgSettingItems...)
}

func (s *Settings) genGetFgPidSettingsRow() {
	y := lineHeight*0 + config.Padding
	app := fyne.CurrentApp()
	s.AutoFillMapId = app.Preferences().BoolWithFallback(PAutoFillMapId, false)
	cbAutoFillMapId := widget.NewCheckWithData(`点击"游玩"时，自动一键唤醒游戏，并填写地图代码`, binding.BindBool(&s.AutoFillMapId))
	cbAutoFillMapId.OnChanged = func(b bool) {
		s.AutoFillMapId = b
		app.Preferences().SetBool(PAutoFillMapId, s.AutoFillMapId)
	}
	cCbAutoFillMapId := container.NewHBox(cbAutoFillMapId)
	cCbAutoFillMapId.Move(fyne.NewPos(config.Padding, y))
	s.fgSettingItems = append(s.fgSettingItems, cCbAutoFillMapId)

	y += lineHeight
	infoLabel := canvas.NewText(`‼️ 如果自动获取不可用，请打开任务管理器查看并粘贴糖豆人进程ID（Pid）`, config.AccentColor)
	infoLabel.Move(fyne.NewPos(config.Padding, y))
	s.fgSettingItems = append(s.fgSettingItems, infoLabel)

	y += lineHeight * 0.5
	fgLabel := widget.NewLabel("糖豆人进程ID：")
	fgLabel.Alignment = fyne.TextAlignTrailing
	fgLabel.Resize(fyne.NewSize(fgLabelWidth, lineHeight))
	fgLabel.Move(fyne.NewPos(config.Padding, y))
	s.fgSettingItems = append(s.fgSettingItems, fgLabel)

	fgPidEntry := searchentry.NewSearchEntry("请填写糖豆人进程ID")
	fgPidEntry.Wrapping = fyne.TextTruncate
	fgPidEntry.Resize(fyne.NewSize(150, 35))
	fgPidLabel := widget.NewLabel("-")
	fgPidLabel.Resize(fyne.NewSize(200, lineHeight))
	fgPidEntry.OnCursorChanged = func() {
		s.FgPid = fgPidEntry.Text
	}
	fgPidEntry.OnSubmitted = func(str string) {
		s.FgPid = str
	}
	s.fgSettingItems = append(s.fgSettingItems, fgPidLabel)
	s.fgSettingItems = append(s.fgSettingItems, fgPidEntry)

	btnGetFgPid := widget.NewButtonWithIcon("自动获取", theme.SettingsIcon(), func() {
		fgPidLabel.SetText("")
		fgPidEntry.Hide()
		fgPidLabel.Show()
		state := false
		go func() {
			for {
				if state {
					break
				}
				if len(fgPidLabel.Text) > 5 {
					fgPidLabel.SetText(".")
				} else {
					fgPidLabel.SetText(fmt.Sprintf("%s.", fgPidLabel.Text))
				}
				time.Sleep(time.Second)
			}
		}()
		go func() {
			ids, err := robotgo.FindIds(FgName)
			if err == nil && len(ids) > 0 {
				pid := fmt.Sprintf("%d", ids[0])
				fgPidLabel.SetText(pid)
				fgPidEntry.SetText(pid)
				s.FgPid = pid
				app.SendNotification(&fyne.Notification{
					Title:   config.AppName,
					Content: "已获取到糖豆人进程ID",
				})
				state = true
				// if pid, err := strconv.ParseInt(s.FgPid, 10, 32); err == nil {
				// 	if err = robotgo.ActivePID(int32(pid)); err != nil {
				// 		app.SendNotification(&fyne.Notification{
				// 			Title:   config.AppName,
				// 			Content: "激活游戏窗口失败",
				// 		})
				// 	}
				// }
			} else {
				fgPidLabel.SetText("未获取到糖豆人进程ID")
				app.SendNotification(&fyne.Notification{
					Title:   config.AppName,
					Content: "未获取到糖豆人进程ID",
				})
				state = true
			}
			fgPidLabel.Hide()
			fgPidEntry.Show()
		}()
	})
	btnGetFgPid.Resize(fyne.NewSize(100, lineHeight/2))
	btnGetFgPid.Move(fyne.NewPos(fgLabel.Size().Width+config.Padding, y+3))
	s.BtnGetFgPid = btnGetFgPid
	fgPidLabel.Move(fyne.NewPos(fgLabel.Size().Width+btnGetFgPid.Size().Width+config.Padding*2, y))
	fgPidEntry.Move(fyne.NewPos(fgLabel.Size().Width+btnGetFgPid.Size().Width+config.Padding*2, y))
	fgPidLabel.Hide()
	s.fgSettingItems = append(s.fgSettingItems, btnGetFgPid)
}

func (s *Settings) genFgAutoClickSettingsRow() {
	app := fyne.CurrentApp()
	y := lineHeight*2.5 + config.Padding
	infoLabel := canvas.NewText(`‼️ 按照"123,123"格式在下方填写对应坐标（截图软件可以定位坐标），然后点击"测试点击"按钮调整坐标，保证能够正确点击相应按钮即可`, config.AccentColor)
	infoLabel.Move(fyne.NewPos(config.Padding, y))
	s.fgSettingItems = append(s.fgSettingItems, infoLabel)

	y += lineHeight * 0.5
	fgLabel := widget.NewLabel("自动点击：")
	fgLabel.Alignment = fyne.TextAlignTrailing
	fgLabel.Resize(fyne.NewSize(fgLabelWidth, lineHeight))
	fgLabel.Move(fyne.NewPos(config.Padding, y))
	s.fgSettingItems = append(s.fgSettingItems, fgLabel)

	// 选择节目按钮点击坐标文本框
	y += lineHeight * 0.5
	pSelectShowPos := app.Preferences().StringWithFallback(PSelectShowPos, "")
	selectShowPosEntry := searchentry.NewSearchEntry("x,y")
	selectShowPosEntry.Wrapping = fyne.TextTruncate
	s.PosSelectShow = &pSelectShowPos
	selectShowPosEntry.Bind(binding.BindString(s.PosSelectShow))
	selectShowPosEntry.Validator = nil
	selectShowPosEntry.Resize(fyne.NewSize(160, 35))
	selectShowPosEntry.Move(fyne.NewPos(config.Padding+fgLabel.Size().Width, y))
	s.fgSettingItems = append(s.fgSettingItems, selectShowPosEntry)
	selectShowPosLabel := widget.NewLabel("选择节目按钮")
	selectShowPosLabel.Alignment = fyne.TextAlignCenter
	selectShowPosLabel.Resize(selectShowPosEntry.Size())
	selectShowPosLabel.Move(fyne.NewPos(selectShowPosEntry.Position().X, selectShowPosEntry.Position().Y-selectShowPosEntry.Size().Height+5))
	s.fgSettingItems = append(s.fgSettingItems, selectShowPosLabel)

	// 输入代码按钮点击坐标文本框
	pEnterMapIdPos := app.Preferences().StringWithFallback(PEnterMapIdPos, "")
	enterMapIdPosEntry := searchentry.NewSearchEntry("x,y")
	enterMapIdPosEntry.Wrapping = fyne.TextTruncate
	s.PosEnterMapId = &pEnterMapIdPos
	enterMapIdPosEntry.Bind(binding.BindString(s.PosEnterMapId))
	enterMapIdPosEntry.Validator = nil
	enterMapIdPosEntry.Resize(fyne.NewSize(160, 35))
	enterMapIdPosEntry.Move(fyne.NewPos(config.Padding*2+fgLabel.Size().Width+selectShowPosEntry.Size().Width, y))
	s.fgSettingItems = append(s.fgSettingItems, enterMapIdPosEntry)
	enterMapIdPosLabel := widget.NewLabel("输入代码按钮")
	enterMapIdPosLabel.Alignment = fyne.TextAlignCenter
	enterMapIdPosLabel.Resize(enterMapIdPosEntry.Size())
	enterMapIdPosLabel.Move(fyne.NewPos(enterMapIdPosEntry.Position().X, enterMapIdPosEntry.Position().Y-enterMapIdPosEntry.Size().Height+5))
	s.fgSettingItems = append(s.fgSettingItems, enterMapIdPosLabel)

	// 代码输入框点击坐标文本框
	pCodeEntryPos := app.Preferences().StringWithFallback(PCodeEntryPos, "")
	codeEntryPosEntry := searchentry.NewSearchEntry("x,y")
	codeEntryPosEntry.Wrapping = fyne.TextTruncate
	s.PosCodeEntry = &pCodeEntryPos
	codeEntryPosEntry.Bind(binding.BindString(s.PosCodeEntry))
	codeEntryPosEntry.Validator = nil
	codeEntryPosEntry.Resize(fyne.NewSize(160, 35))
	codeEntryPosEntry.Move(fyne.NewPos(config.Padding*3+fgLabel.Size().Width+selectShowPosEntry.Size().Width+enterMapIdPosEntry.Size().Width, y))
	s.fgSettingItems = append(s.fgSettingItems, codeEntryPosEntry)
	codeEntryPosLabel := widget.NewLabel("代码输入框")
	codeEntryPosLabel.Alignment = fyne.TextAlignCenter
	codeEntryPosLabel.Resize(codeEntryPosEntry.Size())
	codeEntryPosLabel.Move(fyne.NewPos(codeEntryPosEntry.Position().X, codeEntryPosEntry.Position().Y-codeEntryPosEntry.Size().Height+5))
	s.fgSettingItems = append(s.fgSettingItems, codeEntryPosLabel)

	// 确认按钮点击坐标文本框
	pConfirmBtnPos := app.Preferences().StringWithFallback(PConfirmBtnPos, "")
	confirmBtnPosEntry := searchentry.NewSearchEntry("x,y")
	confirmBtnPosEntry.Wrapping = fyne.TextTruncate
	s.PosConfirmBtn = &pConfirmBtnPos
	confirmBtnPosEntry.Bind(binding.BindString(s.PosConfirmBtn))
	confirmBtnPosEntry.Validator = nil
	confirmBtnPosEntry.Resize(fyne.NewSize(160, 35))
	confirmBtnPosEntry.Move(fyne.NewPos(config.Padding*4+fgLabel.Size().Width+selectShowPosEntry.Size().Width+enterMapIdPosEntry.Size().Width+codeEntryPosEntry.Size().Width, y))
	s.fgSettingItems = append(s.fgSettingItems, confirmBtnPosEntry)
	confirmBtnPosLabel := widget.NewLabel("确认按钮")
	confirmBtnPosLabel.Alignment = fyne.TextAlignCenter
	confirmBtnPosLabel.Resize(confirmBtnPosEntry.Size())
	confirmBtnPosLabel.Move(fyne.NewPos(confirmBtnPosEntry.Position().X, confirmBtnPosEntry.Position().Y-confirmBtnPosEntry.Size().Height+5))
	s.fgSettingItems = append(s.fgSettingItems, confirmBtnPosLabel)
	logger.Debugf("Binding Select Show: %s EnterMapId: %s CodeEntry: %s ConfirmBtn: %s", pSelectShowPos, pEnterMapIdPos, pCodeEntryPos, pConfirmBtnPos)

	// 测试点击按钮
	var clickedEntry *searchentry.SearchEntry
	btnTestPos := widget.NewButtonWithIcon("测试点击", theme.ContentAddIcon(), func() {
		logger.Debugf("Select Show: %s EnterMapId: %s CodeEntry: %s", s.PosSelectShow, s.PosEnterMapId, s.PosCodeEntry)
		if clickedEntry == nil {
			dialog.ShowInformation("提示", "请点击一个坐标文本框再试", *s.Window)
			return
		}
		pos := strings.Split(clickedEntry.Text, ",")
		if len(pos) != 2 {
			dialog.ShowInformation("提示", "坐标格式错误，格式：123,123", *s.Window)
			return
		}
		var xx int
		var yy int
		var err error
		xx, err = strconv.Atoi(pos[0])
		if err != nil {
			dialog.ShowInformation("提示", "坐标格式错误，格式：123,123", *s.Window)
			return
		}
		yy, err = strconv.Atoi(pos[1])
		if err != nil {
			dialog.ShowInformation("提示", "坐标格式错误，格式：123,123", *s.Window)
			return
		}
		go func() {
			if s.FgPid != "" {
				if pid, err := strconv.ParseInt(s.FgPid, 10, 32); err == nil {
					if err = robotgo.ActivePID(int32(pid)); err != nil {
						dialog.ShowInformation("提示", "唤醒游戏失败，请重新获取游戏进程", *s.Window)
					}
				}
			}
			robotgo.Move(xx, yy)
			time.Sleep(500 * time.Millisecond)
			robotgo.Click("left")
		}()
	})
	btnTestPos.Resize(fyne.NewSize(100, 35))
	btnTestPos.Move(fyne.NewPos(config.Padding*5+fgLabel.Size().Width+selectShowPosEntry.Size().Width+enterMapIdPosEntry.Size().Width+codeEntryPosEntry.Size().Width+confirmBtnPosEntry.Size().Width, y))
	s.fgSettingItems = append(s.fgSettingItems, btnTestPos)

	selectShowPosEntry.OnTapped = func(event *fyne.PointEvent) {
		clickedEntry = selectShowPosEntry
	}
	enterMapIdPosEntry.OnTapped = func(event *fyne.PointEvent) {
		clickedEntry = enterMapIdPosEntry
	}
	codeEntryPosEntry.OnTapped = func(event *fyne.PointEvent) {
		clickedEntry = codeEntryPosEntry
	}
	confirmBtnPosEntry.OnTapped = func(event *fyne.PointEvent) {
		clickedEntry = confirmBtnPosEntry
	}
	selectShowPosEntry.OnCursorChanged = func() {
		fyne.CurrentApp().Preferences().SetString(PSelectShowPos, selectShowPosEntry.Text)
	}
	enterMapIdPosEntry.OnCursorChanged = func() {
		fyne.CurrentApp().Preferences().SetString(PEnterMapIdPos, enterMapIdPosEntry.Text)
	}
	codeEntryPosEntry.OnCursorChanged = func() {
		fyne.CurrentApp().Preferences().SetString(PCodeEntryPos, codeEntryPosEntry.Text)
	}
	confirmBtnPosEntry.OnCursorChanged = func() {
		fyne.CurrentApp().Preferences().SetString(PConfirmBtnPos, confirmBtnPosEntry.Text)
	}

	// 快捷键设置
	// 根据给定参数生成连线
	objects := []fyne.CanvasObject{selectShowPosEntry, enterMapIdPosEntry, codeEntryPosEntry, confirmBtnPosEntry}
	playNext := func(str string) {
		keys := strings.Split(str, "+")
		if !s.shortcutChecker(keys, s.isNotify, 3) {
			return
		}
		go func() {
			var modifiers []hotkey.Modifier
			for i := 0; i < len(keys)-1; i++ {
				modifiers = append(modifiers, hotkeys.GetModifier(keys[i]))
			}
			s.registerHotKey(modifiers, s.getKey(keys[len(keys)-1]), func() {
				go func() {
					maps, _ := db.ListMap(1, 1, &model.MapInfo{State: "0"}, `created asc, map_id`)
					if len(maps) > 0 {
						db.UpdateMap(
							model.MapInfo{MapId: maps[0].MapId, State: "1", PlayTime: time.Now()},
							[]string{"State", "PlayTime"},
							&model.MapInfo{State: "0"})
						cbm.CallBackFunc("fg FillMapIdForPlayNext", maps[0].MapId, s)
					}
				}()
			}, s.isNotify, PHotKeyPlayNext)
		}()
	}
	playNextHotKeyStr := app.Preferences().StringWithFallback(PHotKeyPlayNext, "")
	playNext(playNextHotKeyStr)
	s.genEntryWithLink(objects[1:], fmt.Sprintf("例：ctrl+%s+n，回车保存", s.altName), 170, 35, playNext, 2, playNextHotKeyStr)
	playNextOnSelectShow := func(str string) {
		keys := strings.Split(str, "+")
		if !s.shortcutChecker(keys, s.isNotify, 3) {
			return
		}
		go func() {
			var modifiers []hotkey.Modifier
			for i := 0; i < len(keys)-1; i++ {
				modifiers = append(modifiers, hotkeys.GetModifier(keys[i]))
			}
			s.registerHotKey(modifiers, s.getKey(keys[len(keys)-1]), func() {
				go func() {
					maps, _ := db.ListMap(1, 1, &model.MapInfo{State: "0"}, `created asc, map_id`)
					if len(maps) > 0 {
						db.UpdateMap(
							model.MapInfo{MapId: maps[0].MapId, State: "1", PlayTime: time.Now()},
							[]string{"State", "PlayTime"},
							&model.MapInfo{State: "0"})
						cbm.CallBackFunc("fg FillMapIdForPlayNextOnSelectShow", maps[0].MapId, s)
					}
				}()
			}, s.isNotify, PHotKeyPlayNextOnSelectShow)
		}()
	}
	playNextOnSelectShowHotKeyStr := app.Preferences().StringWithFallback(PHotKeyPlayNextOnSelectShow, "")
	playNextOnSelectShow(playNextOnSelectShowHotKeyStr)
	// playNext(playNextOnSelectShowHotKeyStr)
	s.genEntryWithLink(objects, fmt.Sprintf("例：ctrl+%s+p，回车保存", s.altName), 170, 35, playNextOnSelectShow, 1, playNextOnSelectShowHotKeyStr)
	go func() {
		time.Sleep(2 * time.Second)
		*s.isNotify = true
	}()
}

func (s *Settings) shortcutChecker(keys []string, isNotify *bool, keyNumber int) bool {
	if len(keys) < keyNumber {
		if *isNotify {
			dialog.ShowInformation("提示", fmt.Sprintf("快捷键必须至少%d个按键", keyNumber), *s.Window)
		}
		return false
	}
	return true
}

func (s *Settings) registerHotKey(modifiers []hotkey.Modifier, key hotkey.Key, onPress func(), isNotify *bool, preferencesKey string) {
	modifier := ""
	for _, modifierKey := range modifiers {
		modifier = fmt.Sprintf("%s+%s", modifier, hotkeys.GetModifierName(modifierKey))
	}
	if len(modifiers) > 0 {
		modifier = modifier[1:]
	}
	hotKeyStr := fmt.Sprintf("%s+%s", modifier, s.getKeyName(key))
	if s.hotKey[hotKeyStr] != nil {
		_ = s.hotKey[hotKeyStr].Unregister()
	}
	hk := hotkey.New(modifiers, key)
	if err := hk.Register(); err != nil {
		dialog.ShowInformation("提示", fmt.Sprintf("快捷键设置失败：%v", err), *s.Window)
		return
	} else {
		if *isNotify {
			fyne.CurrentApp().Preferences().SetString(preferencesKey, hotKeyStr)
			dialog.ShowInformation("提示", fmt.Sprintf("快捷键[%s]设置成功", hotKeyStr), *s.Window)
		}
		s.propertyLock.Lock()
		s.hotKey[hotKeyStr] = hk
		s.propertyLock.Unlock()
	}
	for range hk.Keydown() {
		onPress()
	}
}

func (s *Settings) genEntryWithLink(objects []fyne.CanvasObject, placeHolder string, width float32, height float32, onSubmit func(s string), lineNo float32, defaultValue string) {
	// 根据objects计算连线起点坐标
	linkStartEndpoint := canvas.NewRectangle(config.ShadowColor)
	linkStartEndpoint.Resize(fyne.NewSize(5, height*lineNo))
	linkStartEndpoint.Move(fyne.NewPos(objects[0].Position().X+objects[0].Size().Width/2, objects[0].Position().Y+objects[0].Size().Height))
	s.fgSettingItems = append(s.fgSettingItems, linkStartEndpoint)
	linkEndEndpoint := canvas.NewRectangle(config.ShadowColor)
	linkEndEndpoint.Resize(fyne.NewSize(5, height*lineNo))
	linkEndEndpoint.Move(fyne.NewPos(objects[len(objects)-1].Position().X+objects[len(objects)-1].Size().Width/2, objects[len(objects)-1].Position().Y+objects[len(objects)-1].Size().Height))
	s.fgSettingItems = append(s.fgSettingItems, linkEndEndpoint)

	link := canvas.NewRectangle(config.ShadowColor)
	// var linkWidth float32
	// for i := 0; i < len(objects); i++ {
	// 	switch i {
	// 	case 0:
	// 		linkWidth += objects[i].Size().Width/2
	// 	case len(objects) - 1:
	// 		linkWidth += objects[i].Size().Width/2 + config.Padding
	// 	default:
	// 		linkWidth += objects[i].Size().Width
	// 	}
	// }
	link.Resize(fyne.NewSize(linkEndEndpoint.Position().X-linkStartEndpoint.Position().X+linkEndEndpoint.Size().Width, 5))
	link.Move(fyne.NewPos(linkStartEndpoint.Position().X, linkStartEndpoint.Position().Y+linkStartEndpoint.Size().Height))
	s.fgSettingItems = append(s.fgSettingItems, link)

	// 生成文本框
	entry := searchentry.NewSearchEntry(placeHolder)
	entry.Wrapping = fyne.TextTruncate
	entry.Resize(fyne.NewSize(width, height))
	entry.Move(fyne.NewPos(link.Position().X+link.Size().Width/2-width/2, link.Position().Y+link.Size().Height-height/2-2.5))
	entry.OnSubmitted = onSubmit
	entry.OnCursorChanged = func() {
		linkStartEndpoint.FillColor = config.ShadowColor
		linkEndEndpoint.FillColor = config.ShadowColor
		link.FillColor = config.ShadowColor
		linkStartEndpoint.Refresh()
		linkEndEndpoint.Refresh()
		link.Refresh()
	}
	entry.OnTapped = func(event *fyne.PointEvent) {
		linkStartEndpoint.FillColor = config.AccentColor
		linkEndEndpoint.FillColor = config.AccentColor
		link.FillColor = config.AccentColor
		linkStartEndpoint.Refresh()
		linkEndEndpoint.Refresh()
		link.Refresh()
	}
	entry.SetText(defaultValue)
	s.fgSettingItems = append(s.fgSettingItems, entry)
}

func (s *Settings) getKey(key string) hotkey.Key {
	switch strings.ToLower(key) {
	case "a":
		return hotkey.Key(hotkey.KeyA)
	case "b":
		return hotkey.Key(hotkey.KeyB)
	case "c":
		return hotkey.Key(hotkey.KeyC)
	case "d":
		return hotkey.Key(hotkey.KeyD)
	case "e":
		return hotkey.Key(hotkey.KeyE)
	case "f":
		return hotkey.Key(hotkey.KeyF)
	case "g":
		return hotkey.Key(hotkey.KeyG)
	case "h":
		return hotkey.Key(hotkey.KeyH)
	case "i":
		return hotkey.Key(hotkey.KeyI)
	case "j":
		return hotkey.Key(hotkey.KeyJ)
	case "k":
		return hotkey.Key(hotkey.KeyK)
	case "l":
		return hotkey.Key(hotkey.KeyL)
	case "m":
		return hotkey.Key(hotkey.KeyM)
	case "n":
		return hotkey.Key(hotkey.KeyN)
	case "o":
		return hotkey.Key(hotkey.KeyO)
	case "p":
		return hotkey.Key(hotkey.KeyP)
	case "q":
		return hotkey.Key(hotkey.KeyQ)
	case "r":
		return hotkey.Key(hotkey.KeyR)
	case "s":
		return hotkey.Key(hotkey.KeyS)
	case "t":
		return hotkey.Key(hotkey.KeyT)
	case "u":
		return hotkey.Key(hotkey.KeyU)
	case "v":
		return hotkey.Key(hotkey.KeyV)
	case "w":
		return hotkey.Key(hotkey.KeyW)
	case "x":
		return hotkey.Key(hotkey.KeyX)
	case "y":
		return hotkey.Key(hotkey.KeyY)
	default:
		return hotkey.Key(hotkey.KeyZ)
	}
}

func (s *Settings) getKeyForFyne(key string) fyne.KeyName {
	switch strings.ToLower(key) {
	case "a":
		return fyne.KeyA
	case "b":
		return fyne.KeyB
	case "c":
		return fyne.KeyC
	case "d":
		return fyne.KeyD
	case "e":
		return fyne.KeyE
	case "f":
		return fyne.KeyF
	case "g":
		return fyne.KeyG
	case "h":
		return fyne.KeyH
	case "i":
		return fyne.KeyI
	case "j":
		return fyne.KeyJ
	case "k":
		return fyne.KeyK
	case "l":
		return fyne.KeyL
	case "m":
		return fyne.KeyM
	case "n":
		return fyne.KeyN
	case "o":
		return fyne.KeyO
	case "p":
		return fyne.KeyP
	case "q":
		return fyne.KeyQ
	case "r":
		return fyne.KeyR
	case "s":
		return fyne.KeyS
	case "t":
		return fyne.KeyT
	case "u":
		return fyne.KeyU
	case "v":
		return fyne.KeyV
	case "w":
		return fyne.KeyW
	case "x":
		return fyne.KeyX
	case "y":
		return fyne.KeyY
	default:
		return fyne.KeyZ
	}
}

func (s *Settings) getKeyName(key hotkey.Key) string {
	switch key {
	case hotkey.Key(hotkey.KeyA):
		return "a"
	case hotkey.Key(hotkey.KeyB):
		return "b"
	case hotkey.Key(hotkey.KeyC):
		return "c"
	case hotkey.Key(hotkey.KeyD):
		return "d"
	case hotkey.Key(hotkey.KeyE):
		return "e"
	case hotkey.Key(hotkey.KeyF):
		return "f"
	case hotkey.Key(hotkey.KeyG):
		return "g"
	case hotkey.Key(hotkey.KeyH):
		return "h"
	case hotkey.Key(hotkey.KeyI):
		return "i"
	case hotkey.Key(hotkey.KeyJ):
		return "j"
	case hotkey.Key(hotkey.KeyK):
		return "k"
	case hotkey.Key(hotkey.KeyL):
		return "l"
	case hotkey.Key(hotkey.KeyM):
		return "m"
	case hotkey.Key(hotkey.KeyN):
		return "n"
	case hotkey.Key(hotkey.KeyO):
		return "o"
	case hotkey.Key(hotkey.KeyP):
		return "p"
	case hotkey.Key(hotkey.KeyQ):
		return "q"
	case hotkey.Key(hotkey.KeyR):
		return "r"
	case hotkey.Key(hotkey.KeyS):
		return "s"
	case hotkey.Key(hotkey.KeyT):
		return "t"
	case hotkey.Key(hotkey.KeyU):
		return "u"
	case hotkey.Key(hotkey.KeyV):
		return "v"
	case hotkey.Key(hotkey.KeyW):
		return "w"
	case hotkey.Key(hotkey.KeyX):
		return "x"
	case hotkey.Key(hotkey.KeyY):
		return "y"
	default:
		return "z"
	}
}

func (s *Settings) getModifier(key string) fyne.KeyModifier {
	switch strings.ToLower(key) {
	case "cmd":
		return fyne.KeyModifierShortcutDefault
	case "ctrl":
		return fyne.KeyModifierControl
	case "option", "alt":
		return fyne.KeyModifierAlt
	default:
		return fyne.KeyModifierShift
	}
}
