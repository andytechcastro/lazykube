package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Modal for show errors
type errorModal struct {
	*tview.Modal
}

// Return the modal for error
func NewErrorModal(dict *resourceDict) *errorModal {
	modal := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			dict.Pages.HidePage("errorModal")
			dict.App.SetFocus(dict.Menu)
		}).SetBackgroundColor(tcell.ColorRed).
		SetTextColor(tcell.ColorBlack)
	dict.Pages.AddPage("errorModal", modal, true, false)

	return &errorModal{
		modal,
	}
}
