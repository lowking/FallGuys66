package searchentry

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type searchEntry struct {
	widget.Entry
	OnKeyUp   func(event *fyne.KeyEvent)
	OnKeyDown func(event *fyne.KeyEvent)
}

func NewSearchEntry(placeHolder string) *searchEntry {
	entry := &searchEntry{}
	entry.PlaceHolder = placeHolder
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *searchEntry) KeyDown(keyEvent *fyne.KeyEvent) {
	if e.OnKeyDown == nil {
		return
	}
	e.OnKeyDown(keyEvent)
}

func (e *searchEntry) KeyUp(keyEvent *fyne.KeyEvent) {
	if e.OnKeyUp == nil {
		return
	}
	e.OnKeyUp(keyEvent)
}
