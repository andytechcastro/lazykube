package tui

import (
	"github.com/rivo/tview"
)

type yamlView struct {
	*tview.TextView
}

func NewYamlView() *yamlView {
	mainTextView := tview.NewTextView()
	mainTextView.SetBorder(true).SetTitle("Info [6]")

	return &yamlView{
		mainTextView,
	}
}
