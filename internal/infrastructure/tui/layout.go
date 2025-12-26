package tui

import (
	"github.com/rivo/tview"
)

type layout struct {
	*tview.Flex
}

func NewLayout(dict *resourceDict) *layout {
	mainContentFlex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(dict.Menu, 0, 1, true).
			AddItem(dict.Namespace, 0, 1, false).
			AddItem(dict.Type, 0, 2, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(dict.Filter, 0, 1, false).
			AddItem(dict.Table, 0, 5, false), 0, 3, false).
		AddItem(dict.View, 0, 3, false)

	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainContentFlex, 0, 1, false).
		AddItem(dict.Keybinding, 3, 0, false)

	return &layout{
		rootFlex,
	}
}
