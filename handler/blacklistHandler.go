package handler

import (
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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/wesovilabs/koazee"
	"time"
)

var blacklistHeader = headertable.TableOpts{
	RefWidth: "reference width",
	ColAttrs: []headertable.ColAttr{
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
			WidthPercent: 500,
			Converter: func(i interface{}, row binding.Struct) string {
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
				Alignment: fyne.TextAlignCenter,
			},
			WidthPercent: 343,
		},
		{
			Name:   "Created",
			Header: "拉黑时间",
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

var blacklistList = [pageSize]model.Blacklist{}
var listBlacklist *[pageSize]model.Blacklist

func refreshBlacklistData(
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
	listLength := len(tBlacklist)

	*recreate = *recreate || cache[recreateKey] == "true" || cacheHt[key] == nil
	tListHeader := cacheListHeader[key]
	if listLength > 0 {
		cache[recreateKey] = "false"
		if cacheHt[key] == nil {
			cacheHt[key] = &headertable.HeaderTable{}
		}
		for i := 0; i < pageSize; i++ {
			if i >= listLength {
				(*listBlacklist)[i] = model.Blacklist{}
			} else {
				(*listBlacklist)[i] = tBlacklist[i]
			}
		}
		if !*recreate {
			logger.Debugf("current index: %d, cacheHt: %v", idx, cacheHt)
			if cacheHt[key].Data != nil {
				cacheHt[key].Data.UnselectAll()
			}
			cacheHt[key].Refresh()
			logger.Infof("refresh finish, total: %v", len(tBlacklist))
			return
		}
		bindingsMap[key] = make([]binding.Struct, pageSize)
		for i := 0; i < pageSize; i++ {
			bindingsMap[key][i] = binding.BindStruct(&((*listBlacklist)[i]))
		}
		tListHeader.Bindings = bindingsMap[key]
		cacheHt[key] = headertable.NewHeaderTable(&tListHeader)
		firstItemAction := func() {
			if s, err := bsCellTempString.Get(); err == nil {
				window.Clipboard().SetContent(s)
			}
		}
		firstMenuItem := fyne.NewMenuItem("", firstItemAction)
		releaseMenuItem := fyne.NewMenuItem("取消拉黑", func() {
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
		tableMenu := fyne.NewMenu("Actions", firstMenuItem, releaseMenuItem)
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
				if uid, err := row.GetValue("Uid"); err == nil {
					_ = bsMapInfoTemp.SetValue("Uid", uid.(string))
					if uid == "" {
						cacheHt[key].Data.UnselectAll()
						return
					}
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
				tableMenu.Items[0].Label = valueString
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
			RefreshMapList(settings, window, tabs, idx, cacheKeyword, WhereMap[idx], OrderMap[idx], false, true)
		}
		listPager := pager.NewPager(cacheCurrentNo[key], pPageSize, &count, &tapped)
		cachePager[key] = listPager
		cBorder := container.NewBorder(nil, listPager, nil, nil, cacheHt[key])
		tabs.Items[idx].Content = cBorder
	} else {
		cache[recreateKey] = "true"
		tabs.Items[idx].Content = utils.MakeEmptyList(theme.PrimaryColor())
	}
}
