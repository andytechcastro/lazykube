package tui

import (
	"github.com/rivo/tview"
)

// KeybindingView is a view that displays keybindings
type KeybindingView struct {
	*tview.TextView
}

// NewKeybindingView creates a new KeybindingView
func NewKeybindingView() *KeybindingView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	textView.SetBorder(true).SetTitle("Keybindings")
	return &KeybindingView{TextView: textView}
}

// SetKeybindings sets the keybindings to be displayed
func (kv *KeybindingView) SetKeybindings(keybindings string) {
	kv.SetText(keybindings)
}
