package donwloader

import (
	"errors"
	"net/url"

	"github.com/romycode/anime-downloader/pkg"
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

var ErrNoDownloaderFound = errors.New("no downloader found")

type AnimeDownloaderFactory struct {
	c *web.Crawler
	s storage.Storage
}

func NewAnimeDownloaderFactory(c *web.Crawler, s storage.Storage) *AnimeDownloaderFactory {
	return &AnimeDownloaderFactory{c: c, s: s}
}

func (a AnimeDownloaderFactory) Build(u url.URL) (pkg.Downloader, error) {
	switch u.Host {
	case string(AnimeFLV2):
		return NewAnimeFlv2Downloader(u, a.c, a.s), nil
	case string(AnimeFLV3):
		return NewAnimeFlv3Downloader(u, a.c, a.s), nil
	default:
		return nil, ErrNoDownloaderFound
	}
}
