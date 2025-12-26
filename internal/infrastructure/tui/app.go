package tui

import (
	"iter"
	"lazykube/internal/adapter/controller"
	"lazykube/internal/infrastructure/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewApp create the base of the tui
func NewApp(clusters iter.Seq[string],
	controller controller.AppController,
	conf *config.Config,
) {
	// This object is creater for the acces of all items
	resourceDict := NewResourceDict()
	resourceDict.Config = conf

	// Create every item
	mainApp := tview.NewApplication()

	// Set global tview styles for a dark theme
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDarkBlue
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorDarkCyan
	tview.Styles.BorderColor = tcell.ColorWhite
	tview.Styles.TitleColor = tcell.ColorWhite
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.InverseTextColor = tcell.ColorBlack

	yamlView := NewYamlView()
	namespaceList := NewNamespaceList(resourceDict)
	menu := NewMenuClusters(clusters, resourceDict)
	typeList := NewTypeList(resourceDict)
	filterInput := NewFilterInputField(resourceDict)
	tableResource := NewTableResource(resourceDict)
	keybindingView := NewKeybindingView()
	logView := NewLogView()
	pages := tview.NewPages()

	// Save all intems
	resourceDict.View = yamlView
	resourceDict.Menu = menu
	resourceDict.Namespace = namespaceList
	resourceDict.App = mainApp
	resourceDict.Type = typeList
	resourceDict.Pages = pages
	resourceDict.Filter = filterInput
	resourceDict.Table = tableResource
	resourceDict.Keybinding = keybindingView
	resourceDict.LogView = logView
	resourceDict.Controller = &controller

	layout := NewLayout(resourceDict)

	pages.AddPage("main", layout, true, true)
	errorModal := NewErrorModal(resourceDict)
	resourceDict.ErrorModal = errorModal

	keybindingsMap := map[string]string{
		"Clusters":   "[red]space[white]: Select | [red]c[white]: Clear | [red]f[white]: Select All | [red]Enter[white]: Apply",
		"Namespaces": "[red]space[white]: Select | [red]c[white]: Clear | [red]f[white]: Select All | [red]Enter[white]: Apply",
		"Types":      "[red]Enter[white]: Apply",
		"Filter":     "[red]Enter[white]: Apply Filter",
		"Resources":  "[red]l[white]: Logs | [red]y[white]: YAML | [red]e[white]: Exec | [red]p[white]: Port Forward | [red]d[white]: Describe | [red]Del[white]: Delete | [red]Enter[white]: View YAML",
		"YAML View":  "[red]q[white]: Close",
		"Logs":       "[red]Esc[white]: Close",
		"Default":    "[red]1[white]: Clusters | [red]2[white]: Namespaces | [red]3[white]: Types | [red]4[white]: Filter | [red]5[white]: Table | [red]6[white]: YAML | [red]q[white]: Quit",
	}

	setFocus := func(p tview.Primitive) {
		mainApp.SetFocus(p)
		var contextTitle string
		switch p {
		case resourceDict.Menu:
			contextTitle = "Clusters"
		case resourceDict.Namespace:
			contextTitle = "Namespaces"
		case resourceDict.Type:
			contextTitle = "Types"
		case resourceDict.Filter:
			contextTitle = "Filter"
		case resourceDict.Table:
			contextTitle = "Resources"
		case resourceDict.View:
			contextTitle = "YAML View"
		case resourceDict.LogView:
			contextTitle = "Logs"
		default:
			contextTitle = "Default"
		}

		if kb, ok := keybindingsMap[contextTitle]; ok {
			resourceDict.Keybinding.SetKeybindings(kb)
		} else {
			resourceDict.Keybinding.SetKeybindings(keybindingsMap["Default"])
		}
	}

	resourceDict.SetFocus = setFocus

	// Set initial focus to the menu
	mainApp.SetFocus(resourceDict.Menu)
	resourceDict.Keybinding.SetKeybindings(keybindingsMap["Clusters"])

	// keys for the entire application
	mainApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			// Consume Ctrl-C globally to prevent application exit
			return nil
		}
		if event.Rune() == '1' {
			setFocus(resourceDict.Menu)
		}
		if event.Rune() == '2' {
			setFocus(resourceDict.Namespace)
		}
		if event.Rune() == '3' {
			setFocus(resourceDict.Type)
		}
		if event.Rune() == '4' {
			setFocus(resourceDict.Filter)
		}
		if event.Rune() == '5' {
			setFocus(resourceDict.Table)
		}
		if event.Rune() == '6' {
			setFocus(resourceDict.View)
		}
		if event.Rune() == 'q' {
			mainApp.Stop()
		}
		return event
	})

	if err := mainApp.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}