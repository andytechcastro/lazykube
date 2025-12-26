package tui

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"lazykube/internal/adapter/controller"
	"lazykube/internal/domain"
	"lazykube/internal/infrastructure/config"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"k8s.io/client-go/tools/remotecommand"
)

var onceDict sync.Once

// The struct for save all items
type resourceDict struct {
	App        *tview.Application
	Pages      *tview.Pages
	Menu       *menuClusters
	View       *yamlView
	Namespace  *namespaceList
	Type       *typeList
	Controller *controller.AppController
	ErrorModal *errorModal
	Filter     *filterInputField
	Table      *tableResource
	Config     *config.Config
	Keybinding *KeybindingView
	LogView    *LogView
	SetFocus   func(p tview.Primitive)
}

// for singleton
var singleInstance *resourceDict

func NewResourceDict() *resourceDict {
	if singleInstance == nil {
		onceDict.Do(
			func() {
				singleInstance = &resourceDict{}
			})
	}
	return singleInstance
}

func (rD *resourceDict) showPortForwardForPod(pod domain.Pod) {
	rD.App.QueueUpdateDraw(func() {
		modal := NewPortForwardModal(func(ports string) {
			rD.Pages.RemovePage("portForward")
			if ports == "" {
				rD.SetFocus(rD.Table)
				return
			}

			if !strings.Contains(ports, ":") {
				rD.ErrorModal.SetText("Invalid port format. Use local:remote.")
				rD.Pages.ShowPage("errorModal")
				return
			}

			stopChan := make(chan struct{}, 1)
			readyChan := make(chan struct{}, 1)

			stdout, stderr, err := rD.Controller.Pod.PortForward(context.Background(), pod.Name, pod.Namespace, pod.Context, []string{ports}, stopChan, readyChan)
			if err != nil {
				rD.ErrorModal.SetText(fmt.Sprintf("Port forward setup failed: %v", err))
				rD.Pages.ShowPage("errorModal")
				return
			}

			go func() {
				select {
				case <-readyChan:
					rD.App.QueueUpdateDraw(func() {
						time.Sleep(250 * time.Millisecond)
						modalText := stdout.String()
						if stderr.Len() > 0 {
							modalText += "\nErrors:\n" + stderr.String()
						}
						if modalText == "" {
							modalText = "Port forward established. No output from command."
						}

						infoModal := tview.NewModal().
							SetText(modalText).
							AddButtons([]string{"Stop"}).
							SetDoneFunc(func(buttonIndex int, buttonLabel string) {
								close(stopChan)
								rD.Pages.HidePage("portForwardInfo")
								rD.SetFocus(rD.Table)
							})
						rD.Pages.AddPage("portForwardInfo", infoModal, true, true)
						rD.SetFocus(infoModal)
					})
				case <-time.After(10 * time.Second):
					close(stopChan)
					rD.App.QueueUpdateDraw(func() {
						rD.ErrorModal.SetText("Port forward timed out.")
						rD.Pages.ShowPage("errorModal")
					})
				}
			}()
		})
		rD.Pages.AddPage("portForward", modal, true, true)
		rD.SetFocus(modal)
	})
}

func (rD *resourceDict) showPortForwardForDeployment(deploymentName, namespace, contextStr string) {
	pods, err := rD.Controller.Deployment.GetPods(context.Background(), deploymentName, namespace, contextStr)
	rD.App.QueueUpdateDraw(func() {
		if err != nil {
			rD.ErrorModal.SetText(fmt.Sprintf("Failed to get pods: %v", err))
			rD.Pages.ShowPage("errorModal")
			return
		}

		if len(pods) == 0 {
			rD.ErrorModal.SetText("No pods found for this deployment.")
			rD.Pages.ShowPage("errorModal")
			return
		}

		if len(pods) == 1 {
			go rD.showPortForwardForPod(pods[0])
			return
		}

		modal, list := NewPodSelectionModal(pods, func(pod domain.Pod) {
			rD.Pages.RemovePage("podSelection")
			if pod.Name != "" {
				go rD.showPortForwardForPod(pod)
			} else {
				rD.SetFocus(rD.Table)
			}
		})
		rD.Pages.AddPage("podSelection", modal, true, true)
		rD.SetFocus(list)
	})
}

func (rD *resourceDict) showLogsForPod(pod domain.Pod, containerName string) {
	var getLogsFn func(string)
	getLogsFn = func(cName string) {
		ctx, cancel := context.WithCancel(context.Background())

		logStream, err := rD.Controller.Pod.GetLogs(ctx, pod.Name, pod.Namespace, pod.Context, cName)
		rD.App.QueueUpdateDraw(func() {
			if err != nil {
				if containerErr, ok := err.(*domain.ContainerSelectionError); ok {
					modal, list := NewContainerSelectionModal(containerErr.Containers, func(container string) {
						rD.Pages.RemovePage("containerSelection")
						if container != "" {
							go getLogsFn(container)
						} else {
							rD.SetFocus(rD.Table)
						}
					})
					rD.Pages.AddPage("containerSelection", modal, true, true)
					rD.SetFocus(list)
				} else {
					rD.ErrorModal.SetText(err.Error())
					rD.Pages.ShowPage("errorModal")
				}
				return
			}

			logView := NewLogView()
			logView.SetTitle(fmt.Sprintf("Logs for %s", pod.Name))
			logPageName := "logs-" + pod.Name

			closeFn := func() {
				cancel()
				rD.Pages.RemovePage(logPageName)
				rD.SetFocus(rD.Table)
			}

			logView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
					closeFn()
				}
				return event
			})

			go func() {
				defer logStream.Close()
				scanner := bufio.NewScanner(logStream)
				for scanner.Scan() {
					line := scanner.Text()
					rD.App.QueueUpdateDraw(func() {
						fmt.Fprintf(logView, "%s\n", line)
						logView.ScrollToEnd()
					})
				}
			}()

			rD.Pages.AddPage(logPageName, logView, true, true)
			rD.SetFocus(logView)
		})
	}
	getLogsFn(containerName)
}

func (rD *resourceDict) showLogsForDeployment(deploymentName, namespace, contextStr string) {
	pods, err := rD.Controller.Deployment.GetPods(context.Background(), deploymentName, namespace, contextStr)
	rD.App.QueueUpdateDraw(func() {
		if err != nil {
			rD.ErrorModal.SetText(fmt.Sprintf("Failed to get pods: %v", err))
			rD.Pages.ShowPage("errorModal")
			return
		}

		if len(pods) == 0 {
			rD.ErrorModal.SetText("No pods found for this deployment.")
			rD.Pages.ShowPage("errorModal")
			return
		}

		if len(pods) == 1 {
			go rD.showLogsForPod(pods[0], "")
			return
		}

		modal, list := NewPodSelectionModal(pods, func(pod domain.Pod) {
			rD.Pages.RemovePage("podSelection")
			if pod.Name != "" {
				go rD.showLogsForPod(pod, "")
			} else {
				rD.SetFocus(rD.Table)
			}
		})
		rD.Pages.AddPage("podSelection", modal, true, true)
		rD.SetFocus(list)
	})
}

func (rD *resourceDict) showExecForPod(pod domain.Pod, containerName string) {
	showTerminal := func(cName string) {
		terminalPageName := fmt.Sprintf("terminal-%s-%s", pod.Name, cName)
		closeFn := func() {
			rD.Pages.RemovePage(terminalPageName)
			rD.SetFocus(rD.Table)
		}
		terminalView := NewTerminalView(rD.App, closeFn)
		streamOptions := terminalView.GetStreamOptions()

		go func() {
			defer closeFn()
			err := rD.Controller.Pod.Exec(context.Background(), pod.Name, pod.Namespace, pod.Context, "sh", cName, false, streamOptions)
			if err != nil {
				rD.App.QueueUpdateDraw(func() {
					rD.ErrorModal.SetText(err.Error())
					rD.Pages.ShowPage("errorModal")
				})
			}
		}()
		rD.Pages.AddPage(terminalPageName, terminalView, true, true)
		rD.SetFocus(terminalView)
	}

	// This first Exec call is a dry run to get container names
	err := rD.Controller.Pod.Exec(context.Background(), pod.Name, pod.Namespace, pod.Context, "sh", containerName, true, remotecommand.StreamOptions{})
	rD.App.QueueUpdateDraw(func() {
		if err != nil {
			if containerErr, ok := err.(*domain.ContainerSelectionError); ok {
				modal, list := NewContainerSelectionModal(containerErr.Containers, func(container string) {
					rD.Pages.RemovePage("containerSelection")
					if container != "" {
						showTerminal(container)
					} else {
						rD.SetFocus(rD.Table)
					}
				})
				rD.Pages.AddPage("containerSelection", modal, true, true)
				rD.SetFocus(list)
			} else {
				rD.ErrorModal.SetText(err.Error())
				rD.Pages.ShowPage("errorModal")
			}
		} else {
			showTerminal(containerName)
		}
	})
}

func (rD *resourceDict) showExecForDeployment(deploymentName, namespace, contextStr string) {
	pods, err := rD.Controller.Deployment.GetPods(context.Background(), deploymentName, namespace, contextStr)
	rD.App.QueueUpdateDraw(func() {
		if err != nil {
			rD.ErrorModal.SetText(fmt.Sprintf("Failed to get pods: %v", err))
			rD.Pages.ShowPage("errorModal")
			return
		}

		if len(pods) == 0 {
			rD.ErrorModal.SetText("No pods found for this deployment.")
			rD.Pages.ShowPage("errorModal")
			return
		}

		if len(pods) == 1 {
			go rD.showExecForPod(pods[0], "")
			return
		}

		modal, list := NewPodSelectionModal(pods, func(pod domain.Pod) {
			rD.Pages.RemovePage("podSelection")
			if pod.Name != "" {
				go rD.showExecForPod(pod, "")
			} else {
				rD.SetFocus(rD.Table)
			}
		})
		rD.Pages.AddPage("podSelection", modal, true, true)
		rD.SetFocus(list)
	})
}

// The function for fill the resource table, used for many items
func (rD *resourceDict) UpdateResources() {
	// Get the info from lists and filter
	contexts := rD.Menu.GetTextSelectedItems()
	namespaces := rD.Namespace.GetTextSelectedItems()
	typeR, _ := rD.Type.GetItemText(rD.Type.GetCurrentItem())
	filter := rD.Filter.GetText()

	// The table must know the type of resource that show
	rD.Table.ResourceType = typeR
	rD.Table.Clear()

	// Fill table headers
	rD.Table.SetCell(0, 0, tview.NewTableCell("NAME").SetSelectable(false))
	rD.Table.SetCell(0, 1, tview.NewTableCell("NAMESPACE").SetSelectable(false))
	rD.Table.SetCell(0, 2, tview.NewTableCell("CLUSTER").SetSelectable(false))
	results := map[string][]map[string]string{}
	switch typeR {
	case "Deployments":
		results, _ = rD.Controller.Deployment.GetFromManyContext(context.Background(), namespaces, contexts)
	case "Pods":
		results, _ = rD.Controller.Pod.GetFromManyContext(context.Background(), namespaces, contexts)
	}

	// Fill the resource table, depends if have or not a filter
	c := 1
	if filter == "" {
		for cluster, res := range results {
			for _, data := range res {
				rD.Table.SetCell(c, 0, tview.NewTableCell(data["name"]))
				rD.Table.SetCell(c, 1, tview.NewTableCell(data["namespace"]))
				rD.Table.SetCell(c, 2, tview.NewTableCell(cluster))
				c++
			}
		}
	} else {
		for cluster, res := range results {
			for _, data := range res {
				if strings.Contains(data["name"], filter) {
					rD.Table.SetCell(c, 0, tview.NewTableCell(data["name"]))
					rD.Table.SetCell(c, 1, tview.NewTableCell(data["namespace"]))
					rD.Table.SetCell(c, 2, tview.NewTableCell(cluster))
					c++
				}
			}
		}
	}
	jsonResult, _ := json.MarshalIndent(results, "", " ")
	rD.View.SetText(string(jsonResult))
	rD.SetFocus(rD.Table)
}

// The event keys for all the list
func (rd *resourceDict) EventList(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEnter {
		rd.UpdateResources()
		return event
	}
	return event
}
