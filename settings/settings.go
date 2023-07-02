package settings

import (
	"FallGuys66/config"
	"FallGuys66/data"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	"time"
)

var (
	lineHeight    = float32(50)
	PAutoGetFgPid = "PAutoGetFgPid"
	FgName        = "Sublime"
)

type Settings struct {
	FgPid        string
	BtnGetFgPid  *widget.Button
	AutoGetFgPid bool

	commonSettingItems []fyne.CanvasObject
	fgSettingItems     []fyne.CanvasObject
}

func NewSettings() *Settings {
	return &Settings{}
}

func (s *Settings) Init() *container.AppTabs {
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
	// 糖豆人窗口捕获
	fgLabel := widget.NewLabel("糖豆人进程ID：")
	y1st := lineHeight*0 + config.Padding
	fgLabel.Resize(fyne.NewSize(100, lineHeight))
	fgLabel.Move(fyne.NewPos(config.Padding, y1st))
	s.fgSettingItems = append(s.fgSettingItems, fgLabel)

	app := fyne.CurrentApp()
	fgPidLabel := widget.NewLabel("-")
	fgPidLabel.Resize(fyne.NewSize(200, lineHeight))
	s.fgSettingItems = append(s.fgSettingItems, fgPidLabel)

	btnGetFgPid := widget.NewButtonWithIcon("自动获取", theme.SettingsIcon(), func() {
		fgPidLabel.SetText("")
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
				s.FgPid = pid
				app.SendNotification(&fyne.Notification{
					Title:   config.AppName,
					Content: "已获取到糖豆人进程ID\n开启一键填充地图ID",
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
				fgPidLabel.SetText("未获取到糖豆人进程ID，无法自动填写地图ID")
				app.SendNotification(&fyne.Notification{
					Title:   config.AppName,
					Content: "未获取到糖豆人进程ID\n无法自动填写地图ID",
				})
				state = true
			}
		}()
	})
	btnGetFgPid.Resize(fyne.NewSize(100, lineHeight/2))
	btnGetFgPid.Move(fyne.NewPos(fgLabel.Size().Width+config.Padding, y1st+3))
	s.BtnGetFgPid = btnGetFgPid
	fgPidLabel.Move(fyne.NewPos(fgLabel.Size().Width+btnGetFgPid.Size().Width+config.Padding*2, y1st))
	s.fgSettingItems = append(s.fgSettingItems, btnGetFgPid)

	return container.NewWithoutLayout(s.fgSettingItems...)
}
