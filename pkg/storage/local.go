package storage

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/romycode/anime-downloader/pkg/errors"
)

type LocalStorage struct {
	path string
	eh   *errors.ErrorHandler
}

// NewLocalStorage creates a new instance of LocalStorage (uses local filesystem)
func NewLocalStorage(path string, eh *errors.ErrorHandler) Storage {
	return &LocalStorage{path: path, eh: eh}
}

// Initialize creates a path that was provided in constructor storage.NewLocalStorage()
func (l LocalStorage) Initialize() {
	err := os.MkdirAll(l.path, 0755)
	if err != nil {
		l.eh.HandleError(err)
	}
}

// CreateFileFromURL create a new file inside the path provided in the constructor
func (l LocalStorage) CreateFileFromURL(name string, url string) error {
	out, err := os.Create(l.path + string(os.PathSeparator) + strings.ReplaceAll(name, " ", "_") + ".mp4")

	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			l.eh.HandleError(err)
		}
	}(out)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
