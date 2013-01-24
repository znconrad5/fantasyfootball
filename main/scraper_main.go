package main

import (
	"flag"
	"github.com/znconrad5/fantasyfootball"
	fscraper "github.com/znconrad5/fantasyfootball/scraper"
	"io"
	"os"
	"runtime"
)

var (
	siteFlag      *string = flag.String("site", "", "The site to get data from.")
	dataDirFlag   *string = flag.String("dataDir", os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/html"), "The directory to put the raw html in.")
	positionsFlag *string = flag.String("positions", "QB,RB,WR,TE,DEF-ST,K", "The comma separated positions to scrape, 'QB', 'RB', 'WR', 'TE', 'LB', 'LB', 'DB', 'DEF-ST', 'K', and/or 'P'.")
	seasonFlag    *int    = flag.Int("season", 2012, "The year to gather player statistics")
	startWeekFlag *int    = flag.Int("startWeek", 1, "The week to start player statistic gathering.")
	endWeekFlag   *int    = flag.Int("endWeek", 17, "The week to end player statistic gathering, inclusive.")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	positions, err := fantasyfootball.ParsePositions(*positionsFlag, ",")
	fantasyfootball.HandleError(err)
	if len(positions) == 0 {
		panic("No positions specified")
	}
	var scraper *fscraper.Scraper
	switch *siteFlag {
	case "nfl":
		scraper = fscraper.NewNflScraper(positions, *seasonFlag, *startWeekFlag, *endWeekFlag)
	default:
		panic("Unexpected flag for --site: " + *siteFlag)
	}
	for page := range scraper.Scrape() {
		var content io.ReadCloser = page.Content
		defer content.Close()
		
		file, err := os.Create(*dataDirFlag + "/" + page.Description + ".html")
		fantasyfootball.HandleError(err)
		defer func() { fantasyfootball.HandleError(file.Close()) }()
		
		io.Copy(file, content)
	}
}