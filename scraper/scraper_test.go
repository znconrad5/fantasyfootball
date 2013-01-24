package scraper

import (
	"github.com/znconrad5/fantasyfootball"
	"testing"
)

func TestNflUrlGenerator(t *testing.T) {
	generator := &nflUrlGenerator{
		positions: []fantasyfootball.Position{fantasyfootball.QB},
		season: 2012,
		startWeek: 1,
		endWeek:31,
	}
	urls := generator.generateUrls()
	_, ok := urls["http://fantasy.nfl.com/research/scoringleaders?position=1&count=48&statCategory=stats&statSeason=2012&statType=weekStats&statWeek=1"]
	if !ok {
		t.Error()
	}
	_, ok = urls["http://fantasy.nfl.com/research/scoringleaders?position=1&count=48&statCategory=stats&statSeason=2012&statType=weekStats&statWeek=2"]
	if !ok {
		t.Error()
	}
	_, ok = urls["http://fantasy.nfl.com/research/scoringleaders?position=1&count=48&statCategory=stats&statSeason=2012&statType=weekStats&statWeek=3"]
	if !ok {
		t.Error()
	}
}
