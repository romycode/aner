package bootstrap

import (
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

func WarmUp(path string) (*web.URLExtractor, storage.Storage) {
	var e = web.NewURLExtractor()

	localStorage := storage.NewLocalStorage(path)

	return e, localStorage
}
