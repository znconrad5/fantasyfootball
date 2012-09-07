package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
)

var dataDir = "../../html"

func main() {
	urlPattern := "http://accuscore.com/fantasy-sports/nfl-fantasy-sports/%v-%v"
	positions := []string{"QB", "RB", "WR", "TE", "LB", "DL", "DB", "DEF-ST", "K", "P"}

	for _, v := range positions {
		for i := 1; i <= 17; i++ {
			week := fmt.Sprintf("Week-%v", i)
			url := fmt.Sprintf(urlPattern, week, v)
			fetch(url, strconv.Itoa(i), v)
		}
		fetch(fmt.Sprintf(urlPattern, "Current-Week", v), "curr", v)
	}

}

func fetch(url string, id string, pos string) {
	fmt.Printf("crawling %v...\n", url)
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