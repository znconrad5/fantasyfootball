package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	urlPattern = "http://accuscore.com/fantasy-sports/nfl-fantasy-sports/%v-%v"
)

var (
	dataDirFlag   = flag.String("dataDir", "./html", "The directory to put the raw html in.")
	positionsFlag = flag.String("positions", "QB,RB,WR,TE,DEF-ST,K", "The comma separated positions to scrape, 'QB', 'RB', 'WR', 'TE', 'LB', 'LB', 'DB', 'DEF-ST', 'K', and/or 'P'")
	startWeekFlag = flag.Int("startWeek", 2, "The week to start player statistic gathering.")
	endWeekFlag   = flag.Int("endWeek", 14, "The week to end player statistic gathering, inclusive.")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	positions := strings.Split(*positionsFlag, ",")

	rateLimiter := time.NewTicker(4 * time.Second)
	var waitGroup sync.WaitGroup
	for _, v := range positions {
		for i := *startWeekFlag; i <= *endWeekFlag; i++ {
			week := fmt.Sprintf("Week-%v", i)
			url := fmt.Sprintf(urlPattern, week, v)
			waitGroup.Add(1)
			go asyncFetch(url, strconv.Itoa(i), v, rateLimiter.C, &waitGroup)
		}
		waitGroup.Add(1)
		go asyncFetch(fmt.Sprintf(urlPattern, "Current-Week", v), "curr", v, rateLimiter.C, &waitGroup)
	}
	waitGroup.Wait()
}

func asyncFetch(url string, id string, pos string, rateLimiter <-chan time.Time, waitGroup *sync.WaitGroup) {
	<-rateLimiter
	fetch(url, id, pos)
	waitGroup.Done()
}

func fetch(url string, id string, pos string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("encountered error crawling %v: %v\n", url, err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("encountered error reading %v: %v\n", url, err)
		return
	}
	filename := *dataDirFlag + "/" + pos + "_" + id + ".html"
	err = ioutil.WriteFile(filename, body, 0600)
	if err != nil {
		fmt.Printf("encountered error saving %v: %v\n", url, err)
		return
	}
}
