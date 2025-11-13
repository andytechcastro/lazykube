package tui

import (
	"lazykube/entities"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tableResource struct {
	*tview.Table
	ResourceType string
}

func NewTableResource(dict *resourceDict) *tableResource {
	table := tview.NewTable()
	table.SetBorder(true).SetTitle("Table Resources [5]")
	table.SetSelectable(true, false)
	table.SetFixed(1, 0)
	table.SetSelectedFunc(func(row int, col int) {
		name := table.GetCell(row, col).Text
		namespace := table.GetCell(row, 1).Text
		kubeContext := table.GetCell(row, 2).Text
		typeR := dict.Table.ResourceType

		results := []byte{}
		var err error
		switch typeR {
		case "Deployments":
			results, err = dict.Controller.Deployment.GetYaml(namespace, name, kubeContext)
		case "Pods":
			results, err = dict.Controller.Pod.GetYaml(namespace, name, kubeContext)
		}
		dict.View.Clear()
		if err != nil {
			dict.ErrorModal.SetText(err.Error())
			dict.Pages.ShowPage("errorModal")
		} else {
			dict.View.SetText(string(results))
			dict.App.SetFocus(dict.View)
		}
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		if row < 1 { // No item selected
			return event
		}
		name := table.GetCell(row, 0).Text
		namespace := table.GetCell(row, 1).Text
		kubeContext := table.GetCell(row, 2).Text
		typeR := dict.Table.ResourceType

		pod := entities.Pod{
			Name:      name,
			Namespace: namespace,
			Context:   kubeContext,
		}

		switch event.Rune() {
		case 'l':
			go func() {
				if typeR == "Pods" {
					dict.showLogsForPod(pod, "")
				} else if typeR == "Deployments" {
					dict.showLogsForDeployment(name, namespace, kubeContext)
				}
			}()
		case 'e':
			go func() {
				if typeR == "Pods" {
					dict.showExecForPod(pod, "")
				} else if typeR == "Deployments" {
					dict.showExecForDeployment(name, namespace, kubeContext)
				}
			}()
		case 'y':
			results := []byte{}
			var err error
			switch typeR {
			case "Deployments":
				results, err = dict.Controller.Deployment.GetYaml(namespace, name, kubeContext)
			case "Pods":
				results, err = dict.Controller.Pod.GetYaml(namespace, name, kubeContext)
			}
			dict.View.Clear()
			if err != nil {
				dict.ErrorModal.SetText(err.Error())
				dict.Pages.ShowPage("errorModal")
			} else {
				dict.View.SetText(string(results))
				dict.App.SetFocus(dict.View)
			}
		case 'p':
			go func() {
				switch typeR {
				case "Pods":
					dict.showPortForwardForPod(pod)
				case "Deployments":
					dict.showPortForwardForDeployment(name, namespace, kubeContext)
				}
			}()
		}
		return event
	})
	return &tableResource{
		Table: table,
	}
}
