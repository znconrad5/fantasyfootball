package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var dataDir = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/html")

func main() {
	urlPattern := "http://accuscore.com/fantasy-sports/nfl-fantasy-sports/%v-%v"
	positions := []string{"QB", "RB", "WR", "TE", "LB", "DL", "DB", "DEF-ST", "K", "P"}

	var waitGroup sync.WaitGroup

	var asyncFetch = func(urlf1 string, idf1 string, posf1 string, wgf1 *sync.WaitGroup) {
		fetch(urlf1, idf1, posf1)
		wgf1.Done()
	}

	for _, v := range positions {
		for i := 1; i <= 17; i++ {
			week := fmt.Sprintf("Week-%v", i)
			url := fmt.Sprintf(urlPattern, week, v)
			waitGroup.Add(1)
			go asyncFetch(url, strconv.Itoa(i), v, &waitGroup)
		}
		waitGroup.Add(1)
		go asyncFetch(fmt.Sprintf(urlPattern, "Current-Week", v), "curr", v, &waitGroup)
	}
	waitGroup.Wait()
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
	filename := dataDir + "/" + pos + "_" + id + ".html"
	err = ioutil.WriteFile(filename, body, 0600)
	if err != nil {
		fmt.Printf("encountered error saving %v: %v\n", url, err)
		return
	}
}
