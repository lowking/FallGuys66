package settings

import (
	"FallGuys66/config"
	"FallGuys66/data"
	"FallGuys66/live/douyu/lib/logger"
	"FallGuys66/widgets/searchentry"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"strconv"
	"strings"
	"time"
)

var (
	lineHeight     = float32(50)
	PAutoGetFgPid  = "PAutoGetFgPid"
	PAutoFillMapId = "PAutoFillMapId"
	PSelectShowPos = "PSelectShowPos"
	PEnterMapIdPos = "PEnterMapIdPos"
	PCodeEntryPos  = "PCodeEntryPos"
	PConfirmBtnPos = "PConfirmBtnPos"
	FgName         = "FallGuys_client"
	fgLabelWidth   = float32(100)
)

type Settings struct {
	FgPid         string
	BtnGetFgPid   *widget.Button
	AutoGetFgPid  bool
	AutoFillMapId bool
	Window        *fyne.Window

	PosSelectShow *string
	PosEnterMapId *string
	PosCodeEntry  *string
	PosConfirmBtn *string

	commonSettingItems []fyne.CanvasObject
	fgSettingItems     []fyne.CanvasObject
}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) Init(window *fyne.Window) *container.AppTabs {
	s.Window = window
	settingTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("通用", theme.SettingsIcon(), s.GenCommonSettings()),
		container.NewTabItemWithIcon("　糖豆人　", data.FgLogo, s.GenFgSettings()),
	)
	settingTabs.SetTabLocation(container.TabLocationLeading)

	return settingTabs
}

func (s *Settings) GenCommonSettings() *fyne.Container {
	app := fyne.CurrentApp()
	s.AutoGetFgPid = app.Preferences().BoolWithFallback(PAutoGetFgPid, false)
	y1st := lineHeight*0 + config.Padding
	cbAutoGetFgPid := widget.NewCheckWithData("启动时，自动获取糖豆人进程ID", binding.BindBool(&s.AutoGetFgPid))
	cbAutoGetFgPid.OnChanged = func(b bool) {
		s.AutoGetFgPid = b
		app.Preferences().SetBool(PAutoGetFgPid, s.AutoGetFgPid)
	}
	cbAutoGetFgPid.SetChecked(s.AutoGetFgPid)
	cCbAutoGetFgPid := container.NewHBox(cbAutoGetFgPid)
	cCbAutoGetFgPid.Move(fyne.NewPos(config.Padding, y1st))
	s.commonSettingItems = append(s.commonSettingItems, cCbAutoGetFgPid)

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

	y = lineHeight*1 + config.Padding
	infoLabel := canvas.NewText(`‼️ 如果自动获取不可用，请打开任务管理器查看并粘贴糖豆人进程ID（Pid）`, config.AccentColor)
	infoLabel.Move(fyne.NewPos(config.Padding, y))
	s.fgSettingItems = append(s.fgSettingItems, infoLabel)

	y = lineHeight*1.5 + config.Padding
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

	y = lineHeight*3 + config.Padding
	fgLabel := widget.NewLabel("自动点击：")
	fgLabel.Alignment = fyne.TextAlignTrailing
	fgLabel.Resize(fyne.NewSize(fgLabelWidth, lineHeight))
	fgLabel.Move(fyne.NewPos(config.Padding, y))
	s.fgSettingItems = append(s.fgSettingItems, fgLabel)

	// 选择节目按钮点击坐标文本框
	y = lineHeight*3.5 + config.Padding
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
	s.genEntryWithLink(objects[1:], "例：ctrl+alt+n，回车保存", 150, 35, func(str string) {
		keys := strings.Split(str, "+")
		if len(keys) < 3 {
			dialog.ShowInformation("提示", "快捷键必须至少3个按键", *s.Window)
			return
		}
	}, 2)
	s.genEntryWithLink(objects, "例：ctrl+alt+p，回车保存", 150, 35, func(str string) {
		keys := strings.Split(str, "+")
		if len(keys) < 3 {
			dialog.ShowInformation("提示", "快捷键必须至少3个按键", *s.Window)
			return
		}
	}, 1)
}

func (s *Settings) genEntryWithLink(objects []fyne.CanvasObject, placeHolder string, width float32, height float32, onSubmit func(s string), lineNo float32) {
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
	s.fgSettingItems = append(s.fgSettingItems, entry)
}
