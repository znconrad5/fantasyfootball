package scraper

import (
	"github.com/znconrad5/fantasyfootball"
	"io"
	"net/http"
	"sync"
	"time"
)

type urlGenerator interface {
	generateUrls() map[string]string
}

type ScrapedPage struct {
	Content     io.ReadCloser
	Description string
}

type Scraper struct {
	urlGenerator
	rateLimiter *time.Ticker
}

func (scraper *Scraper) Scrape() <-chan *ScrapedPage {
	pages := make(chan *ScrapedPage)
	var waitGroup sync.WaitGroup
	for url, desc := range scraper.generateUrls() {
		waitGroup.Add(1)
		go func(urlToFetch, descOfUrl string) {
			defer waitGroup.Done()
			<-scraper.rateLimiter.C
			response := fetch(urlToFetch)
			pages <- &ScrapedPage{
				Content:     response.Body,
				Description: descOfUrl,
			}
		}(url, desc)
	}
	go func() {
		waitGroup.Wait()
		close(pages)
	}()
	return pages
}

func fetch(url string) *http.Response {
	res, err := http.Get(url)
	fantasyfootball.HandleError(err)
	return res
}