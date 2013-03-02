package scraper

import (
	"fmt"
	"github.com/znconrad5/fantasyfootball"
	"time"
)

func NewAccuscoreScraper(positions []fantasyfootball.Position, startWeek int, endWeek int) *Scraper {
	return &Scraper{
		&accuscoreUrlGenerator{
			positions: positions,
			startWeek: startWeek,
			endWeek:   endWeek,
		},
		time.NewTicker(2 * time.Second),
	}
}

const accuscoreUrlFormat string = "http://accuscore.com/fantasy-sports/nfl-fantasy-sports/%v-%s"

type accuscoreUrlGenerator struct {
	positions []fantasyfootball.Position
	startWeek int
	endWeek   int
}

func (generator *accuscoreUrlGenerator) generateUrls() map[string]string {
	urlMap := make(map[string]string)
	for _, position := range generator.positions {
		for week := generator.startWeek; week <= generator.endWeek; week++ {
			url := fmt.Sprintf(accuscoreUrlFormat, fmt.Sprintf("Week-%d", week), position)
			desc := fmt.Sprintf("ACCUSCORE_%s_%d", position, week)
			urlMap[url] = desc
		}
		url := fmt.Sprintf(accuscoreUrlFormat, "Current-Week", position)
		desc := fmt.Sprintf("ACCUSCORE_%s_%s", position, "curr")
		urlMap[url] = desc
	}
	return urlMap
}