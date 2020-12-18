package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/bridge"
	"github.com/Jacalz/wormhole-gui/internal/bridge/widgets"
)

type send struct {
	contentPicker dialog.Dialog

	fileChoice      *widget.Button
	fileDialog      *dialog.FileDialog
	directoryChoice *widget.Button
	directoryDialog *dialog.FileDialog
	textChoice      *widget.Button

	contentToSend *widget.Button
	sendList      *widgets.SendList

	bridge      *bridge.Bridge
	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newSend(a fyne.App, w fyne.Window, b *bridge.Bridge, as *AppSettings) *send {
	return &send{app: a, window: w, bridge: b, appSettings: as}
}

func (s *send) onFileSend() {
	s.contentPicker.Hide()
	s.fileDialog.Show()
}

func (s *send) onDirSend() {
	s.contentPicker.Hide()
	s.directoryDialog.Show()
}

func (s *send) onTextSend() {
	s.contentPicker.Hide()
	s.sendList.SendText()
}

func (s *send) onContentToSend() {
	s.contentPicker.Show()
}

func (s *send) buildUI() *fyne.Container {
	s.fileChoice = &widget.Button{Text: "File", Icon: theme.FileIcon(), OnTapped: s.onFileSend}
	s.directoryChoice = &widget.Button{Text: "Directory", Icon: theme.FolderOpenIcon(), OnTapped: s.onDirSend}
	s.textChoice = &widget.Button{Text: "Text", Icon: theme.DocumentCreateIcon(), OnTapped: s.onTextSend}

	choiceContent := container.NewGridWithColumns(1, s.fileChoice, s.directoryChoice, s.textChoice)
	s.contentPicker = dialog.NewCustom("Pick a content type", "Cancel", choiceContent, s.window)

	s.sendList = widgets.NewSendList(s.bridge)
	s.contentToSend = &widget.Button{Text: "Add content to send", Icon: theme.ContentAddIcon(), OnTapped: s.onContentToSend}

	s.fileDialog = dialog.NewFileOpen(s.sendList.OnFileSelect, s.window)
	s.directoryDialog = dialog.NewFolderOpen(s.sendList.OnDirSelect, s.window)

	box := container.NewVBox(s.contentToSend, &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, s.sendList)
}

func (s *send) tabItem() *container.TabItem {
	return container.NewTabItemWithIcon("Send", theme.MailSendIcon(), s.buildUI())
}
