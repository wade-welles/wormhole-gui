package bridge

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
	"github.com/mholt/archiver/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

// NewReceive runs a receive using wormhole-william and handles types accordingly.
func (b *Bridge) NewReceive(code string, uri chan fyne.URI) error {
	msg, err := b.Receive(context.Background(), code)
	if err != nil {
		fyne.LogError("Error on receiving data", err)
		msg.Reject()
		return err
	}

	if msg.Type == wormhole.TransferText {
		content, err := ioutil.ReadAll(msg)
		if err != nil {
			fyne.LogError("Error on reading received data", err)
			return err
		}

		b.displayReceivedText(content)

		uri <- storage.NewURI("Text Snippet")
	}

	path := filepath.Join(b.DownloadPath, msg.Name)

	switch msg.Type {
	case wormhole.TransferFile:
		file, err := os.Create(path)
		if err != nil {
			fyne.LogError("Error on creating file", err)
			msg.Reject()
			return err
		}

		defer func() {
			if cerr := file.Close(); cerr != nil {
				fyne.LogError("Error on closing file", err)
				err = cerr
			}
		}()

		_, err = io.Copy(file, ioutil.NopCloser(msg))
		if err != nil {
			fyne.LogError("Error on copying contents to file", err)
			return err
		}

		uri <- storage.NewFileURI(path)
	case wormhole.TransferDirectory:
		tmp, err := ioutil.TempFile("", msg.Name+".zip.tmp")
		if err != nil {
			fyne.LogError("Error on creating tempfile", err)
			msg.Reject()
			return err
		}

		defer func() {
			if cerr := tmp.Close(); cerr != nil {
				fyne.LogError("Error on closing file", err)
				err = cerr
			}

			if rerr := os.Remove(tmp.Name()); rerr != nil {
				fyne.LogError("Error on removing temp file", err)
				err = rerr
			}
		}()

		_, err = io.Copy(tmp, ioutil.NopCloser(msg))
		if err != nil {
			fyne.LogError("Error on copying contents to file", err)
			return err
		}

		err = archiver.NewZip().Unarchive(tmp.Name(), path)
		if err != nil {
			fyne.LogError("Error on unzipping contents", err)
			return err
		}

		uri <- storage.NewFileURI(path)
	}

	return nil
}
