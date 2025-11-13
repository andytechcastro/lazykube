package tui

import (
	"lazykube/entities"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewPodSelectionModal creates a modal for selecting a pod from a list.
// 'onSelect' is called with the chosen pod, or with a zero-value pod if cancelled.
func NewPodSelectionModal(pods []entities.Pod, onSelect func(pod entities.Pod)) (*tview.Flex, *tview.List) {
	list := tview.NewList()
	list.SetSelectedBackgroundColor(tcell.ColorLightSkyBlue)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetBackgroundColor(tcell.ColorGray)
	list.SetBorder(true)
	list.SetTitle("Select a Pod")

	for _, p := range pods {
		pod := p // Capture loop variable
		list.AddItem(pod.Name, "Namespace: "+pod.Namespace, 0, func() {
			onSelect(pod)
		})
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			onSelect(entities.Pod{}) // Zero-value pod indicates cancellation
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			case 'k':
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			}
		}
		return event
	})

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(list, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	flex.SetBackgroundColor(tcell.ColorGray)

	return flex, list
}
