package donwloader

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

type AnimeFlv3Downloader struct {
	u url.URL
	c *web.Crawler
	s storage.Storage
}

func NewAnimeFlv3Downloader(u url.URL, c *web.Crawler, s storage.Storage) *AnimeFlv3Downloader {
	return &AnimeFlv3Downloader{
		u: u,
		c: c,
		s: s,
	}
}

func (a AnimeFlv3Downloader) GetEpisodes() ([]string, error) {
	siteURL, err := url.Parse(fmt.Sprintf("%s://%s", a.u.Scheme, a.u.Host))
	if err != nil {
		log.Fatalln(err)
	}
	return a.c.GetAllElementAttributeByQuery(a.u.String(), "li.fa-play-circle > a", "href", siteURL.String(), true)
}

func (a AnimeFlv3Downloader) DownloadEpisodes(episodes []string) {
	preDownloadURLs := make(map[string]string)
	for _, episode := range episodes {
		downloadURLs, err := a.c.GetAllElementAttributeByQuery(episode, "a.Button.Sm.fa-download", "href", "", true)
		if err != nil {
			log.Fatalln(err)
		}

		name, err := a.c.GetElementTextByQuery(episode, "#XpndCn > div.CpCnA > div.CapiTop > h2", false)
		if err != nil {
			log.Fatalln(err)
		}

		for _, downloadURL := range downloadURLs {
			if strings.Contains(downloadURL, "streamtape.com") {
				preDownloadURLs[name] = downloadURL
				break
			}
		}
	}

	var wg sync.WaitGroup
	for name, preDownloadURL := range preDownloadURLs {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			iframeWithURL, err := a.c.GetElementTextByQuery(url, "#content > div > div:nth-child(6) > div > textarea", false)
			if err != nil {
				log.Fatalln(err)
			}
			iframe, err := goquery.NewDocumentFromReader(strings.NewReader(iframeWithURL))
			if err != nil {
				log.Fatalln(err)
			}

			downloadURL, _ := iframe.Find("iframe").Attr("src")

			if name != "" {
				_ = a.s.CreateFileFromURL(name, downloadURL)
			}
		}(preDownloadURL)
	}

	wg.Wait()
}
