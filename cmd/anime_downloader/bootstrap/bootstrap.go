package bootstrap

import (
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

func WarmUp(path string) (*web.Crawler, storage.Storage) {
	c := web.NewCrawler()
	s := storage.NewLocalStorage(path)

	return c, s
}
