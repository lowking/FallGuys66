package utils

import (
	"FallGuys66/common/cbm"
	"FallGuys66/live/douyu/lib/logger"
	"FallGuys66/settings"
	"github.com/go-vgo/robotgo"
	"runtime/debug"
	"strconv"
	"strings"
)

func init() {
	cbm.RegisterCallBack("fg FillMapIdForPlayNext", FillMapIdForPlayNext)
	cbm.RegisterCallBack("fg FillMapIdForPlayNextOnSelectShow", FillMapIdForPlayNextOnSelectShow)
}

func FillMapId(id string, settings *settings.Settings) {
	if !settings.AutoFillMapId || settings.FgPid == "" || settings.PosSelectShow == nil || settings.PosEnterMapId == nil || settings.PosCodeEntry == nil || settings.PosConfirmBtn == nil || *settings.PosSelectShow == "" || *settings.PosEnterMapId == "" || *settings.PosCodeEntry == "" || *settings.PosConfirmBtn == "" {
		return
	}
	pid, err := strconv.ParseInt(settings.FgPid, 10, 32)
	if err != nil {
		return
	}
	posSelectShow := strings.Split(*settings.PosSelectShow, ",")
	posEnterMapId := strings.Split(*settings.PosEnterMapId, ",")
	posCodeEntry := strings.Split(*settings.PosCodeEntry, ",")
	posConfirmBtn := strings.Split(*settings.PosConfirmBtn, ",")
	if len(posSelectShow) != 2 || len(posEnterMapId) != 2 || len(posCodeEntry) != 2 || len(posConfirmBtn) != 2 {
		return
	}
	x1, err := strconv.Atoi(posSelectShow[0])
	if err != nil {
		return
	}
	x2, err := strconv.Atoi(posEnterMapId[0])
	if err != nil {
		return
	}
	x3, err := strconv.Atoi(posCodeEntry[0])
	if err != nil {
		return
	}
	x4, err := strconv.Atoi(posConfirmBtn[0])
	if err != nil {
		return
	}
	y1, err := strconv.Atoi(posSelectShow[1])
	if err != nil {
		return
	}
	y2, err := strconv.Atoi(posEnterMapId[1])
	if err != nil {
		return
	}
	y3, err := strconv.Atoi(posCodeEntry[1])
	if err != nil {
		return
	}
	y4, err := strconv.Atoi(posConfirmBtn[1])
	if err != nil {
		return
	}
	err = robotgo.ActivePID(int32(pid))
	robotgo.MilliSleep(500)
	if err != nil {
		return
	}
	defer func() {
		err := recover()
		if err != nil {
			logger.Errorf("执行自动填写id错误：%v", err)
			debug.PrintStack()
		}
	}()
	// 点击选择节目
	robotgo.MoveClick(x1, y1, "left")
	robotgo.MilliSleep(500)
	// 点击输入代码
	robotgo.MoveClick(x2, y2, "left")
	robotgo.MilliSleep(300)
	// 点击代码文本框
	robotgo.MoveClick(x3, y3, "left")
	robotgo.MilliSleep(100)
	// switch runtime.GOOS {
	// case "darwin":
	// 	robotgo.KeyTap("a", "command")
	// default:
	// 	robotgo.KeyTap("a", "ctrl")
	// }
	// robotgo.DragSmooth(x3+300, y3)
	// 输入地图代码
	robotgo.TypeStr(id)
	robotgo.MilliSleep(300)
	robotgo.MoveClick(x4, y4, "left")
}

func FillMapIdForPlayNext(id string, settings *settings.Settings) {
	if settings.PosEnterMapId == nil || settings.PosCodeEntry == nil || settings.PosConfirmBtn == nil || *settings.PosEnterMapId == "" || *settings.PosCodeEntry == "" || *settings.PosConfirmBtn == "" {
		return
	}

	posEnterMapId := strings.Split(*settings.PosEnterMapId, ",")
	posCodeEntry := strings.Split(*settings.PosCodeEntry, ",")
	posConfirmBtn := strings.Split(*settings.PosConfirmBtn, ",")
	if len(posEnterMapId) != 2 || len(posCodeEntry) != 2 || len(posConfirmBtn) != 2 {
		return
	}

	x2, err := strconv.Atoi(posEnterMapId[0])
	if err != nil {
		return
	}
	x3, err := strconv.Atoi(posCodeEntry[0])
	if err != nil {
		return
	}
	x4, err := strconv.Atoi(posConfirmBtn[0])
	if err != nil {
		return
	}
	y2, err := strconv.Atoi(posEnterMapId[1])
	if err != nil {
		return
	}
	y3, err := strconv.Atoi(posCodeEntry[1])
	if err != nil {
		return
	}
	y4, err := strconv.Atoi(posConfirmBtn[1])
	if err != nil {
		return
	}

	defer func() {
		err := recover()
		if err != nil {
			logger.Errorf("执行自动填写id(PlayNext)错误：%v", err)
			debug.PrintStack()
		}
	}()
	// 点击输入代码
	robotgo.MoveClick(x2, y2, "left")
	robotgo.MilliSleep(300)
	// 点击代码文本框
	robotgo.MoveClick(x3, y3, "left")
	robotgo.MilliSleep(100)
	// switch runtime.GOOS {
	// case "darwin":
	// 	robotgo.KeyTap("a", "command")
	// default:
	// 	robotgo.KeyTap("a", "ctrl")
	// }
	// robotgo.DragSmooth(x3+300, y3)
	// 输入地图代码
	robotgo.TypeStr(id)
	robotgo.MilliSleep(300)
	robotgo.MoveClick(x4, y4, "left")
}

func FillMapIdForPlayNextOnSelectShow(id string, settings *settings.Settings) {
	if settings.PosSelectShow == nil || settings.PosEnterMapId == nil || settings.PosCodeEntry == nil || settings.PosConfirmBtn == nil || *settings.PosSelectShow == "" || *settings.PosEnterMapId == "" || *settings.PosCodeEntry == "" || *settings.PosConfirmBtn == "" {
		return
	}

	posSelectShow := strings.Split(*settings.PosSelectShow, ",")
	posEnterMapId := strings.Split(*settings.PosEnterMapId, ",")
	posCodeEntry := strings.Split(*settings.PosCodeEntry, ",")
	posConfirmBtn := strings.Split(*settings.PosConfirmBtn, ",")
	if len(posSelectShow) != 2 || len(posEnterMapId) != 2 || len(posCodeEntry) != 2 || len(posConfirmBtn) != 2 {
		return
	}
	x1, err := strconv.Atoi(posSelectShow[0])
	if err != nil {
		return
	}
	x2, err := strconv.Atoi(posEnterMapId[0])
	if err != nil {
		return
	}
	x3, err := strconv.Atoi(posCodeEntry[0])
	if err != nil {
		return
	}
	x4, err := strconv.Atoi(posConfirmBtn[0])
	if err != nil {
		return
	}
	y1, err := strconv.Atoi(posSelectShow[1])
	if err != nil {
		return
	}
	y2, err := strconv.Atoi(posEnterMapId[1])
	if err != nil {
		return
	}
	y3, err := strconv.Atoi(posCodeEntry[1])
	if err != nil {
		return
	}
	y4, err := strconv.Atoi(posConfirmBtn[1])
	if err != nil {
		return
	}

	defer func() {
		err := recover()
		if err != nil {
			logger.Errorf("执行自动填写id(OnSelectShow)错误：%v", err)
			debug.PrintStack()
		}
	}()
	// 点击选择节目
	robotgo.MoveClick(x1, y1, "left")
	robotgo.MilliSleep(500)
	// 点击输入代码
	robotgo.MoveClick(x2, y2, "left")
	robotgo.MilliSleep(300)
	// 点击代码文本框
	robotgo.MoveClick(x3, y3, "left")
	robotgo.MilliSleep(100)
	// switch runtime.GOOS {
	// case "darwin":
	// 	robotgo.KeyTap("a", "command")
	// default:
	// 	robotgo.KeyTap("a", "ctrl")
	// }
	// robotgo.DragSmooth(x3+300, y3)
	// 输入地图代码
	robotgo.TypeStr(id)
	robotgo.MilliSleep(300)
	robotgo.MoveClick(x4, y4, "left")
}
