package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/bridge"
)

var emptyRecvItem = &RecvItem{}

// RecvItem is the item that is being received
type RecvItem struct {
	URI    fyne.URI
	Status string
}

// RecvList is a list of progress bars that track send progress.
type RecvList struct {
	widget.List

	bridge *bridge.Bridge

	Items []RecvItem
}

// Length returns the length of the data.
func (p *RecvList) Length() int {
	return len(p.Items)
}

// CreateItem creates a new item in the list.
func (p *RecvList) CreateItem() fyne.CanvasObject {
	return container.New(&listLayout{}, widget.NewFileIcon(nil), widget.NewLabel("Waiting for filename..."), newRecvProgress())
}

// UpdateItem updates the data in the list.
func (p *RecvList) UpdateItem(i int, item fyne.CanvasObject) {
	item.(*fyne.Container).Objects[0].(*widget.FileIcon).SetURI(p.Items[i].URI)
	item.(*fyne.Container).Objects[1].(*widget.Label).SetText(p.Items[i].URI.Name())
	item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*recvProgress).SetStatus(p.Items[i].Status)
}

// RemoveItem removes the item at the specified index.
func (p *RecvList) RemoveItem(i int) {
	copy(p.Items[i:], p.Items[i+1:])
	p.Items[p.Length()-1] = *emptyRecvItem // Make sure that GC run on removed element
	p.Items = p.Items[:p.Length()-1]
	p.Refresh()
}

// OnSelected handles removing items and stopping send (in the future)
func (p *RecvList) OnSelected(i int) {
	dialog.ShowConfirm("Remove from list", "Do you wish to remove the item from the list?", func(remove bool) {
		if remove {
			p.RemoveItem(i)
			p.Refresh()
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])

	p.Unselect(i)
}

// NewReceive adds data about a new send to the list and then returns the channel to update the code.
func (p *RecvList) NewReceive(code string) {
	p.Items = append(p.Items, RecvItem{URI: storage.NewURI("Waiting for filename..."), Status: "start"})
	p.Refresh()

	uri := make(chan fyne.URI)
	index := p.Length() - 1

	go func() {
		p.Items[index].URI = <-uri
		p.Refresh()
	}()

	go func(code string) {
		if err := p.bridge.NewReceive(code, uri); err != nil {
			p.Items[index].Status = "Failed"
			fyne.LogError("Error on sending file", err)
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
		} else {
			p.Items[index].Status = "Completed"
		}

		p.Refresh()
	}(code)
}

// NewRecvList greates a list of progress bars.
func NewRecvList(bridge *bridge.Bridge) *RecvList {
	p := &RecvList{bridge: bridge}
	p.List.Length = p.Length
	p.List.CreateItem = p.CreateItem
	p.List.UpdateItem = p.UpdateItem
	p.List.OnSelected = p.OnSelected
	p.ExtendBaseWidget(p)

	return p
}
