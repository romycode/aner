package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/romycode/anime-downloader/cmd/anime_downloader/bootstrap"
	"github.com/romycode/anime-downloader/pkg/downloader"
	"github.com/romycode/anime-downloader/pkg/storage"
	"github.com/romycode/anime-downloader/pkg/web"
)

var wd string
var crawler *web.Crawler
var localStorage storage.Storage

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	wd, _ = os.Getwd()
	animeURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	crawler, localStorage = bootstrap.WarmUp(fmt.Sprintf("%s/Anime/%s", wd, os.Args[2]))
	localStorage.Initialize()

	downloader, err := donwloader.NewAnimeDownloaderFactory(crawler, localStorage).Build(*animeURL)
	if err != nil {
		log.Fatalln(err)
	}

	episodes, err := downloader.GetEpisodes()
	if err != nil {
		log.Fatalln(err)
	}

	downloader.DownloadEpisodes(episodes)
}
