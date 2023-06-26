package handler

import (
	"FallGuys66/db"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"github.com/PaulWaldo/fyne-headertable/headertable"
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
					return "---"
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
					return "♥️"
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
			WidthPercent: 390,
			Converter: func(i interface{}) string {
				t := i.(string)
				rowLen := 40
				ret := warpStr(t, rowLen)
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

func warpStr(str string, rowLen int) string {
	runes := []rune(str)
	if len(runes) > rowLen {
		times := len(runes) / rowLen
		ts := ""
		for i := 1; i <= times; i++ {
			ts = fmt.Sprintf("%s%s\n", ts, string(runes[:rowLen]))
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

func RefreshMapList(tabs *container.AppTabs, idx int, where string, order string) {
	// 查询数据库获取最新列表
	key := fmt.Sprintf("map%d", idx)
	listMap := db.ListMap(1, 18, where, order)
	listLength := len(listMap)
	tListHeader := listHeader
	if listLength > 0 {
		if len(bindingsMap[key]) > 0 && bindingsMap[key][0] != nil {
			if val, err := bindingsMap[key][0].GetValue("MapId"); err == nil {
				if val == cache[key] {
					return
				}
			}
		}
		bindingsMap[key] = make([]binding.Struct, listLength)
		for i := 0; i < listLength; i++ {
			bindingsMap[key][i] = binding.BindStruct(&listMap[i])
		}
		tListHeader.Bindings = bindingsMap[key]
		ht := headertable.NewHeaderTable(&tListHeader)
		tabs.Items[idx].Content = container.NewMax(ht)
		cache[key] = listMap[0].MapId
	}
}
