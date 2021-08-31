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
var urlExtractor *web.URLExtractor
var localStorage storage.Storage

func main() {
	wd, _ = os.Getwd()

	urlExtractor, localStorage = bootstrap.WarmUp(wd + "/Anime/" + os.Args[2])
	localStorage.Initialize()

	animeURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	episodes := urlExtractor.GetFromStaticWebsiteAttributeValueFromAllElementsByQuery(animeURL.String(), "li.fa-play-circle > a", "href", animeURL.Scheme+"://"+animeURL.Host)
	downloadEpisodes(episodes)
}

func downloadEpisodes(episodes []string) {
	var preDownloadURLs []string
	for _, episode := range episodes {
		preDownloadURLs = append(preDownloadURLs, urlExtractor.GetFromStaticWebsiteAttributeValueFromAllElementsByQuery(episode, "a.BtnNw-a", "href", "https:")...)
	}

	var wg sync.WaitGroup
	for _, preDownloadURL := range preDownloadURLs {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			name, downloadURL := urlExtractor.GetFromDynamicWebsiteAttributeValueFromElementByQuery(url, "#content-download > div:nth-child(1) > div:nth-child(3) > a", "href")
			fmt.Printf("Downloading episode %s on url %s \n", name, downloadURL)
			err = localStorage.CreateFileFromURL(name, downloadURL)
			if err != nil {
				log.Fatalln(err)
			}
		}(preDownloadURL)
	}

	wg.Wait()
}
