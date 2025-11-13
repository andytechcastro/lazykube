package tui

import (
	"io"
	// Added context import
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/client-go/tools/remotecommand"
)

// TerminalView es un componente TUI que simula una terminal para la ejecución de comandos en pods.
type TerminalView struct {
	*tview.Flex
	textView *tview.TextView
	input    *Input
	app      *tview.Application

	// Pipes para streaming de I/O
	stdinReader  *io.PipeReader
	stdinWriter  *io.PipeWriter
	stdoutPipe   *io.PipeReader
	stdoutWriter *io.PipeWriter
}

// NewTerminalView crea una nueva TerminalView.
func NewTerminalView(app *tview.Application, closeFn func()) *TerminalView {
	v := &TerminalView{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
		app:  app,
	}

	// Configuración de pipes
	stdinReader, stdinWriter := io.Pipe()
	v.stdinReader = stdinReader
	v.stdinWriter = stdinWriter
	v.stdoutPipe, v.stdoutWriter = io.Pipe()

	// Configuración de TextView para el historial de salida
	v.textView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			v.app.Draw()
			v.textView.ScrollToEnd()
		})

	v.input = NewInput(v.stdinWriter)

	// Leer continuamente del pipe de stdout y escribir en el TextView
	go func() {
		if _, err := io.Copy(tview.ANSIWriter(v.textView), v.stdoutPipe); err != nil {
			// Manejar error, quizás mostrar en un modal
		}
	}()

	v.SetBorder(true).SetTitle("Pod Terminal (Presiona Esc para cerrar)")

	// Manejar la entrada del teclado para el contenedor Flex (solo Esc y Ctrl+C para cerrar/interrumpir)
	v.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			closeFn()
			return nil
		case tcell.KeyCtrlC:
			// Enviar ASCII ETX para Ctrl+C al shell remoto
			_, _ = v.stdinWriter.Write([]byte{''})
			return nil // Consumir el evento
		}
		return event
	})

	inputArea := tview.NewFlex().SetDirection(tview.FlexColumn)
	inputArea.AddItem(tview.NewTextView().SetText("$ "), 2, 0, false)
	inputArea.AddItem(v.input, 0, 1, true)

	// Añadir TextView e InputField al layout Flex
	v.AddItem(v.textView, 0, 1, false) // TextView ocupa la mayor parte del espacio, no enfocable por defecto
	v.AddItem(inputArea, 3, 0, true)    // InputField ocupa 3 línea, es enfocable

	// Establecer el foco inicial en el campo de entrada
	v.app.SetFocus(v.input)

	return v
}

// GetStreamOptions devuelve las opciones de stream para el comando remoto.
func (v *TerminalView) GetStreamOptions() remotecommand.StreamOptions {
	return remotecommand.StreamOptions{
		Stdin:  v.stdinReader,
		Stdout: v.stdoutWriter,
		Stderr: v.stdoutWriter, // Enviar stderr al mismo writer que stdout
		Tty:    true,
	}
}

// Stop cierra los pipes.
func (v *TerminalView) Stop() {
	_ = v.stdinWriter.Close()
	_ = v.stdoutPipe.Close()
	_ = v.stdoutWriter.Close()
}
