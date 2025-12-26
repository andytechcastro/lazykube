package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewContainerSelectionModal(containers []string, onSelect func(container string)) (*tview.Flex, *tview.List) {
	list := tview.NewList()
	list.SetSelectedBackgroundColor(tcell.ColorLightSkyBlue)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetBackgroundColor(tcell.ColorGray) // Set background color for the list to Gray

	for _, c := range containers {
		container := c
		list.AddItem(container, "", 0, func() {
			onSelect(container)
		})
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			onSelect("") // Empty string indicates cancellation
			return nil    // Consume the event
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				// Simulate a KeyDown event
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			case 'k':
				// Simulate a KeyUp event
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			}
		}
		return event // Return original event for default handling
	})

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(list, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	flex.SetBackgroundColor(tcell.ColorGray) // Set background color for the outer flex to Gray

	return flex, list
}
