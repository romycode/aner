package bootstrap

import (
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

func WarmUp(path string) (*web.Crawler, storage.Storage) {
	var e = web.NewCrawler()

	localStorage := storage.NewLocalStorage(path)

	return e, localStorage
}
