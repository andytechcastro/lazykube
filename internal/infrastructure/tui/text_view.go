package tui

import "github.com/rivo/tview"

type textView struct {
	*tview.TextView
}

func NewTextView() *textView {
	mainTextView := tview.NewTextView()
	return &textView{
		mainTextView,
	}
}
