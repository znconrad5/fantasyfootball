package main

import (
	"encoding/csv"
	"fmt"
	"github.com/znconrad5/fantasyfootball"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var inputDir = "./html"
var testOutDir = "./parsed"

var rowRegex = regexp.MustCompile("(?s)<tr[^>]*>.*?</tr>")
var dataRegex = regexp.MustCompile(">[\\s\\r\\n]*([^<>]*?\\w+[^<>]*?)[\\s\\r\\n]*<")
var posRegex = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"pos\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
var weekRegex = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"split\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
var parseWeek = regexp.MustCompile("^Week-(\\d+)$")
var currWeekString = "curr"
var currWeekRegex = regexp.MustCompile(currWeekString)

func main() {
	files, err := ioutil.ReadDir(inputDir)
	fantasyfootball.HandleError(err)

	var asyncParse = func(in string, out string, wg *sync.WaitGroup, weeks chan<- int) {
		weeks <- parseFile(in, out)
		wg.Done()
	}

	var waitGroup sync.WaitGroup
	weekChan := make(chan int)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		waitGroup.Add(1)
		go asyncParse(fmt.Sprintf("%s/%s", inputDir, file.Name()), testOutDir, &waitGroup, weekChan)
	}

	minWeekChan := make(chan int, 1)
	go func(weeks <-chan int, out chan<- int) {
		min := math.MaxInt16
		for week := range weeks {
			if week > 0 && week < min {
				min = week
			}
		}
		out <- min
	}(weekChan, minWeekChan)
	waitGroup.Wait()
	close(weekChan)

	//fix weeks name
	minWeek := <-minWeekChan
	outFiles, err := ioutil.ReadDir(testOutDir)
	fantasyfootball.HandleError(err)
	for _, file := range outFiles {
		if file.IsDir() {
			continue
		}
		if fileName := file.Name(); currWeekRegex.MatchString(fileName) {
			os.Rename(fmt.Sprintf("%s/%s", testOutDir, fileName), fmt.Sprintf("%s/%s", testOutDir, currWeekRegex.ReplaceAllString(fileName, fmt.Sprintf("%v", minWeek-1))))
		}
	}
}

func parseFile(in string, out string) int {
	content, err := ioutil.ReadFile(in)
	fantasyfootball.HandleError(err)

	//extract position from page
	posMatch := posRegex.FindSubmatch(content)
	if posMatch == nil {
		log.Fatalf("unable to parse position from file: %s", in)
	}
	pos := posMatch[1]

	//extract week from page
	weekMatch := weekRegex.FindSubmatch(content)
	if weekMatch == nil {
		log.Fatalf("unable to parse week from file: %s", in)
	}
	var week string
	if parseWeekMatch := parseWeek.FindSubmatch(weekMatch[1]); parseWeekMatch != nil {
		week = fmt.Sprintf("%s", parseWeekMatch[1])
	} else {
		week = currWeekString
	}

	//prepare csv file
	file, err := os.Create(fmt.Sprintf("%s/%s_%v.txt", testOutDir, pos, week))
	fantasyfootball.HandleError(err)
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	csvWriter.Comma = '\t'

	//extract player info from page
	rows := rowRegex.FindAll(content, -1)
	for _, row := range rows {
		columns := dataRegex.FindAllSubmatch(row, -1)
		csvRow := make([]string, len(columns))
		for i, column := range columns {
			csvRow[i] = fmt.Sprintf("%s", column[1])
		}
		csvWriter.Write(csvRow)
	}
	weekInt, err := strconv.ParseInt(week, 10, 0)
	if err != nil {
		weekInt = -1
	}
	return int(weekInt)
}
