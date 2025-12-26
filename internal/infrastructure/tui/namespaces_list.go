package tui



type namespaceList struct {
	*ListMultiSelection
}

func NewNamespaceList(dict *resourceDict) *namespaceList {
	mainList := NewListMultiSelection()
	for _, namespace := range dict.Config.DefaultNamespaces {
		mainList.AddItem(namespace)
	}
	mainList.SetCurrentItem(1)
	mainList.SetBorder(true).SetTitle("Namespaces [2]")

	mainList.SetDoneFunc(dict.EventList)

	return &namespaceList{
		mainList,
	}
}
