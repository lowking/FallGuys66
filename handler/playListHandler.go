package handler

import (
	"FallGuys66/common/cbm"
	"FallGuys66/config"
	"FallGuys66/db"
	"FallGuys66/db/model"
	"FallGuys66/live/douyu/lib/logger"
	"FallGuys66/settings"
	"FallGuys66/utils"
	"FallGuys66/widgets/headertable"
	"FallGuys66/widgets/pager"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/wesovilabs/koazee"
	"golang.org/x/text/width"
	"math"
	"strings"
	"time"
)

func init() {
	cbm.RegisterCallBack("sv isCleanMapId", func(b bool) { isCleanMapId = b })
}

var isCleanMapId = false

// "｟", "｠", "《", "》", "（", "）", "＜", "＞", "［", "］", "｢", "｣", "〈", "〉", "「", "」", "『", "』", "【", "】", "〔", "〕", "〖", "〗", "〘", "〙", "〚", "〛", "‘", "’", "‛", "“", "”", "„", "‟", "\"", "(", ")", "<", ">", "[", "]",
var punctuation = []string{"＂", "＃", "＄", "％", "＆", "＇", "＊", "＋", "，", "－", "／", "：", "；", "＝", "＠", "＼", "＾", "＿", "｀", "｛", "｜", "｝", "～", "､", "　", "、", "〃", "〜", "〝", "〞", "〟", "〰", "〾", "〿", "–", "—", "…", "‧", "﹏", "﹑", "﹔", "·", "．", "！", "？", "｡", "。", "!", "#", "$", "%", "&", "'", "*", "+", ",", "-", ".", "/", ":", ";", "=", "?", "@", "\\", "^", "_", "`", "{", "|", "}", "~"}
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
			WidthPercent: 195,
			Converter: func(i interface{}, row binding.Struct) string {
				if level, err := row.GetValue("Level"); err == nil && level != "" {
					return fmt.Sprintf("%s (Lv. %v)", i.(string), level)
				}
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
			WidthPercent: 100,
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
			WidthPercent: 121,
			Converter: func(i interface{}, row binding.Struct) string {
				t := i.(string)
				if livePlatform, err := row.GetValue("LivePlatform"); err == nil && t != "" {
					return fmt.Sprintf("%s (%v)", t, livePlatform)
				}
				return t
			},
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
			Converter: func(i interface{}, row binding.Struct) string {
				t := i.(string)
				if t == "1" {
					return "🎮"
				} else {
					return ""
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
			Converter: func(i interface{}, row binding.Struct) string {
				t := i.(string)
				if t == "1" {
					return "★"
				} else {
					return ""
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
				Wrapping:  fyne.TextTruncate,
			},
			WidthPercent: 200,
			Converter:    dmConvertor,
		},
		{
			Name:   "Created",
			Header: "投稿时间",
			HeaderStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignCenter,
				TextStyle: fyne.TextStyle{Bold: true},
			},
			DataStyle: headertable.CellStyle{
				Alignment: fyne.TextAlignTrailing,
			},
			WidthPercent: 135,
			Converter: func(i interface{}, row binding.Struct) string {
				t := i.(time.Time)
				if t.IsZero() {
					return ""
				}
				return t.Format("2006-01-02 15:04:05")
			},
		},
	},
}

func dmConvertor(i interface{}, row binding.Struct) string {
	t := i.(string)
	if isCleanMapId {
		for _, mapId := range mapRe.FindAllString(t, -1) {
			t = strings.ReplaceAll(t, mapId, "")
		}
		source := i.(string)
		index := strings.Index(source, t)
		if index >= len(source)-len(t) {
			t = strings.TrimLeftFunc(strings.TrimSpace(t), func(r rune) bool {
				return utils.In(punctuation, string(r))
			})
		} else {
			t = strings.TrimRightFunc(strings.TrimSpace(t), func(r rune) bool {
				return utils.In(punctuation, string(r))
			})
		}
	}
	return strings.TrimSpace(t)
}

func wrapStr(s string, rowLen int, isSingleLine bool) string {
	w := 0
	ts := ""
	for _, r := range s {
		if w >= rowLen {
			if isSingleLine {
				return ts + " ..."
			} else {
				ts = fmt.Sprintf("%s\n", ts)
				w = 0
			}
		}
		switch width.LookupRune(r).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			w += 2
		case width.EastAsianHalfwidth, width.EastAsianNarrow, width.Neutral, width.EastAsianAmbiguous:
			w += 1
		}
		ts = fmt.Sprintf("%s%s", ts, string(r))
	}
	return ts
}

var pSize = 17
var pPageSize = &pSize

const pageSize = 17

var (
	WhereMap = make(map[int]interface{})
	OrderMap = make(map[int]string)
)

var cache = make(map[string]string)
var cacheHt = make(map[string]*headertable.HeaderTable)
var cacheKeyword = ""
var cacheKeywordForFoundCell = ""
var cacheListHeader = make(map[string]headertable.TableOpts)
var cacheCurrentNo = make(map[string]*int)
var cachePager = make(map[string]*pager.Pager)
var bindingsMap = make(map[string][]binding.Struct, 1)

// 点击时临时存储单元格内容
var cellTempString string
var bsCellTempString = binding.BindString(&cellTempString)
var mapInfoTemp model.MapInfo
var bsMapInfoTemp = binding.BindStruct(&mapInfoTemp)

var listMapPlay = [pageSize]model.MapInfo{}
var listMapPlayed = [pageSize]model.MapInfo{}
var listMapStar = [pageSize]model.MapInfo{}
var searchResult = [pageSize]model.MapInfo{}
var listMap *[pageSize]model.MapInfo
var tListMap = make(map[string][]model.MapInfo)
var tBlacklist []model.Blacklist
var whereString = "map_id like ? or nn like ? or uid like ? or rid like ? or txt like ?"
var whereStringForList = []string{"MapId", "Nn", "Uid", "Rid", "Txt"}
var allowTapManual = true
var previousCellX = 0
var previousCellY = 0
var previousTabIndex = 0

func RefreshMapList(
	settings *settings.Settings,
	window fyne.Window,
	tabs *container.AppTabs,
	idx int,
	keyWord string,
	where interface{},
	order string,
	recreate bool,
	fromPager bool,
) {
	// 重置一些参数
	defer func() {
		previousTabIndex = tabs.SelectedIndex()
	}()
	allowTapManual = true
	settings.SelectedCell = false
	key := fmt.Sprintf("map%d", idx)
	if _, ok := cacheCurrentNo[key]; !ok {
		no := 1
		cacheCurrentNo[key] = &no
	}
	if keyWord != "" {
		keyWord = strings.ToLower(keyWord)
		if cacheKeyword != keyWord {
			previousCellX = 0
		}
		// 搜索当前列表，有的画就高亮，否则执行搜索
		if previousCellX != -1 {
			var i, j int
			if cacheKeywordForFoundCell == keyWord && previousTabIndex == tabs.SelectedIndex() {
				i = previousCellX
				j = previousCellY + 1
				if j >= len(listHeader.ColAttrs) {
					i = previousCellX + 1
					j = 0
				}
			}
			findCellKey := fmt.Sprintf("map%d", tabs.SelectedIndex())
			table := cacheHt[findCellKey]
			header := cacheListHeader[findCellKey]
			if tabs.SelectedIndex() == 4 {
				tFindCellKey := fmt.Sprintf("map%d", 0)
				table = cacheHt[tFindCellKey]
				header = cacheListHeader[tFindCellKey]
			}
			// 每行数据开始搜索
			isFound := false
			for ; table != nil && i < len(table.TableOpts.Bindings); i++ {
				bds := table.TableOpts.Bindings[i]
				// 搜索指定列
				for ; j < len(header.ColAttrs); j++ {
					colAttr := header.ColAttrs[j]
					// 定位到指定列
					if !utils.In(whereStringForList, colAttr.Name) {
						continue
					}
					// 获取列的值
					if t, err := bds.GetValue(colAttr.Name); err == nil {
						// 根据列的转换器处理文本
						if colAttr.Converter != nil {
							t = colAttr.Converter(t, bds)
						}
						// 列文本匹配搜索关键字
						if strings.Index(strings.ToLower(t.(string)), keyWord) == -1 {
							continue
						}
						if tabs.SelectedIndex() == 4 {
							tabs.SelectIndex(0)
							time.Sleep(200 * time.Millisecond)
						}
						settings.SelectedCell = true
						allowTapManual = false
						previousCellX = i
						previousCellY = j
						isFound = true
						flashCell(table, i, j)
						allowTapManual = true
						cacheKeywordForFoundCell = keyWord
						cacheKeyword = keyWord
						return
					}
				}
				j = 0
			}
			cacheKeyword = keyWord
			// 未在列表搜索到相关信息就返回，再次搜索才到数据库中查询
			if !isFound && previousCellX != -1 {
				if tabs.SelectedIndex() != 4 {
					previousCellX = -1
					return
				}
			}
		}
		cacheKeyword = keyWord
		tabs.SelectIndex(idx)
	}
	// 查询数据库获取最新列表
	isRefresh := false
	var count int64
	switch idx {
	case 0, 1, 2:
		tListMap[key], count = db.ListMap(*cacheCurrentNo[key], pageSize, where.(*model.MapInfo), order, idx == 0)
		if len(tListMap[key]) == 0 {
			// 如果刷新发现没数据，页码-1再次刷新，直到页码为1
			if *cacheCurrentNo[key] <= 1 {
				break
			}
			*cacheCurrentNo[key]--
			RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], recreate, fromPager)
			return
		}
	case 3:
		tBlacklist, count = db.ListBlacklist(*cacheCurrentNo[key], pageSize, where.(*model.Blacklist), order)
		if len(tBlacklist) == 0 {
			// 如果刷新发现没数据，页码-1再次刷新，直到页码为1
			if *cacheCurrentNo[key] <= 1 {
				break
			}
			*cacheCurrentNo[key]--
			RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], recreate, fromPager)
			return
		}
	case 4:
		if cacheKeyword == "" {
			return
		}
		cacheKeyword = strings.TrimSpace(cacheKeyword)
		if cacheKeyword != "" {
			rWhere := ""
			for _, s := range strings.Split(cacheKeyword, " ") {
				if s == "" {
					continue
				}
				rWhere = fmt.Sprintf(
					`%s or %s`,
					strings.ReplaceAll(whereString, "?", fmt.Sprintf(`"%%%s%%"`, s)),
					rWhere,
				)
			}
			if rWhere != "" {
				rWhere = rWhere[:len(rWhere)-3]
			}
			// 每次搜索重置第一页
			if !fromPager {
				*cacheCurrentNo[key] = 1
				isRefresh = true
			}
			tListMap[key], count = db.SearchMap(*cacheCurrentNo[key], pageSize, rWhere, order)
			previousCellX = 0
			previousCellY = 0
		}
	}
	// 总页数变了就得刷新一次页码
	if cachePager[key] != nil {
		if isRefresh || (cachePager[key] != nil && int(math.Ceil(float64(count)/float64(pageSize))) != *cachePager[key].PageCount) {
			cachePager[key] = pager.NewPager(cacheCurrentNo[key], pPageSize, &count, cachePager[key].OnTapped)
			cBorder := container.NewBorder(nil, cachePager[key], nil, nil, cacheHt[key])
			tabs.Items[idx].Content = cBorder
		} else {
			*cachePager[key].Total = count
			(*cachePager[key].Items)[1].(*widget.Label).SetText(fmt.Sprintf("共 %d 条", count))
		}
	}
	time.Sleep(100 * time.Millisecond)
	refreshData(settings, window, tabs, idx, where, order, &recreate, count, fromPager)
}

func flashCell(table *headertable.HeaderTable, row int, col int) {
	for i := 0; i < 2; i++ {
		if i%2 == 0 {
			table.Data.Select(widget.TableCellID{
				Row: row,
				Col: col,
			})
			time.Sleep(250 * time.Millisecond)
		} else {
			table.Data.Unselect(widget.TableCellID{
				Row: row,
				Col: col,
			})
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func refreshData(
	settings *settings.Settings,
	window fyne.Window,
	tabs *container.AppTabs,
	idx int,
	where interface{},
	order string,
	recreate *bool,
	count int64,
	fromPager bool,
) {
	key := fmt.Sprintf("map%d", idx)
	recreateKey := fmt.Sprintf("map%dRecreate", idx)
	listLength := len(tListMap[key])
	switch idx {
	case 0:
		listMap = &listMapPlay
		if _, ok := cacheListHeader[key]; !ok {
			cacheListHeader[key] = listHeader
		}
	case 1:
		listMap = &listMapPlayed
		if _, ok := cacheListHeader[key]; !ok {
			tListHeader := listHeader
			tListHeader.ColAttrs = make([]headertable.ColAttr, len(listHeader.ColAttrs))
			copy(tListHeader.ColAttrs, listHeader.ColAttrs)
			tListHeader.ColAttrs[7].Name = "PlayTime"
			tListHeader.ColAttrs[7].Header = "游玩时间"
			cacheListHeader[key] = tListHeader
		}
	case 2:
		listMap = &listMapStar
		if _, ok := cacheListHeader[key]; !ok {
			cacheListHeader[key] = listHeader
		}
	case 3:
		listBlacklist = &blacklistList
		if _, ok := cacheListHeader[key]; !ok {
			cacheListHeader[key] = blacklistHeader
		}
		refreshBlacklistData(settings, window, tabs, idx, where, order, recreate, count, fromPager)
		return
	case 4:
		listMap = &searchResult
		if _, ok := cacheListHeader[key]; !ok {
			cacheListHeader[key] = listHeader
		}
	}
	// logger.Debugf("%v", listMap)
	*recreate = *recreate || cache[recreateKey] == "true" || cacheHt[key] == nil
	tListHeader := cacheListHeader[key]
	if listLength > 0 {
		cache[recreateKey] = "false"
		if cacheHt[key] == nil {
			cacheHt[key] = &headertable.HeaderTable{}
		}
		for i := 0; i < pageSize; i++ {
			if i >= listLength {
				(*listMap)[i] = model.MapInfo{}
			} else {
				(*listMap)[i] = tListMap[key][i]
			}
		}
		if !*recreate {
			logger.Debugf("current index: %d, cacheHt: %v", idx, cacheHt)
			if cacheHt[key].Data != nil {
				cacheHt[key].Data.UnselectAll()
			}
			if !settings.SelectedCell {
				cacheHt[key].Refresh()
			}
			logger.Infof("refresh finish, total: %v", len(tListMap[key]))
			return
		}
		bindingsMap[key] = make([]binding.Struct, pageSize)
		for i := 0; i < pageSize; i++ {
			bindingsMap[key][i] = binding.BindStruct(&((*listMap)[i]))
		}
		tListHeader.Bindings = bindingsMap[key]
		cacheHt[key] = headertable.NewHeaderTable(&tListHeader)
		firstItemAction := func() {
			if s, err := bsCellTempString.Get(); err == nil {
				window.Clipboard().SetContent(strings.Split(s, " ")[0])
			}
		}
		firstMenuItem := fyne.NewMenuItem("", firstItemAction)
		playMenuItem := fyne.NewMenuItem("游玩", func() {
			if s, err := bsMapInfoTemp.GetValue("MapId"); err == nil {
				previousCellX = -1
				window.Clipboard().SetContent(s.(string))
				go utils.FillMapId(s.(string), settings)
				go func() {
					db.UpdateMap(
						model.MapInfo{MapId: s.(string), State: "1", PlayTime: time.Now()},
						[]string{"State", "PlayTime"},
						&model.MapInfo{State: "0"})
					RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], false, false)
				}()
			}
		})
		blockMenuItem := fyne.NewMenuItem("拉黑", func() {
			previousCellX = -1
			var uid, nn string
			var err error
			var s interface{}
			if s, err = bsMapInfoTemp.GetValue("Nn"); err == nil {
				nn = s.(string)
			}
			if s, err = bsMapInfoTemp.GetValue("Uid"); err == nil {
				uid = s.(string)
				go func() {
					db.InsertBlacklist(model.Blacklist{
						Uid:     uid,
						Nn:      nn,
						Created: time.Now(),
					})
					RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], false, false)
					dialog.ShowInformation("提示", fmt.Sprintf("已拉黑[%s]，未玩列表不再显示与他相关的", nn), *settings.Window)
				}()
			}
		})
		releaseMenuItem := fyne.NewMenuItem("取消拉黑", func() {
			previousCellX = -1
			var uid, nn string
			var err error
			var s interface{}
			if s, err = bsMapInfoTemp.GetValue("Nn"); err == nil {
				nn = s.(string)
			}
			if s, err = bsMapInfoTemp.GetValue("Uid"); err == nil {
				uid = s.(string)
				go func() {
					db.DeleteBlacklist(model.Blacklist{
						Uid: uid,
					})
					RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], false, false)
					dialog.ShowInformation("提示", fmt.Sprintf("已将[%s]从黑名单移除", nn), *settings.Window)
				}()
			}
		})
		starMenuItem := fyne.NewMenuItem("收藏", func() {
			if s, err := bsMapInfoTemp.GetValue("MapId"); err == nil {
				previousCellX = -1
				window.Clipboard().SetContent(s.(string))
				go utils.FillMapId(s.(string), settings)
				go func() {
					db.UpdateMap(
						model.MapInfo{MapId: s.(string), Star: "1"},
						[]string{"Star"},
						&model.MapInfo{Star: "0"})
					RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], false, false)
				}()
			}
		})
		unStarMenuItem := fyne.NewMenuItem("取消收藏", func() {
			if s, err := bsMapInfoTemp.GetValue("MapId"); err == nil {
				previousCellX = -1
				window.Clipboard().SetContent(s.(string))
				go utils.FillMapId(s.(string), settings)
				go func() {
					db.UpdateMap(
						model.MapInfo{MapId: s.(string), Star: "0"},
						[]string{"Star"},
						&model.MapInfo{Star: "1"})
					RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], false, false)
				}()
			}
		})
		tableMenu := fyne.NewMenu("Actions", firstMenuItem, playMenuItem, starMenuItem, unStarMenuItem)
		cacheHt[key].Header.OnSelected = func(id widget.TableCellID) {
			cacheHt[key].Header.UnselectAll()
		}
		cacheHt[key].Data.OnSelected = func(id widget.TableCellID) {
			// 如果搜索选中了的话，不触发搜索的select
			if settings.SelectedCell && !allowTapManual {
				return
			}
			settings.SelectedCell = false
			allowTapManual = true
			row := bindingsMap[key][id.Row]
			colKey := tListHeader.ColAttrs[id.Col].Name
			if value, err := row.GetValue(colKey); err == nil {
				// 每次点击记录信息
				if mid, err := row.GetValue("MapId"); err == nil {
					_ = bsMapInfoTemp.SetValue("MapId", mid.(string))
					if mid == "" {
						cacheHt[key].Data.UnselectAll()
						return
					}
				}
				if uid, err := row.GetValue("Uid"); err == nil {
					_ = bsMapInfoTemp.SetValue("Uid", uid.(string))
				}
				if nn, err := row.GetValue("Nn"); err == nil {
					_ = bsMapInfoTemp.SetValue("Nn", nn.(string))
				}
				// 处理单元格字符
				converter := tListHeader.ColAttrs[id.Col].Converter
				valueString := ""
				if converter != nil {
					valueString = converter(value, row)
				} else {
					valueString = value.(string)
				}

				star, _ := row.GetValue("Star")
				// 根据收藏状态，调整菜单显示
				deleteTargetItem := starMenuItem
				addTargetItem := unStarMenuItem
				if star == "0" {
					deleteTargetItem = unStarMenuItem
					addTargetItem = starMenuItem
				}
				addOrDeleteMenuItem(tableMenu, deleteTargetItem, addTargetItem)
				// 查库，判断用户是否被拉黑，再处理菜单
				uid, _ := bsMapInfoTemp.GetValue("Uid")
				c := db.CountBlacklist(&model.Blacklist{Uid: uid.(string)})
				deleteTargetItem = blockMenuItem
				addTargetItem = releaseMenuItem
				if c == 0 {
					deleteTargetItem = releaseMenuItem
					addTargetItem = blockMenuItem
				}
				addOrDeleteMenuItem(tableMenu, deleteTargetItem, addTargetItem)

				// 处理点击状态和收藏，移除第一个菜单
				shouldBe := 4
				switch colKey {
				case "State", "Star", "Level":
					if len(tableMenu.Items) == shouldBe {
						_, stream := koazee.StreamOf(tableMenu.Items).Pop()
						tableMenu.Items = stream.Out().Val().([]*fyne.MenuItem)
						shouldBe--
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

				if valueString == "" || valueString == "★" || valueString == "🎮" {
					_, streamOf := koazee.StreamOf(tableMenu.Items).Pop()
					out, _ := streamOf.Sort(func(a, b *fyne.MenuItem) int {
						if len(a.Label) > len(b.Label) {
							return -1
						}
						return 1
					}).Pop()
					valueString = out.Val().(*fyne.MenuItem).Label
				}
				xx, yy := getCellPos(fyne.NewPos(0, 220), id.Col, id.Row, tListHeader, cacheHt[key].RefWidth, 36.1, valueString)
				widget.NewPopUpMenu(tableMenu, window.Canvas()).ShowAtPosition(fyne.NewPos(xx, yy))
			}
			cacheHt[key].Data.UnselectAll()
		}
		tapped := func(pageNo int) {
			previousCellX = 0
			RefreshMapList(settings, window, tabs, idx, "", WhereMap[idx], OrderMap[idx], false, true)
		}
		listPager := pager.NewPager(cacheCurrentNo[key], pPageSize, &count, &tapped)
		cachePager[key] = listPager
		cBorder := container.NewBorder(nil, listPager, nil, nil, cacheHt[key])
		tabs.Items[idx].Content = cBorder
	} else {
		cache[recreateKey] = "true"
		tabs.Items[idx].Content = utils.MakeEmptyList(config.AccentColor)
	}
}

func addOrDeleteMenuItem(tableMenu *fyne.Menu, deleteTargetItem *fyne.MenuItem, addTargetItem *fyne.MenuItem) {
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
}

func getCellPos(base fyne.Position, x int, y int, header headertable.TableOpts, refWidth float32, cellHeight float32, valueString string) (float32, float32) {
	xx := base.X
	yy := base.Y
	for i := 0; i <= x; i++ {
		if i == x {
			cellTextWidth := widget.NewLabel(valueString).MinSize().Width
			switch header.ColAttrs[i].DataStyle.Alignment {
			case fyne.TextAlignLeading:
				// xx += float32(header.ColAttrs[i].WidthPercent)/100*refWidth + 6
			case fyne.TextAlignCenter:
				xx += (float32(header.ColAttrs[i].WidthPercent)/100*refWidth - cellTextWidth) / 2
			case fyne.TextAlignTrailing:
				xx += float32(header.ColAttrs[i].WidthPercent)/100*refWidth - cellTextWidth
			}
		} else {
			xx += float32(header.ColAttrs[i].WidthPercent)/100*refWidth + 6
		}
	}

	return xx, yy + (float32(y)+0.7)*cellHeight - cellHeight
}
