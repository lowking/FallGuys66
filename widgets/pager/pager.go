package pager

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math"
	"strconv"
)

const limit = 10

var _ fyne.Widget = &Pager{}

type Pager struct {
	widget.BaseWidget

	CurrentPageNo *int
	PageSize      *int
	Total         *int64
	PageCount     *int

	OnTapped *func(pageNo int)

	BtnPreviosPage *widget.Button
	BtnNextPage    *widget.Button

	Items *[]fyne.CanvasObject
}

func NewPager(currentPageNo *int, pageSize *int, total *int64, onTapped *func(pageNo int)) *Pager {
	p := &Pager{
		CurrentPageNo: currentPageNo,
		PageSize:      pageSize,
		Total:         total,
	}
	tap := *onTapped
	p.OnTapped = &tap
	btnPreviousPage := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		p.SelectPage(*p.CurrentPageNo - 1)
		updateBackBtn(p)
	})
	p.BtnPreviosPage = btnPreviousPage
	btnNextPage := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		p.SelectPage(*p.CurrentPageNo + 1)
		updateNextBtn(p)
	})
	p.BtnNextPage = btnNextPage
	p.Items = &[]fyne.CanvasObject{}
	*p.Items = append(*p.Items, layout.NewSpacer())
	*p.Items = append(*p.Items, p.BtnPreviosPage)

	p.Init(p.Total, p.PageSize)
	updateNextBtn(p)
	updateBackBtn(p)
	p.ExtendBaseWidget(p)

	return p
}

func updateNextBtn(p *Pager) {
	if !p.BtnNextPage.Disabled() && *p.CurrentPageNo >= *p.PageCount {
		p.BtnNextPage.Disable()
	}
	if p.BtnPreviosPage.Disabled() && *p.CurrentPageNo > 1 {
		p.BtnPreviosPage.Enable()
	}
}

func updateBackBtn(p *Pager) {
	if !p.BtnPreviosPage.Disabled() && *p.CurrentPageNo <= 1 {
		p.BtnPreviosPage.Disable()
	}
	if p.BtnNextPage.Disabled() && *p.CurrentPageNo < *p.PageCount {
		p.BtnNextPage.Enable()
	}
}

func (p *Pager) Init(total *int64, pageSize *int) {
	ceil := int(math.Ceil(float64(*total) / float64(*pageSize)))
	p.PageCount = &ceil
	start := 2
	half := limit / 2
	start = int(math.Max(2, float64(*p.CurrentPageNo-half)))
	end := int(math.Min(float64(start+limit-1), float64(*p.PageCount-1)))
	start = int(math.Max(2, float64(end-limit+1)))

	// items前2个 占位和前一页按钮
	isFirst := len(*p.Items) == 2
	idx := 3
	if isFirst {
		*p.Items = append(*p.Items, widget.NewButton("首页", func() {
			p.SelectPage(1)
			updateBackBtn(p)
		}))
	}
	for i := start; i <= end; i++ {
		ti := i
		var currentItem *widget.Button
		if isFirst {
			currentItem = widget.NewButton(strconv.Itoa(i), func() {
				p.SelectPage(ti)
				updateNextBtn(p)
				updateBackBtn(p)
			})
			*p.Items = append(*p.Items, currentItem)
		} else {
			currentItem = (*p.Items)[idx].(*widget.Button)
			currentItem.SetText(strconv.Itoa(i))
			currentItem.OnTapped = func() {
				p.SelectPage(ti)
				updateNextBtn(p)
				updateBackBtn(p)
			}
		}
		if *p.CurrentPageNo == i && !currentItem.Disabled() {
			currentItem.Disable()
		}
		if *p.CurrentPageNo != i && currentItem.Disabled() {
			currentItem.Enable()
		}
		idx++
	}
	if isFirst || *p.PageCount+2 > len(*p.Items) {
		*p.Items = append(*p.Items, widget.NewButton("尾页", func() {
			p.SelectPage(*p.PageCount)
			updateNextBtn(p)
		}))
		*p.Items = append(*p.Items, p.BtnNextPage)
	} else if *p.PageCount+2 < len(*p.Items) {
		// 需要清楚尾部的按钮
	}
	// 只有一页，把按钮置灰
	if *p.PageCount == 1 {
		p.BtnPreviosPage.Disable()
		p.BtnNextPage.Disable()
	}
}

func (p *Pager) SelectPage(no int) {
	*p.CurrentPageNo = no
	(*p.OnTapped)(*p.CurrentPageNo)
	p.Init(p.Total, p.PageSize)
}

// ====================================== renderer ======================================
var _ fyne.WidgetRenderer = pagerRenderer{}

type pagerRenderer struct {
	pager     *Pager
	container *fyne.Container
}

func (p *Pager) CreateRenderer() fyne.WidgetRenderer {
	box := container.NewHBox(*p.Items...)
	return pagerRenderer{
		pager:     p,
		container: box,
	}
}

func (r pagerRenderer) MinSize() fyne.Size {
	minWidth := float32(0)
	for _, object := range r.Objects() {
		minWidth += object.Size().Width
	}
	return fyne.NewSize(minWidth-11, 35)
}

func (r pagerRenderer) Layout(s fyne.Size) {
	r.container.Resize(s)
}

func (r pagerRenderer) Destroy() {
}

func (r pagerRenderer) Refresh() {
}

func (r pagerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.container}
}
