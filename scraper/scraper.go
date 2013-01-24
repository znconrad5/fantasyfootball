package scraper

import (
	"fmt"
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
			pages <- &ScrapedPage{
				Content:     fetch(urlToFetch).Body,
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

func NewNflScraper(positions []fantasyfootball.Position, season int, startWeek int, endWeek int) *Scraper {
	return &Scraper{
		&nflUrlGenerator{
			positions: positions,
			season:    season,
			startWeek: startWeek,
			endWeek:   endWeek,
		},
		time.NewTicker(2 * time.Second),
	}
}

const nflUrlFormat string = "http://fantasy.nfl.com/research/scoringleaders?position=%d&count=%d&statCategory=stats&statSeason=%d&statType=weekStats&statWeek=%d"

type nflPositionData struct {
	queryValue int
	count      int
}

var nflPositionMapping = map[fantasyfootball.Position]nflPositionData{
	fantasyfootball.QB:     nflPositionData{1, 48},
	fantasyfootball.RB:     nflPositionData{2, 96},
	fantasyfootball.WR:     nflPositionData{3, 128},
	fantasyfootball.TE:     nflPositionData{4, 64},
	fantasyfootball.K:      nflPositionData{7, 48},
	fantasyfootball.DST: 	nflPositionData{8, 32},
}

type nflUrlGenerator struct {
	positions []fantasyfootball.Position
	season    int
	startWeek int
	endWeek   int
}

func (generator *nflUrlGenerator) generateUrls() map[string]string {
	urlMap := make(map[string]string)
	for _, position := range generator.positions {
		for week := generator.startWeek; week <= generator.endWeek; week++ {
			nflPositionMapping := nflPositionMapping[position]
			url := fmt.Sprintf(nflUrlFormat, nflPositionMapping.queryValue, nflPositionMapping.count, generator.season, week)
			desc := fmt.Sprintf("NFL_%v_%v_%v", generator.season, position, week)
			urlMap[url] = desc
		}
	}
	return urlMap
}