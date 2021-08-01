package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

type (
	EpisodeInformation struct {
		EpisodeName string
		DownloadURL string
	}
)

var (
	logger         = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	animeflvURL, _ = url.Parse("https://www2.animeflv.to")
	animePath      = ""
)

func main() {
	dir, _ := os.Getwd()

	animePath = dir + "/Anime/" + os.Args[2] + "/"
	err := os.MkdirAll(animePath, 0755)
	handleError(err)

	animeURL, err := url.Parse(os.Args[1])
	handleError(err)

	episodeURLs := fetchEpisodesURLs(animeURL)
	downloadEpisodes(episodeURLs)
}

func fetchEpisodesURLs(animeURL *url.URL) []string {
	var episodeURLs []string

	c := colly.NewCollector(
		colly.MaxDepth(500),
		colly.Async(true),
	)
	// Find and visit all links
	c.OnHTML("li", func(li *colly.HTMLElement) {
		if li.Attr("class") == "fa-play-circle" {
			episodeURLs = append(episodeURLs, animeflvURL.String()+li.ChildAttr("a", "href"))
		}
	})

	err := c.Visit(animeURL.String())
	handleError(err)
	c.Wait()

	return episodeURLs
}

func downloadEpisodes(episodeURLs []string) {
	var tmpURLs []string

	c := colly.NewCollector(
		colly.MaxDepth(500),
		colly.Async(true),
	)
	c.OnHTML("a", func(a *colly.HTMLElement) {
		if a.Attr("class") == "BtnNw-a" {
			tmpURLs = append(tmpURLs, "https:"+a.Attr("href"))
		}
	})

	for _, episodeURL := range episodeURLs {
		err := c.Visit(episodeURL)
		handleError(err)
	}
	c.Wait()

	var wg sync.WaitGroup
	for _, tmpURL := range tmpURLs {
		wg.Add(1)
		episodeData := getURL(tmpURL)
		go downloadEpisode(episodeData.EpisodeName, episodeData.DownloadURL, &wg)
	}
}

func getURL(tmpURL string) EpisodeInformation {
	requestURL, _ := url.Parse(tmpURL)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Headless,
		chromedp.NoFirstRun,
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var ok bool
	var directDownloadURL string
	err := chromedp.Run(ctx,
		chromedp.Navigate(requestURL.String()),
		chromedp.WaitVisible(`#content-download > div:nth-child(1) > div:nth-child(3) > a`, chromedp.ByQuery),
		chromedp.AttributeValue(
			"#content-download > div:nth-child(1) > div:nth-child(3) > a",
			"href",
			&directDownloadURL,
			&ok,
			chromedp.ByQuery,
		),
	)
	handleError(err)

	episodeName := requestURL.Query().Get("title")
	return EpisodeInformation{
		EpisodeName: episodeName,
		DownloadURL: directDownloadURL,
	}
}

func downloadEpisode(episodeName string, downloadURL string, wg *sync.WaitGroup) {
	logger.Printf("Downloading: %s , from: %s\n", episodeName, downloadURL)

	out, err := os.Create(animePath + episodeName + ".mp4")
	handleError(err)
	defer func(out *os.File) {
		err := out.Close()
		handleError(err)
	}(out)

	resp, err := http.Get(downloadURL)
	handleError(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	handleError(err)

	wg.Done()
}

func handleError(err error) {
	if err != nil {
		logger.Fatalln(err)
	}
}
