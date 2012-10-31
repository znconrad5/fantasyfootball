package main

import (
	"flag"
	"fmt"
	"github.com/znconrad5/fantasyfootball"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	dataDirFlag   = flag.String("dataDir", os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/html"), "The directory to put the raw html in.")
	positionsFlag = flag.String("positions", "QB,RB,WR,TE,DEF-ST,K", "The comma separated positions to scrape, 'QB', 'RB', 'WR', 'TE', 'LB', 'LB', 'DB', 'DEF-ST', 'K', and/or 'P'")
	startWeekFlag = flag.Int("startWeek", 6, "The week to start player statistic gathering.")
	endWeekFlag   = flag.Int("endWeek", 14, "The week to end player statistic gathering, inclusive.")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	positions := strings.Split(*positionsFlag, ",")
	scraper := NewAccuscoreScraper(*startWeekFlag, *endWeekFlag, positions, 4*time.Second)
	contentChannel := scraper.Start()
	persister := &FilePersister{
		dataDir: *dataDirFlag,
		keyGen: func(data *ScraperOutput) string {
			return data.pairs["position"] + "_" + data.pairs["week"] + ".html"
		},
	}
	for content := range contentChannel {
		defer content.Close()
		persister.Persist(content)
	}
}

type ScraperOutput struct {
	content io.ReadCloser
	source  string
	pairs   map[string]string
}

func (scraperOutput *ScraperOutput) Read(p []byte) (n int, err error) {
	return scraperOutput.content.Read(p)
}

func (scraperOutput *ScraperOutput) Close() error {
	return scraperOutput.content.Close()
}

type Scraper interface {
	Start() <-chan *ScraperOutput
}

type AccuscoreScraper struct {
	startWeek   int
	endWeek     int
	positions   []string
	rateLimiter *time.Ticker
	urlPattern  string
}

const (
	urlPattern = "http://accuscore.com/fantasy-sports/nfl-fantasy-sports/%v-%v"
)

func NewAccuscoreScraper(startWeek int, endWeek int, positions []string, period time.Duration) *AccuscoreScraper {
	accuscoreScraper := &AccuscoreScraper{
		startWeek: startWeek,
		endWeek:   endWeek,
		positions: positions,
	}
	accuscoreScraper.rateLimiter = time.NewTicker(period)
	accuscoreScraper.urlPattern = urlPattern
	return accuscoreScraper
}

func (scraper *AccuscoreScraper) Start() <-chan *ScraperOutput {
	content := make(chan *ScraperOutput)
	go func(out chan<- *ScraperOutput) {
		var waitGroup sync.WaitGroup
		for _, v := range scraper.positions {
			for i := scraper.startWeek; i <= scraper.endWeek; i++ {
				week := fmt.Sprintf("Week-%v", i)
				url := fmt.Sprintf(scraper.urlPattern, week, v)
				waitGroup.Add(1)
				go scraper.asyncFetch(url, strconv.Itoa(i), v, &waitGroup, out)
			}
			waitGroup.Add(1)
			go scraper.asyncFetch(fmt.Sprintf(scraper.urlPattern, "Current-Week", v), "curr", v, &waitGroup, out)
		}
		waitGroup.Wait()
		close(out)
	}(content)
	return content
}

func (scraper *AccuscoreScraper) asyncFetch(url string, id string, pos string, waitGroup *sync.WaitGroup, out chan<- *ScraperOutput) {
	defer waitGroup.Done()
	<-scraper.rateLimiter.C
	out <- scraper.fetch(url, id, pos)
}

func (scraper *AccuscoreScraper) fetch(url string, id string, pos string) *ScraperOutput {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("encountered error crawling %v: %v\n", url, err)
		return nil
	}
	return &ScraperOutput{
		content: res.Body,
		source:  url,
		pairs: map[string]string{
			"week":     id,
			"position": pos,
		},
	}
}

type Persister interface {
	Persist(key string, data *ScraperOutput)
}

type FilePersister struct {
	dataDir string
	keyGen  func(data *ScraperOutput) string
}

func (persister *FilePersister) Persist(data *ScraperOutput) {
	file, err := os.Create(persister.dataDir + "/" + persister.keyGen(data))
	fantasyfootball.HandleError(err)
	defer file.Close()
	io.Copy(file, data)
}
