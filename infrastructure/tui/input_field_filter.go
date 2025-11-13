package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type filterInputField struct {
	*tview.InputField
}

func NewFilterInputField(dict *resourceDict) *filterInputField {
	inputField := tview.NewInputField().
		SetLabel("Filter: ").
		SetFieldBackgroundColor(tcell.ColorGray)
	inputField.SetBorder(true).
		SetTitle("Filter [4]")

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			dict.UpdateResources()
		}
		return event
	})

	return &filterInputField{
		inputField,
	}
}
