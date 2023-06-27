package handler

import (
	"FallGuys66/db"
	"FallGuys66/utils"
	"FallGuys66/widgets/headertable"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"time"
)

var listHeader = headertable.TableOpts{
	RefWidth: "reference width",
	ColAttrs: []headertable.ColAttr{
		{
			Name:   "MapId",
			Header: "地图代码",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
			},
			WidthPercent: 120,
		},
		{
			Name:   "Nn",
			Header: "用户名",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
			},
			WidthPercent: 120,
			Converter: func(i interface{}) string {
				return i.(string)
			},
		},
		{
			Name:   "Uid",
			Header: "用户ID",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignTrailing,
			},
			WidthPercent: 90,
		},
		{
			Name:   "Rid",
			Header: "投稿直播间",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignTrailing,
			},
			WidthPercent: 70,
		},
		{
			Name:   "Level",
			Header: "等级",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignTrailing,
			},
			WidthPercent: 40,
		},
		{
			Name:   "State",
			Header: "状态",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
			},
			WidthPercent: 40,
			Converter: func(i interface{}) string {
				t := i.(string)
				if t == "0" {
					return ""
				} else {
					// return "---"
					return "🎮"
				}
			},
		},
		{
			Name:   "Star",
			Header: "收藏",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
			},
			WidthPercent: 40,
			Converter: func(i interface{}) string {
				t := i.(string)
				if t == "0" {
					return ""
				} else {
					return "★"
				}
			},
		},
		{
			Name:   "Txt",
			Header: "原始弹幕",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignLeading,
			},
			WidthPercent: 200,
			Converter: func(i interface{}) string {
				t := i.(string)
				rowLen := 20
				ret := warpStr(t, rowLen, true)
				return ret
			},
		},
		{
			Name:   "Created",
			Header: "投稿时间",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignLeading,
			},
			WidthPercent: 135,
			Converter: func(i interface{}) string {
				t := i.(time.Time)
				return t.Format("2006-01-02 15:04:05")
			},
		},
	},
}

func warpStr(str string, rowLen int, isSingleLine bool) string {
	runes := []rune(str)
	if len(runes) > rowLen {
		times := len(runes) / rowLen
		ts := ""
		for i := 1; i <= times; i++ {
			ts = fmt.Sprintf("%s%s\n", ts, string(runes[:rowLen]))
			if isSingleLine {
				return fmt.Sprintf("%s ...", ts[:len(ts)-1])
			}
			runes = runes[rowLen:]
			if i == times {
				ts = fmt.Sprintf("%s%s\n", ts, string(runes))
			}
		}
		return ts[:len(ts)-1]
	}
	return str
}

var cache = make(map[string]string)
var bindingsMap = make(map[string][]binding.Struct, 1)

// 点击时临时存储单元格内容
var cellTempString string
var bsCellTempString = binding.BindString(&cellTempString)
var mapIdTempString string
var bsMapIdTempString = binding.BindString(&mapIdTempString)

func RefreshMapList(driver fyne.Driver, window fyne.Window, tabs *container.AppTabs, idx int, where string, order string) {
	// 查询数据库获取最新列表
	key := fmt.Sprintf("map%d", idx)
	listMap := db.ListMap(1, 18, where, order)
	listLength := len(listMap)
	tListHeader := listHeader
	if listLength > 0 {
		// 已经去掉了轮训刷新，无需校验缓存
		// if len(bindingsMap[key]) > 0 && bindingsMap[key][0] != nil {
		// 	if val, err := bindingsMap[key][0].GetValue("MapId"); err == nil {
		// 		if val == cache[key] {
		// 			return
		// 		}
		// 	}
		// }
		bindingsMap[key] = make([]binding.Struct, listLength)
		for i := 0; i < listLength; i++ {
			bindingsMap[key][i] = binding.BindStruct(&listMap[i])
		}
		tListHeader.Bindings = bindingsMap[key]
		ht := headertable.NewHeaderTable(&tListHeader)
		tabs.Items[idx].Content = container.NewMax(ht)
		cache[key] = listMap[0].MapId
		tableMenu := fyne.NewMenu(
			"file",
			fyne.NewMenuItem("", func() {
				if s, err := bsCellTempString.Get(); err == nil {
					window.Clipboard().SetContent(s)
				}
			}),
			fyne.NewMenuItem("游玩", func() {
				if s, err := bsMapIdTempString.Get(); err == nil {
					window.Clipboard().SetContent(s)
					go utils.FillMapId(s)
				}
			}),
		)
		ht.Data.OnSelected = func(id widget.TableCellID) {
			row := bindingsMap[key][id.Row]
			colKey := tListHeader.ColAttrs[id.Col].Name
			if value, err := row.GetValue(colKey); err == nil {
				tablePos := driver.AbsolutePositionForObject(ht.Data)

				// 每次点击记录地图Id
				if mid, err := row.GetValue("MapId"); err == nil {
					_ = bsMapIdTempString.Set(mid.(string))
				}
				// 处理单元格字符
				valueString := ""
				if s, ok := value.(string); ok {
					valueString = s
				} else if s, ok := value.(time.Time); ok {
					valueString = s.Format("2006-01-02 15:04:05")
				}
				tableMenu.Items[0].Label = valueString
				_ = bsCellTempString.Set(valueString)

				xx, yy := getCellPos(tablePos, id.Col, id.Row, tListHeader, 36.1)
				widget.NewPopUpMenu(tableMenu, window.Canvas()).ShowAtPosition(fyne.NewPos(xx, yy))
			}
		}
	}
}

func getCellPos(base fyne.Position, x int, y int, header headertable.TableOpts, cellHeight float32) (float32, float32) {
	xx := base.X
	yy := base.Y
	for i := 0; i < x; i++ {
		xx += float32(header.ColAttrs[i].WidthPercent) + 15
	}

	return xx + 15, yy + (float32(y)+0.7)*cellHeight
}
