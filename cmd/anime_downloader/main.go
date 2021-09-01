package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/romycode/anime-downloader/cmd/anime_downloader/bootstrap"
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

var err error
var wd string
var urlExtractor *web.Crawler
var localStorage storage.Storage

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	wd, _ = os.Getwd()
	animeURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	siteURL, err := url.Parse(fmt.Sprintf("%s://%s", animeURL.Scheme, animeURL.Host))
	if err != nil {
		log.Fatalln(err)
	}

	urlExtractor, localStorage = bootstrap.WarmUp(fmt.Sprintf("%s/Anime/%s", wd, os.Args[2]))
	localStorage.Initialize()

	episodes, err := urlExtractor.GetAllElementAttributeByQuery(animeURL.String(), "li.fa-play-circle > a", "href", siteURL.String(), false)
	if err != nil {
		log.Fatalln(err)
	}

	downloadEpisodes(episodes)
}

func downloadEpisodes(episodes []string) {
	var preDownloadURLs []string
	for _, episode := range episodes {
		toAdd, err := urlExtractor.GetElementAttributeByQuery(episode, "a.BtnNw-a", "href", "https:", false)
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

			name, err := urlExtractor.GetElementTextByQuery(url, "#title", false)
			if err != nil {
				log.Fatalln(err)
			}

			downloadURL, err := urlExtractor.GetElementAttributeByQuery(url, "#content-download > div:nth-child(1) > div:nth-child(3) > a", "href", "", true)
			if err != nil {
				log.Fatalln(err)
			}



			fmt.Printf("Downloading episode %s on url %s \n", name, downloadURL)
			err = localStorage.CreateFileFromURL(name, downloadURL)
			if err != nil {
				log.Fatalln(err)
			}
		}(preDownloadURL)
	}

	wg.Wait()
}
