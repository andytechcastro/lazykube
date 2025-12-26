package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type typeList struct {
	*tview.List
}

func NewTypeList(dict *resourceDict) *typeList {
	mainList := tview.NewList().ShowSecondaryText(false)
	mainList.AddItem("Deployments", "", rune(0), nil)
	mainList.AddItem("Pods", "", rune(0), nil)

	mainList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, rune(0), tcell.ModNone)
		}
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, rune(0), tcell.ModNone)
		}
		return dict.EventList(event)
	})
	mainList.SetBorder(true).SetTitle("Resource's Type [3]")
	return &typeList{
		mainList,
	}
}
