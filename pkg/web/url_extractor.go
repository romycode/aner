package web

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
	"github.com/romycode/anime-downloader/pkg/errors"
)

type URLExtractor struct {
	eh *errors.ErrorHandler
	c  *colly.Collector
}

func NewURLExtractor(eh *errors.ErrorHandler) *URLExtractor {
	return &URLExtractor{
		eh: eh,
		c: colly.NewCollector(
			colly.Async(true),
		),
	}
}

func (u URLExtractor) GetFromStaticWebsiteAttributeValueFromAllElementsByQuery(targetURL string, sel string, attr string, prefix string) []string {
	var urls []string
	u.c.OnHTML(sel, func(e *colly.HTMLElement) {
		if prefix != "" {
			urls = append(urls, prefix+e.Attr(attr))
		}
	})

	err := u.c.Visit(targetURL)
	if err != nil {
		u.eh.HandleError(err)
	}

	u.c.Wait()

	return urls
}

func (u URLExtractor) GetFromStaticWebsiteAttributeValueFromElementByQuery(targetURL string, sel string, attr string, prefix string) string {
	var downloadURL string
	u.c.OnHTML(sel, func(e *colly.HTMLElement) {
		if prefix != "" {
			downloadURL = prefix + e.Attr(attr)
		}
	})

	err := u.c.Visit(targetURL)
	if err != nil {
		u.eh.HandleError(err)
	}

	u.c.Wait()

	return downloadURL
}

func (u URLExtractor) GetFromDynamicWebsiteAttributeValueFromElementByQuery(siteURL string, sel string, attr string) (string, string) {
	parsedURL, _ := url.Parse(siteURL)

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

	var ok bool
	var name string
	var downloadURL string

	_ = chromedp.Run(ctx,
		chromedp.Navigate(parsedURL.String()),
		chromedp.Text(
			"#title",
			&name,
			chromedp.ByQuery,
		),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.AttributeValue(
			sel,
			attr,
			&downloadURL,
			&ok,
			chromedp.ByQuery,
		),
	)

	return name, downloadURL
}
