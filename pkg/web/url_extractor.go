package web

import (
	"net/url"

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
			colly.MaxDepth(500),
			colly.Async(true),
		),
	}
}

func (u URLExtractor) GetAttributeValueFromElementByQuery(targetURL string, sel string, attr string) []string {
	parsedURL, _ := url.Parse(targetURL)

	var urls []string
	u.c.OnHTML(sel, func(e *colly.HTMLElement) {
		urls = append(urls, parsedURL.Scheme + "://" + parsedURL.Host + e.Attr(attr))
	})

	err := u.c.Visit(targetURL)
	u.eh.HandleError(err)

	u.c.Wait()
	return urls
}
