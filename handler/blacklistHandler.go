package handler

import (
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
			WidthPercent: 400,
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
			WidthPercent: 352,
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
				valueString := ""
				if s, ok := value.(string); ok {
					valueString = s
				} else if s, ok := value.(time.Time); ok {
					valueString = s.Format("2006-01-02 15:04:05")
				}
				tableMenu.Items[0].Label = valueString
				_ = bsCellTempString.Set(valueString)

				xx, yy := getCellPos(fyne.NewPos(0, 220), id.Col, id.Row, tListHeader, 36.1)
				widget.NewPopUpMenu(tableMenu, window.Canvas()).ShowAtPosition(fyne.NewPos(xx+float32(tListHeader.ColAttrs[id.Col].WidthPercent/2), yy))
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
		tabs.Items[idx].Content = utils.MakeEmptyList(config.AccentColor)
	}
}
