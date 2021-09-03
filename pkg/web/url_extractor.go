package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Crawler struct {
	htmlJsCache map[string]string
	htmlCache   map[string]string
	mu          sync.Mutex
}

var AttributeNotFound = errors.New("attribute not found")

func NewCrawler() *Crawler {
	return &Crawler{
		htmlJsCache: map[string]string{},
		htmlCache:   map[string]string{},
		mu:          sync.Mutex{},
	}
}

func (u *Crawler) GetAllElementAttributeByQuery(targetURL string, sel string, attr string, prefix string, needsJs bool) ([]string, error) {
	document := u.fetchURL(targetURL, sel, needsJs)

	html, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return []string{}, err
	}

	var urls []string
	html.Find(sel).Each(func(_ int, selection *goquery.Selection) {
		value, ok := selection.Attr(attr)
		if !ok {
			log.Fatalln(AttributeNotFound)
		}

		urls = append(urls, fmt.Sprintf("%s%s", prefix, value))
	})

	return urls, nil
}

func (u *Crawler) GetElementAttributeByQuery(targetURL string, sel string, attr string, prefix string, needsJs bool) (string, error) {
	urls, err := u.GetAllElementAttributeByQuery(targetURL, sel, attr, prefix, needsJs)
	if err != nil {
		return "", err
	}
	if len(urls) > 1 {
		return "", errors.New(fmt.Sprintf("foound multiple elements for selector: %s", sel))
	}
	if len(urls) == 0 {
		urls, _ = u.GetAllElementAttributeByQuery(targetURL, sel, attr, prefix, needsJs)
	}

	return urls[0], nil
}

func (u *Crawler) GetElementTextByQuery(targetURL string, sel string, needsJs bool) (string, error) {
	document := u.fetchURL(targetURL, sel, needsJs)

	html, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return "", err
	}

	var text string
	html.Find(sel).Each(func(_ int, selection *goquery.Selection) {
		text = selection.Text()
	})
	return text, nil
}

func (u *Crawler) fetchURL(targetURL string, sel string, needJs bool) string {
	var document string
	if needJs {
		if html, ok := u.htmlJsCache[targetURL]; ok {
			return html
		}

		opts := []chromedp.ExecAllocatorOption{
			chromedp.Headless,
			chromedp.NoFirstRun,
			chromedp.NoSandbox,
			chromedp.DisableGPU,
			chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"),
		}
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		defer cancel()
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		_ = chromedp.Run(ctx,
			chromedp.Navigate(targetURL),
			chromedp.Sleep(500*time.Millisecond),
			chromedp.WaitVisible(sel, chromedp.ByQuery),
			chromedp.InnerHTML(
				"html",
				&document,
				chromedp.ByQuery,
			),
		)

		u.mu.Lock()
		u.htmlJsCache[targetURL] = document
		u.mu.Unlock()
	} else {
		if html, ok := u.htmlCache[targetURL]; ok {
			return html
		}
		if html, ok := u.htmlJsCache[targetURL]; ok {
			return html
		}

		res, err := http.Get(targetURL)
		if err != nil {
			log.Fatal(err)
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		document = buf.String()

		u.mu.Lock()
		u.htmlCache[targetURL] = document
		u.mu.Unlock()
	}
	return document
}
