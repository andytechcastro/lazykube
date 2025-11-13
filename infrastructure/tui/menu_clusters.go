package tui

import (
	"iter"
)

type menuClusters struct {
	*ListMultiSelection
}

func NewMenuClusters(valList iter.Seq[string],
	dict *resourceDict,
) *menuClusters {
	mainList := NewListMultiSelection()
	for cluster := range valList {
		mainList.AddItem(cluster)
	}
	mainList.SetDoneFunc(dict.EventList)
	mainList.SetSelectetItemFunc(func(selectedItems []string) {
		dict.Namespace.Clear()
		if len(selectedItems) > 1 || len(selectedItems) == 0 {
			dict.Namespace.AddItem("default")
			for _, namespace := range dict.Config.DefaultNamespaces {
				dict.Namespace.AddItem(namespace)
			}
		} else {
			context := selectedItems[0]
			namespaces, err := dict.Controller.Namespace.GetAll(context)
			if err != nil {
				dict.ErrorModal.SetText(err.Error())
				dict.Pages.ShowPage("errorModal")
				return
			}
			for _, namespace := range namespaces {
				dict.Namespace.AddItem(namespace)
			}
			intList := dict.Namespace.FindItems("default")
			dict.Namespace.SetCurrentItem(intList[0])
		}
	})
	mainList.SetBorder(true).SetTitle("Clusters [1]")
	return &menuClusters{
		mainList,
	}
}
