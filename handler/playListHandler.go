package handler

import (
	"FallGuys66/config"
	"FallGuys66/db"
	"FallGuys66/db/model"
	"FallGuys66/utils"
	"FallGuys66/widgets/headertable"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/wesovilabs/koazee"
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

var ht *headertable.HeaderTable

const pageSize = 18

var listMap [pageSize]model.MapInfo

func RefreshMapList(driver fyne.Driver, window fyne.Window, tabs *container.AppTabs, idx int, where string, order string, recreate bool) {
	// 查询数据库获取最新列表
	key := fmt.Sprintf("map%d", idx)
	recreateKey := fmt.Sprintf("map%dRecreate", idx)
	// 外部传入recreate或者缓存中是true，则重新new
	tListMap := db.ListMap(1, pageSize, where, order)
	listLength := len(tListMap)
	recreate = recreate || cache[recreateKey] == "true" || listLength < len(listMap)
	tListHeader := listHeader
	if listLength > 0 {
		cache[recreateKey] = "false"
		for i := 0; i < listLength; i++ {
			listMap[i] = tListMap[i]
		}
		if !recreate {
			ht.Refresh()
			return
		}
		bindingsMap[key] = make([]binding.Struct, listLength)
		for i := 0; i < listLength; i++ {
			bindingsMap[key][i] = binding.BindStruct(&listMap[i])
		}
		tListHeader.Bindings = bindingsMap[key]
		ht = headertable.NewHeaderTable(&tListHeader)
		cache[key] = listMap[0].MapId
		firstItemAction := func() {
			if s, err := bsCellTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
			}
		}
		firstMenuItem := fyne.NewMenuItem("", firstItemAction)
		playMenuItem := fyne.NewMenuItem("游玩", func() {
			if s, err := bsMapIdTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
				go utils.FillMapId(s)
				go func() {
					db.UpdateMap(s, `state="1"`, `and mapId=? and state="0"`)
					RefreshMapList(driver, window, tabs, idx, where, order, false)
				}()
			}
		})
		starMenuItem := fyne.NewMenuItem("收藏", func() {
			if s, err := bsMapIdTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
				go utils.FillMapId(s)
				go func() {
					db.UpdateMap(s, `star="1"`, `and mapId=? and star="0"`)
					RefreshMapList(driver, window, tabs, idx, where, order, false)
				}()
			}
		})
		unStarMenuItem := fyne.NewMenuItem("取消收藏", func() {
			if s, err := bsMapIdTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
				go utils.FillMapId(s)
				go func() {
					db.UpdateMap(s, `star="0"`, `and mapId=? and star="1"`)
					RefreshMapList(driver, window, tabs, idx, where, order, false)
				}()
			}
		})
		tableMenu := fyne.NewMenu("Actions", firstMenuItem, playMenuItem, starMenuItem, unStarMenuItem)
		ht.Data.OnSelected = func(id widget.TableCellID) {
			row := bindingsMap[key][id.Row]
			colKey := tListHeader.ColAttrs[id.Col].Name
			if value, err := row.GetValue(colKey); err == nil {
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

				star, _ := row.GetValue("Star")
				// 根据收藏状态，调整菜单显示
				deleteTargetItem := starMenuItem
				addTargetItem := unStarMenuItem
				if star == "0" {
					deleteTargetItem = unStarMenuItem
					addTargetItem = starMenuItem
				}
				index := -1
				for i, item := range tableMenu.Items {
					if item.Label == deleteTargetItem.Label {
						index = i
						break
					}
				}
				if index > -1 {
					if slice, err := utils.DeleteSlice(tableMenu.Items, index); err == nil {
						tableMenu.Items = slice.([]*fyne.MenuItem)
					}
				}
				index = -1
				for i, item := range tableMenu.Items {
					if item.Label == addTargetItem.Label {
						index = i
						break
					}
				}
				if index == -1 {
					tableMenu.Items = append(tableMenu.Items, addTargetItem)
				}

				// 处理点击状态和收藏，移除第一个菜单
				shouldBe := 3
				switch colKey {
				case "State", "Star", "Level":
					if len(tableMenu.Items) == shouldBe {
						_, stream := koazee.StreamOf(tableMenu.Items).Pop()
						tableMenu.Items = stream.Out().Val().([]*fyne.MenuItem)
					}
				default:
					if len(tableMenu.Items) < shouldBe {
						newTableMenuItems := []*fyne.MenuItem{
							firstMenuItem,
						}
						tableMenu.Items = append(newTableMenuItems, tableMenu.Items...)
					}
					tableMenu.Items[0].Label = valueString
				}

				_ = bsCellTempString.Set(valueString)

				xx, yy := getCellPos(fyne.NewPos(0, 220), id.Col, id.Row, tListHeader, 36.1)
				widget.NewPopUpMenu(tableMenu, window.Canvas()).ShowAtPosition(fyne.NewPos(xx, yy))
			}
		}
		tabs.Items[idx].Content = container.NewMax(ht)
	} else {
		cache[recreateKey] = "true"
		tabs.Items[idx].Content = utils.MakeEmptyList(config.AccentColor)
	}
}

func RefreshMapList1(driver fyne.Driver, window fyne.Window, tabs *container.AppTabs, idx int, where string, order string) {
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
		ht = headertable.NewHeaderTable(&tListHeader)
		cache[key] = listMap[0].MapId
		firstItemAction := func() {
			if s, err := bsCellTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
			}
		}
		firstMenuItem := fyne.NewMenuItem("", firstItemAction)
		playMenuItem := fyne.NewMenuItem("游玩", func() {
			if s, err := bsMapIdTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
				go utils.FillMapId(s)
				go func() {
					db.UpdateMap(s, `state="1"`, `and mapId=? and state="0"`)
					RefreshMapList(driver, window, tabs, idx, where, order, false)
				}()
			}
		})
		starMenuItem := fyne.NewMenuItem("收藏", func() {
			if s, err := bsMapIdTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
				go utils.FillMapId(s)
				go func() {
					db.UpdateMap(s, `star="1"`, `and mapId=? and star="0"`)
					RefreshMapList(driver, window, tabs, idx, where, order, false)
				}()
			}
		})
		unStarMenuItem := fyne.NewMenuItem("取消收藏", func() {
			if s, err := bsMapIdTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
				go utils.FillMapId(s)
				go func() {
					db.UpdateMap(s, `star="0"`, `and mapId=? and star="1"`)
					RefreshMapList(driver, window, tabs, idx, where, order, false)
				}()
			}
		})
		tableMenu := fyne.NewMenu("Actions", firstMenuItem, playMenuItem, starMenuItem, unStarMenuItem)
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

				star, _ := row.GetValue("Star")
				// 根据收藏状态，调整菜单显示
				deleteTargetItem := starMenuItem
				addTargetItem := unStarMenuItem
				if star == "0" {
					deleteTargetItem = unStarMenuItem
					addTargetItem = starMenuItem
				}
				index := -1
				for i, item := range tableMenu.Items {
					if item.Label == deleteTargetItem.Label {
						index = i
						break
					}
				}
				if index > -1 {
					if slice, err := utils.DeleteSlice(tableMenu.Items, index); err == nil {
						tableMenu.Items = slice.([]*fyne.MenuItem)
					}
				}
				index = -1
				for i, item := range tableMenu.Items {
					if item.Label == addTargetItem.Label {
						index = i
						break
					}
				}
				if index == -1 {
					tableMenu.Items = append(tableMenu.Items, addTargetItem)
				}

				// 处理点击状态和收藏，移除第一个菜单
				shouldBe := 3
				switch colKey {
				case "State", "Star", "Level":
					if len(tableMenu.Items) == shouldBe {
						_, stream := koazee.StreamOf(tableMenu.Items).Pop()
						tableMenu.Items = stream.Out().Val().([]*fyne.MenuItem)
					}
				default:
					if len(tableMenu.Items) < shouldBe {
						newTableMenuItems := []*fyne.MenuItem{
							firstMenuItem,
						}
						tableMenu.Items = append(newTableMenuItems, tableMenu.Items...)
					}
					tableMenu.Items[0].Label = valueString
				}

				_ = bsCellTempString.Set(valueString)

				xx, yy := getCellPos(tablePos, id.Col, id.Row, tListHeader, 36.1)
				widget.NewPopUpMenu(tableMenu, window.Canvas()).ShowAtPosition(fyne.NewPos(xx, yy))
			}
		}
		tabs.Items[idx].Content = container.NewMax(ht)
	} else {
		tabs.Items[idx].Content = utils.MakeEmptyList(config.AccentColor)
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
