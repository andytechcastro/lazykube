package tui

import (
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Input struct {
	*tview.TextArea
}

func NewInput(stdinWriter io.Writer) *Input {
	textArea := tview.NewTextArea()
	textArea.SetBackgroundColor(tcell.ColorGray)

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			line := textArea.GetText()
			// Escribir comando + salto de línea al pipe de stdin
			_, _ = stdinWriter.Write([]byte(line + "\r"))
			textArea.SetText("", true) // Limpiar campo de entrada
			return nil
		}
		return event
	})

	return &Input{
		TextArea: textArea,
	}
}
