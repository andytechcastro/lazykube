package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewPortForwardModal creates a modal window for entering port-forwarding information.
// It takes a callback function 'onOk' which is called with the port string (e.g., "8080:80")
// when the user clicks OK, or with an empty string on cancellation.
func NewPortForwardModal(onOk func(ports string)) *tview.Grid {
	form := tview.NewForm()
	form.SetBackgroundColor(tcell.ColorGray)
	form.SetBorder(true)
	form.SetTitle("Port Forward")

	// Add an input field for the ports
	form.AddInputField("Ports (local:remote)", "", 40, nil, nil)

	// Add buttons
	form.AddButton("OK", func() {
		inputField := form.GetFormItem(0).(*tview.InputField)
		ports := inputField.GetText()
		onOk(ports)
	})
	form.AddButton("Cancel", func() {
		onOk("") // Indicate cancellation
	})

	grid := tview.NewGrid().
		SetRows(0, 7, 0).                 // Top padding, fixed height, bottom padding
		SetColumns(0, 80, 0).              // Left padding, fixed width, right padding
		AddItem(form, 1, 1, 1, 1, 0, 0, true) // Add form to the center cell

	return grid
}
