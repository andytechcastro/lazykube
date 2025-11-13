package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Item struct {
	Text     string
	Selected bool
}

type ListMultiSelection struct {
	*tview.Box
	items            []*Item
	currentItem      int
	itemOffset       int
	horizontalOffset int
	changed          func(index int, text string)
	selectItem       func(selectedItems []string)
	done             func(key *tcell.EventKey) *tcell.EventKey
}

func NewListMultiSelection() *ListMultiSelection {
	lMS := &ListMultiSelection{
		Box:        tview.NewBox(),
		items:      make([]*Item, 0),
		changed:    func(index int, text string) {},
		selectItem: func(selectedItems []string) {},
		done:       func(key *tcell.EventKey) *tcell.EventKey { return key },
	}
	return lMS
}

func (lMS *ListMultiSelection) AddItem(text string) *ListMultiSelection {
	lMS.items = append(lMS.items, &Item{
		Text:     text,
		Selected: true,
	})
	return lMS
}

func (lMS *ListMultiSelection) FindItems(text string) []int {
	return []int{0}
}

func (lMS *ListMultiSelection) Clear() {
	lMS.items = make([]*Item, 0)
}

func (lMS *ListMultiSelection) GetCurrentItem() int {
	return lMS.currentItem
}

func (lMS *ListMultiSelection) SetCurrentItem(index int) {
	lMS.currentItem = index
}

func (lMS *ListMultiSelection) GetItemText(index int) (string, error) {
	return lMS.items[index].Text, nil
}

func (lMS *ListMultiSelection) GetItemSelected(index int) bool {
	return lMS.items[index].Selected
}

func (lMS *ListMultiSelection) GetTextSelectedItems() []string {
	selecteItems := make([]string, 0)
	for _, item := range lMS.items {
		if item.Selected {
			selecteItems = append(selecteItems, item.Text)
		}
	}
	return selecteItems
}

func (lMS *ListMultiSelection) NextItem() {
	lMS.currentItem++
}

func (lMS *ListMultiSelection) PreviewItem() {
	lMS.currentItem--
}

func (lMS *ListMultiSelection) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return lMS.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		lastItem := lMS.currentItem
		switch event.Key() {
		case tcell.KeyEnter:
			if lMS.done != nil {
				lMS.done(event)
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				if lMS.currentItem+1 < len(lMS.items) {
					lMS.currentItem++
				}
			case 'k':
				if lMS.currentItem > 0 {
					lMS.currentItem--
				}
			case ' ':
				lMS.items[lMS.currentItem].Selected = !lMS.items[lMS.currentItem].Selected
				lMS.selectItem(lMS.GetTextSelectedItems())
			case 'c':
				for index := range lMS.items {
					lMS.items[index].Selected = false
				}
				lMS.selectItem(lMS.GetTextSelectedItems())
			case 'f':
				for index := range lMS.items {
					lMS.items[index].Selected = true
				}
				lMS.selectItem(lMS.GetTextSelectedItems())
			}
		}
		if lastItem != lMS.currentItem && lMS.currentItem < len(lMS.items) {
			if lMS.changed != nil {
				item := lMS.items[lMS.currentItem]
				lMS.changed(lMS.currentItem, item.Text)
			}
		}
	})
}

func (lMS *ListMultiSelection) Draw(screen tcell.Screen) {
	lMS.Box.DrawForSubclass(screen, lMS)

	x, y, width, height := lMS.GetInnerRect()
	bottomLimit := y + height
	_, totalHeight := screen.Size()
	if bottomLimit > totalHeight {
		bottomLimit = totalHeight
	}

	if height == 0 {
		return
	}

	if lMS.currentItem < lMS.itemOffset {
		lMS.itemOffset = lMS.currentItem
	} else {
		if lMS.currentItem-lMS.itemOffset >= height {
			lMS.itemOffset = lMS.currentItem + 1 - height
		}
	}
	if lMS.horizontalOffset < 0 {
		lMS.horizontalOffset = 0
	}

	for index, item := range lMS.items {
		if index < lMS.itemOffset {
			continue
		}

		if y >= bottomLimit {
			break
		}

		var checkbox string

		if item.Selected {
			checkbox = "[\u2713] " // Checkmark symbol
		} else {
			checkbox = "[ ] "
		}

		line := checkbox + item.Text
		if index == lMS.currentItem {
			line = "[::r]" + line
		}

		tview.Print(screen, line, x, y, width, tview.AlignLeft, tcell.ColorWhite)
		y++

	}

}

func (lMS *ListMultiSelection) SetChangedFunc(handler func(index int, text string)) *ListMultiSelection {
	lMS.changed = handler
	return lMS
}

func (lMS *ListMultiSelection) SetSelectetItemFunc(handler func(selectedItems []string)) *ListMultiSelection {
	lMS.selectItem = handler
	return lMS
}

func (lMS *ListMultiSelection) SetDoneFunc(handler func(key *tcell.EventKey) *tcell.EventKey) *ListMultiSelection {
	lMS.done = handler
	return lMS
}

