package donwloader

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

type AnimeFlv2Downloader struct {
	u url.URL
	c *web.Crawler
	s storage.Storage
}

func NewAnimeFlv2Downloader(animeURL url.URL, ue *web.Crawler, s storage.Storage) *AnimeFlv2Downloader {
	return &AnimeFlv2Downloader{
		u: animeURL,
		c: ue,
		s: s,
	}
}

func (a AnimeFlv2Downloader) GetEpisodes() ([]string, error) {
	siteURL, err := url.Parse(fmt.Sprintf("%s://%s", a.u.Scheme, a.u.Host))
	if err != nil {
		log.Fatalln(err)
	}

	return a.c.GetAllElementAttributeByQuery(a.u.String(), "li.fa-play-circle > a", "href", siteURL.String(), false)
}

func (a AnimeFlv2Downloader) DownloadEpisodes(episodes []string) {
	var preDownloadURLs []string
	for _, episode := range episodes {
		toAdd, err := a.c.GetElementAttributeByQuery(episode, "a.BtnNw-a", "href", "https:", false)
		if err != nil {
			log.Fatalln(err)
		}

		preDownloadURLs = append(preDownloadURLs, toAdd)
	}

	var wg sync.WaitGroup
	for _, preDownloadURL := range preDownloadURLs {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			var err error
			var name, downloadURL string

			name, err = a.c.GetElementTextByQuery(url, "#title", false)
			if err != nil {
				log.Fatalln(err)
			}

			downloadURL, err = a.c.GetElementAttributeByQuery(url, "#content-download > div:nth-child(1) > div:nth-child(3) > a", "href", "", true)
			if err != nil {
				log.Fatalln(err)
			}

			if name != "" {
				_ = a.s.CreateFileFromURL(name, downloadURL)
			}
		}(preDownloadURL)

		time.Sleep(500 * time.Millisecond)
	}

	wg.Wait()
}
