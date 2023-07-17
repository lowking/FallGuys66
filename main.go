//go:generate rsrc -ico data/assets/logo-white.ico -manifest data/assets/main.exe.manifest -o data/assets/main.syso
package main

import (
	"FallGuys66/config"
	"FallGuys66/data"
	_ "FallGuys66/db"
	"FallGuys66/db/model"
	"FallGuys66/handler"
	"FallGuys66/live/bilibili"
	"FallGuys66/live/douyu/DMconfig/config"
	"FallGuys66/live/douyu/DYtype"
	"FallGuys66/live/douyu/client"
	"FallGuys66/live/douyu/lib/logger"
	"FallGuys66/settings"
	"FallGuys66/utils"
	"FallGuys66/widgets/searchentry"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/tutorials"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	lock "github.com/viney-shih/go-lock"
)

const preferenceCurrentTutorial = "currentTutorial"
const PDefaultLiveHostNo = "DefaultLiveHostNo"
const PLiveHosts = "LiveHosts"
const PLivePlatform = "PLivePlatform"

var topWindow fyne.Window
var tabs *container.AppTabs
var setting *settings.Settings
var version = "1.7.0"
var driver fyne.Driver
var window fyne.Window
var biliClient *gobilibili.BiliBiliClient
var enableQueryLoop = true
var connectLock = lock.NewCASMutex()
var searchLock = lock.NewCASMutex()

func main() {
	offsetX := float32(-50)
	offsetY := float32(3)
	offset := float32(-6)
	application := app.NewWithID("pro.lowking.fallguys66")
	application.Settings().SetTheme(&data.MyTheme{
		Regular:    data.FontSmileySansOblique,
		Bold:       data.FontSmileySansOblique,
		Italic:     data.FontSmileySansOblique,
		BoldItalic: data.FontSmileySansOblique,
		Monospace:  data.FontSmileySansOblique,
	})
	application.SetIcon(data.LogoWhite)
	driver = application.Driver()
	// 托盘图标
	makeTray(application)
	handler.WhereMap[0] = &model.MapInfo{State: "0"}
	handler.OrderMap[0] = `created asc, map_id`
	handler.WhereMap[1] = &model.MapInfo{State: "1"}
	handler.OrderMap[1] = `play_time desc, map_id`
	handler.WhereMap[2] = &model.MapInfo{Star: "1"}
	handler.OrderMap[2] = `created desc, map_id`
	handler.WhereMap[3] = &model.Blacklist{}
	handler.OrderMap[3] = `created desc`
	handler.WhereMap[4] = &model.MapInfo{}
	handler.OrderMap[4] = `created desc, map_id`
	logLifecycle(application)
	window = application.NewWindow(config.AppName)
	topWindow = window

	// 菜单
	// window.SetMainMenu(makeMenu(application, window))
	window.SetMaster()

	logo := canvas.NewImageFromResource(data.LogoWhite)
	logo.SetMinSize(fyne.NewSize(config.LogoSize, config.LogoSize))
	logoBlack := canvas.NewImageFromResource(data.LogoBlack)
	logoBlack.SetMinSize(fyne.NewSize(config.LogoSize, config.LogoSize))

	cLogo := container.NewCenter(logo)
	cLogo.Add(logoBlack)
	cLogo.Resize(fyne.NewSize(100, 100))
	cLogo.Move(fyne.NewPos(config.Padding, config.Padding+10))
	logoBlack.Hide()
	go func() {
		for {
			flashEle(rand.Intn(5), logo, logoBlack)
			randTime := rand.Intn(900) + 100
			time.Sleep(time.Duration(randTime) * time.Millisecond)
		}
	}()

	// 需要渲染的元素
	var elements []fyne.CanvasObject
	logger.Info("Loading background ...")
	bg := canvas.NewImageFromResource(data.Bg)
	bgHeight := cLogo.Size().Height + config.Padding*3 + offset + 60
	bg.Resize(fyne.NewSize(config.AppSize.Width, bgHeight))
	bg.Move(fyne.NewPos(offset, offset))
	elements = append(elements, bg)
	elements = append(elements, cLogo)

	// 直播间label
	logger.Info("Initializing live host label ...")
	lLiveHostLabel := canvas.NewText("直播间：", theme.PrimaryColor())
	lLiveHostShadowLabel := canvas.NewText("直播间：", config.ShadowColor)
	lLiveHostLabel.TextSize = 35
	lLiveHostShadowLabel.TextSize = 35
	lLiveHostLabel.Move(fyne.NewPos(config.ToolbarPaddingLeft, config.ToolbarPaddingTop))
	lLiveHostShadowLabel.Move(fyne.NewPos(config.ToolbarPaddingLeft+1, config.ToolbarPaddingTop+1))
	elements = append(elements, lLiveHostShadowLabel)
	elements = append(elements, lLiveHostLabel)

	optionSize := fyne.NewSize(140, 37)
	// 直播间输入框
	logger.Info("Initializing live host number label ...")
	defaultLiveHostNo := application.Preferences().StringWithFallback(PDefaultLiveHostNo, "")
	tLiveHostNo := canvas.NewText(defaultLiveHostNo, color.RGBA{
		R: 106,
		G: 135,
		B: 89,
		A: 255,
	})
	tLiveHostNo.TextSize = 35
	tLiveHostNo.Move(fyne.NewPos(cLogo.Size().Width+config.ToolbarPaddingLeft+optionSize.Width+config.Padding*2+offsetX, config.ToolbarPaddingTop+offsetY-3))
	tLiveHostNo.Hide()
	elements = append(elements, tLiveHostNo)
	liveHosts := []string{
		"156277",
	}
	pLiveHosts := application.Preferences().StringWithFallback(PLiveHosts, "")
	if pLiveHosts != "" {
		liveHosts = strings.Split(pLiveHosts, "|")
	}
	logger.Info("Initializing live host selector ...")
	// bindLiveHosts := binding.BindStringList(&liveHosts)
	liveHostOption := widget.NewSelectEntry(liveHosts)
	liveHostOption.TextStyle.Bold = true
	liveHostOption.SetPlaceHolder("请输入/选择直播间号")
	liveHost := ""
	liveHostOption.OnChanged = func(s string) {
		liveHost = s
	}
	liveHostOption.SetText(defaultLiveHostNo)
	cLiveHostOption := container.NewMax(liveHostOption)
	cLiveHostOption.Resize(optionSize)
	cLiveHostOption.Move(fyne.NewPos(cLogo.Size().Width+config.ToolbarPaddingLeft+optionSize.Width+config.Padding*2+offsetX, config.ToolbarPaddingTop+offsetY))
	elements = append(elements, cLiveHostOption)

	// 初始化平台选项
	logger.Info("Initializing live platform radio ...")
	livePlatformOptions := []string{"Douyu", "Bilibili"}
	livePlatform := application.Preferences().StringWithFallback(PLivePlatform, "Douyu")
	rLivePlatform := widget.NewRadioGroup(livePlatformOptions, func(p string) {
		livePlatform = p
		application.Preferences().SetString(PLivePlatform, p)
	})
	rLivePlatform.SetSelected(livePlatform)
	rLivePlatform.Move(fyne.NewPos(cLiveHostOption.Position().X+cLiveHostOption.Size().Width+config.Padding, cLiveHostOption.Position().Y))
	elements = append(elements, rLivePlatform)

	// 初始化copyright
	logger.Info("Initializing copyright info ...")
	lCopyrightL := canvas.NewText("斗鱼ID：石疯悦耳", color.RGBA{
		R: 32,
		G: 32,
		B: 32,
		A: 100,
	})
	lCopyrightL.TextSize = 14
	// lCopyrightL.Move(fyne.NewPos(0, appSize.Height-25))
	lCopyrightL.Alignment = fyne.TextAlignTrailing
	lCopyrightL.Move(fyne.NewPos(config.AppSize.Width-10, 140))
	lCopyrightR := canvas.NewText("© lowking 2023. All Rights Reserved.", color.RGBA{
		R: 32,
		G: 32,
		B: 32,
		A: 100,
	})
	lCopyrightR.TextSize = 14
	lCopyrightR.Alignment = fyne.TextAlignTrailing
	// lCopyrightR.Move(fyne.NewPos(appSize.Width-10, appSize.Height-25))
	lCopyrightR.Move(fyne.NewPos(config.AppSize.Width-130, 140))
	elements = append(elements, lCopyrightL)
	elements = append(elements, lCopyrightR)

	logger.Info("Initializing table stuff ...")
	var toggleElements []fyne.CanvasObject
	// 给列表行添加背景
	rowHeight := float32(36.2)
	startY := cLogo.Size().Height + 137
	for i := float32(0); i+2 < (config.AppSize.Height-startY)/rowHeight; i += 2 {
		tableRowLine := canvas.NewRectangle(theme.HoverColor())
		tableRowLine.Resize(fyne.NewSize(config.AppSize.Width+10, rowHeight))
		tableRowLine.Move(fyne.NewPos(-10, startY+rowHeight*i))
		toggleElements = append(toggleElements, tableRowLine)
	}

	// 直线分割表头
	tableHeaderLine := canvas.NewRectangle(theme.ShadowColor())
	tableHeaderLine.Resize(fyne.NewSize(config.AppSize.Width+20, 6))
	tableHeaderLine.Move(fyne.NewPos(-10, cLogo.Size().Height+98))
	toggleElements = append(toggleElements, tableHeaderLine)
	tableFooterLine := canvas.NewRectangle(theme.ShadowColor())
	tableFooterLine.Resize(fyne.NewSize(config.AppSize.Width+20, 6))
	tableFooterLine.Move(fyne.NewPos(-10, cLogo.Size().Height+98))
	tableFooterLine.Move(fyne.NewPos(-10, config.AppSize.Height-47))
	toggleElements = append(toggleElements, tableFooterLine)

	elements = append(elements, toggleElements...)

	// 初始化tab列表
	logger.Info("Initializing tab list ...")
	setting = settings.NewSettings()
	tabs = container.NewAppTabs(
		container.NewTabItem("未玩列表", utils.MakeEmptyList(theme.PrimaryColor())),
		container.NewTabItem("已玩列表", utils.MakeEmptyList(theme.PrimaryColor())),
		container.NewTabItem("收藏列表", utils.MakeEmptyList(theme.PrimaryColor())),
		container.NewTabItem("黑名单", utils.MakeEmptyList(theme.PrimaryColor())),
		container.NewTabItem("搜索结果", utils.MakeEmptyList(theme.PrimaryColor())),
		container.NewTabItem("设置", setting.Init(&window, &driver)),
	)
	tabs.OnSelected = func(item *container.TabItem) {
		toggleElementsVisible := true
		switch item.Text {
		case "未玩列表":
			go handler.RefreshMapList(setting, window, tabs, 0, "", handler.WhereMap[0], handler.OrderMap[0], false, false)
		case "已玩列表":
			go handler.RefreshMapList(setting, window, tabs, 1, "", handler.WhereMap[1], handler.OrderMap[1], false, false)
		case "收藏列表":
			go handler.RefreshMapList(setting, window, tabs, 2, "", handler.WhereMap[2], handler.OrderMap[2], false, false)
		case "黑名单":
			go handler.RefreshMapList(setting, window, tabs, 3, "", handler.WhereMap[3], handler.OrderMap[3], false, false)
		case "搜索结果":
			go handler.RefreshMapList(setting, window, tabs, 4, "", handler.WhereMap[4], handler.OrderMap[4], false, false)
		default:
			toggleElementsVisible = false
		}
		for _, element := range toggleElements {
			if toggleElementsVisible {
				element.Show()
			} else {
				element.Hide()
			}
		}
	}
	cTabList := container.NewBorder(nil, nil, nil, nil, tabs)
	cTabList.Move(fyne.NewPos(offset, cLogo.Size().Height+config.Padding*3+offset))
	cTabList.Resize(fyne.NewSize(config.AppSize.Width, config.AppSize.Height-cLogo.Size().Height-config.Padding*3))
	elements = append(elements, cTabList)

	// 连接直播间按钮
	logger.Info("Initializing connect button ...")
	var webSocketClient client.DyBarrageWebSocketClient
	var btnCon *widget.Button
	btnCon = widget.NewButtonWithIcon("连接", theme.NavigateNextIcon(), func() {
		if !connectLock.TryLockWithTimeout(100 * time.Millisecond) {
			return
		}
		defer connectLock.Unlock()
		if liveHost == "" {
			btnConSetDefault(btnCon, liveHostOption, rLivePlatform, tLiveHostNo)
			dialog.ShowInformation("提示", "请填写直播间号", window)
			return
		} else if _, err := strconv.Atoi(liveHost); err != nil && liveHost != "dev" {
			btnConSetDefault(btnCon, liveHostOption, rLivePlatform, tLiveHostNo)
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
			rLivePlatform.Disable()
			connectSuc := false
			time.Sleep(1 * time.Second)
			if liveHost != "dev" {
				// 连接弹幕
				logger.Infof("connecting %s: %s", livePlatform, liveHost)
				deferFunc := func() {
					err := recover()
					if err != nil {
						dialog.ShowInformation("提示", "连接异常断开，请重新连接！", window)
						btnConSetDefault(btnCon, liveHostOption, rLivePlatform, tLiveHostNo)
						btnCon.Refresh()
						// 删除下拉框对应直播间
						if index := utils.Index(liveHost, liveHosts); index != -1 {
							slice, _ := utils.DeleteSlice(liveHosts, index)
							options := slice.([]string)
							liveHostOption.SetOptions(options)
							liveHostOption.Refresh()
							application.Preferences().SetString(PLiveHosts, strings.Join(options, "|"))
						}
						debug.PrintStack()
					}
				}
				switch livePlatform {
				case "Douyu":
					spiderConfig := &DMconfig.DMconfig{
						Rid:            liveHost,
						LoginMsg:       "type@=loginreq/room_id@=%s/dfl@=sn@A=105@Sss@A=1/username@=%s/uid@=%s/ver@=20190610/aver@=218101901/ct@=0/",
						LoginJoinGroup: "type@=joingroup/rid@=%s/gid@=-9999/",
						Url:            "wss://danmuproxy.douyu.com:8506/",
					}
					out := make(chan client.Item, 30)
					webSocketClient = client.DyBarrageWebSocketClient{
						ItemIn: out,
						Config: spiderConfig,
						MsgBreakers: DYtype.CodeBreakershandler{
							IsLive: false,
						},
					}
					go func() {
						defer deferFunc()
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
								handler.FilterMap(msg, livePlatform)
							default:
								logger.ShowJson("[%s]not handle msg: %s", msg.Type, msg)
							}
						}
					}()
				case "Bilibili":
					go func() {
						defer deferFunc()
						if biliClient == nil {
							biliClient = gobilibili.NewBiliBiliClient()
							biliClient.RegHandleFunc(gobilibili.CmdDanmuMsg, func(c *gobilibili.Context) bool {
								info := c.GetDanmuInfo()
								logger.ShowJson("[%v]handle msg: %s", rLivePlatform.Disabled(), info)
								handler.FilterMap(client.Item{
									Rid:   liveHost,
									Uid:   strconv.Itoa(info.UID),
									Txt:   info.Text,
									Nn:    info.Uname,
									Level: strconv.Itoa(info.Level),
									Tst:   false,
								}, livePlatform)

								return false
							})
						}
						liveHostInt, _ := strconv.Atoi(liveHost)
						if err := biliClient.ConnectServer(liveHostInt); err != nil {
							panic(err)
						}
					}()
				}
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
							Tst:     true,
							Payload: nil,
						}, "")
					}
					logger.Debugf("mock %d map", count)
				}()
			}
			if connectSuc {
				// 连接成功
				btnCon.Importance = widget.HighImportance
				btnCon.Text = "成功"
				btnCon.Icon = theme.ConfirmIcon()
				liveHostOption.Hide()
				tLiveHostNo.Show()
				tLiveHostNo.Text = liveHost
				tLiveHostNo.Refresh()
				// 更新下拉框
				if liveHost != "dev" && !utils.In(liveHosts, liveHost) {
					liveHosts = append(liveHosts, liveHost)
					liveHostOption.SetOptions(liveHosts)
					liveHostOption.Refresh()
					application.Preferences().SetString(PLiveHosts, strings.Join(liveHosts, "|"))
				}
				// 保存默认值
				application.Preferences().SetString(PDefaultLiveHostNo, liveHost)
			} else {
				// 连接失败
				btnCon.Importance = widget.DangerImportance
				btnCon.Text = "错误"
				btnCon.Icon = theme.CancelIcon()
			}
		case widget.HighImportance:
			// 连接中 已连接，再次点击断开连接，状态设置成未连接
			btnConSetDefault(btnCon, liveHostOption, rLivePlatform, tLiveHostNo)
			if liveHost != "dev" {
				switch livePlatform {
				case "Douyu":
					go webSocketClient.Stop()
				case "Bilibili":
					biliClient.Stop()
					logger.Infof("disconnect %s", liveHost)
				}
			}
		}
		btnCon.Refresh()
	})
	btnCon.Importance = widget.MediumImportance
	cBtnCon := container.NewMax(btnCon)
	cBtnCon.Resize(fyne.NewSize(65, 36))
	cBtnCon.Move(fyne.NewPos(cLogo.Size().Width+config.ToolbarPaddingLeft+65+config.Padding*2+offsetX, config.ToolbarPaddingTop+offsetY))
	setting.BtnCon = btnCon
	elements = append(elements, cBtnCon)

	// 版本信息
	logger.Info("Initializing version info ...")
	versionText := canvas.NewText(fmt.Sprintf("v%s", version), config.VersionColor)
	versionText.TextSize = 14
	versionText.Alignment = fyne.TextAlignTrailing
	versionText.Move(fyne.NewPos(config.AppSize.Width-10, 0))
	elements = append(elements, versionText)

	// 搜索相关
	logger.Info("Initializing search stuff ...")
	searchEles := generateSearchContainer()
	elements = append(elements, searchEles...)

	// 打开数据文件夹
	logger.Info("Initializing open folder button ...")
	btOpenDataFolder := widget.NewButtonWithIcon("打开数据文件夹", theme.FolderOpenIcon(), func() {
		if err := open.Run(config.UserConfigDir); err != nil {
			dialog.ShowInformation("提示", fmt.Sprintf("打开[%s]失败，错误：%v", config.UserConfigDir, err), window)
		}
	})
	btOpenDataFolder.Resize(fyne.NewSize(150, 35))
	btOpenDataFolder.Move(fyne.NewPos(-10, config.AppSize.Height-btOpenDataFolder.Size().Height-6))
	elements = append(elements, btOpenDataFolder)

	// 说明，从远程获取
	logger.Info("Initializing remark ...")
	remarkText := widget.NewRichTextFromMarkdown(config.RemarkText)
	remarkText.Wrapping = fyne.TextWrapBreak
	remarkText.Resize(fyne.NewSize(540, 100))
	remarkText.Move(fyne.NewPos(560, 0))
	elements = append(elements, remarkText)
	go func() {
		for {
			remark := getRemoteRemark("https://gist.githubusercontent.com/lowking/b45a90d64c15ddde25d28fc11e15a464/raw/remark")
			if remark != "" {
				remarkText.ParseMarkdown(remark)
			}
			time.Sleep(30 * 60 * time.Second)
		}
	}()

	gridLayout := container.NewWithoutLayout(elements...)

	window.SetContent(gridLayout)
	window.Resize(config.AppSize)
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.ShowAndRun()
}

func generateSearchContainer() []fyne.CanvasObject {
	height := float32(37)
	y := float32(71)
	keyWordEntry := searchentry.NewSearchEntry("多个关键词用空格隔开")
	keyWordEntry.Wrapping = fyne.TextTruncate
	keyWordEntry.Resize(fyne.NewSize(250, height))
	keyWordEntry.Move(fyne.NewPos(config.ToolbarPaddingLeft, y))
	idx := 4
	keyWordEntry.OnSubmitted = func(_ string) {
		go doSearch(idx, keyWordEntry)
	}
	keyWordEntry.OnCursorChanged = func() {
		if keyWordEntry.CursorColumn == keyWordEntry.PreviousCursorColumn {
			return
		}
		if setting.SelectedCell {
			setting.SelectedCell = false
		}
		keyWordEntry.PreviousCursorColumn = keyWordEntry.CursorColumn
	}
	searchBtn := widget.NewButtonWithIcon("搜索", theme.SearchIcon(), func() {
		go doSearch(idx, keyWordEntry)
	})
	searchBtn.Resize(fyne.NewSize(90, height))
	searchBtn.Move(fyne.NewPos(config.ToolbarPaddingLeft+keyWordEntry.Size().Width+config.Padding, y))

	setting.OtherEntry["keyWordEntry"] = keyWordEntry
	return []fyne.CanvasObject{keyWordEntry, searchBtn}
}

func doSearch(idx int, keyWordEntry *searchentry.SearchEntry) {
	if !searchLock.TryLockWithTimeout(100 * time.Millisecond) {
		return
	}
	defer searchLock.Unlock()
	handler.RefreshMapList(setting, window, tabs, idx, keyWordEntry.Text, handler.WhereMap[idx], handler.OrderMap[idx], false, false)
}

func getRemoteRemark(url string) string {
	ret := ""
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ret
	}
	c := http.Client{}
	response, err := c.Do(req)
	if err != nil {
		if strings.Index(url, "https://ghproxy.com/") == 0 {
			return ret
		}
		url = fmt.Sprintf("https://ghproxy.com/%s", url)
		return getRemoteRemark(url)
	}
	if response.StatusCode != 200 {
		if strings.Index(url, "https://ghproxy.com/") == 0 {
			return ret
		}
		url = fmt.Sprintf("https://ghproxy.com/%s", url)
		return getRemoteRemark(url)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ret
	}
	return string(body)
}

func btnConSetDefault(btnCon *widget.Button, liveHostOption *widget.SelectEntry, livePlatform *widget.RadioGroup, liveHostNo *canvas.Text) {
	btnCon.Importance = widget.MediumImportance
	btnCon.Text = "连接"
	btnCon.Icon = theme.NavigateNextIcon()
	liveHostOption.Show()
	liveHostNo.Hide()
	livePlatform.Enable()
}

func flashEle(times int, elements ...*canvas.Image) {
	for i := 0; i < times; i++ {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		for _, element := range elements {
			if element.Visible() {
				element.Hide()
			} else {
				element.Show()
			}
			// 不要一直刷新，内存会暴涨🤣
			// element.Refresh()
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
			handler.RefreshMapList(setting, window, tabs, 0, "", handler.WhereMap[0], handler.OrderMap[0], true, false)
		}
		handler.RefreshMapList(setting, window, tabs, 1, "", handler.WhereMap[1], handler.OrderMap[1], true, false)
		handler.RefreshMapList(setting, window, tabs, 2, "", handler.WhereMap[2], handler.OrderMap[2], true, false)

		if setting.AutoGetFgPid {
			go setting.BtnGetFgPid.OnTapped()
		}
		if setting.AutoConnect {
			go setting.BtnCon.OnTapped()
		}
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	refreshLock := lock.NewCASMutex()
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
		if tabs.SelectedIndex() >= 3 {
			return
		}
		if !refreshLock.TryLockWithTimeout(100 * time.Millisecond) {
			return
		}
		defer refreshLock.Unlock()
		// 每次聚焦窗口刷新列表
		refreshList()

		// 聚焦窗口开启每秒自动刷新
		enableQueryLoop = true
		go func() {
			for {
				time.Sleep(time.Second)
				if !enableQueryLoop {
					break
				}
				refreshList()
			}
		}()
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
		enableQueryLoop = false
	})
}

func refreshList() {
	if setting.SelectedCell {
		return
	}
	switch tabs.SelectedIndex() {
	case 0:
		go handler.RefreshMapList(setting, window, tabs, 0, "", handler.WhereMap[0], handler.OrderMap[0], false, false)
	case 1:
		go handler.RefreshMapList(setting, window, tabs, 1, "", handler.WhereMap[1], handler.OrderMap[1], false, false)
	case 2:
		go handler.RefreshMapList(setting, window, tabs, 2, "", handler.WhereMap[2], handler.OrderMap[2], false, false)
	}
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		resetLiveHostOptions := fyne.NewMenuItem("还原直播间下拉框", func() {
			a.Preferences().SetString(PLiveHosts, "")
			dialog.ShowInformation("提示", "已还原，重启生效", window)
		})
		menu := fyne.NewMenu("菜单", resetLiveHostOptions)
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
