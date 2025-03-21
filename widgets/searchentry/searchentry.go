package searchentry

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SearchEntry struct {
	widget.Entry
	OnKeyUp   func(event *fyne.KeyEvent)
	OnKeyDown func(event *fyne.KeyEvent)
	OnTapped  func(event *fyne.PointEvent)

	PreviousCursorColumn int
}

type Tappable interface {
	Tapped(*fyne.PointEvent)
}

func NewSearchEntry(placeHolder string) *SearchEntry {
	entry := &SearchEntry{}
	entry.PlaceHolder = placeHolder
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *SearchEntry) KeyDown(keyEvent *fyne.KeyEvent) {
	if e.OnKeyDown == nil {
		return
	}
	e.OnKeyDown(keyEvent)
}

func (e *SearchEntry) KeyUp(keyEvent *fyne.KeyEvent) {
	if e.OnKeyUp == nil {
		return
	}
	e.OnKeyUp(keyEvent)
}

func (e *SearchEntry) Tapped(event *fyne.PointEvent) {
	if e.OnTapped == nil {
		return
	}
	e.OnTapped(event)
}

func (e *SearchEntry) Focus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(e); c != nil {
		c.Focus(e)
	}
}

func (e *SearchEntry) UnFocus() {
	if c := fyne.CurrentApp().Driver().CanvasForObject(e); c != nil {
		c.Unfocus()
	}
}
