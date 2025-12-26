package tui

import (
	"github.com/rivo/tview"
)

type LogView struct {
	*tview.TextView
}

func NewLogView() *LogView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	textView.SetBorder(true).SetTitle("Logs")
	return &LogView{
		TextView: textView,
	}
}

func (y *LogView) SetContent(content string) {
	y.SetText(content)
	y.ScrollToEnd()
}
