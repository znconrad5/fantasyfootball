package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/znconrad5/fantasyfootball"
	"io"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"sync"
)

const (
	CurrWeekString = "curr"
)

var (
	inputDir      = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/html")
	testOutDir    = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed")
	rowRegex      = regexp.MustCompile("(?s)<tr[^>]*>.*?</tr>")
	dataRegex     = regexp.MustCompile(">[\\s\\r\\n]*([^<>]*?\\w+[^<>]*?)[\\s\\r\\n]*<")
	posRegex      = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"pos\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
	weekRegex     = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"split\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
	parseWeek     = regexp.MustCompile("^Week-(\\d+)$")
	currWeekRegex = regexp.MustCompile(CurrWeekString)
)

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
			newName := fmt.Sprintf("%s/%s", testOutDir, currWeekRegex.ReplaceAllString(fileName, fmt.Sprintf("%v", minWeek-1)))
			os.Remove(newName)
			os.Rename(fmt.Sprintf("%s/%s", testOutDir, fileName), newName)
		}
	}
}

func parseFile(in string, out string) int {
	content, err := os.Open(in)
	fantasyfootball.HandleError(err)
	defer content.Close()

	//create outputfile, don't know exactly what to call it yet
	file, err := ioutil.TempFile(testOutDir, "parsed_")
	fantasyfootball.HandleError(err)

	//parse the file
	week, pos := parse(content, file)
	defer os.Remove(file.Name())
	defer file.Close()

	//name the file based off it's contents
	var weekString string
	if week == -1 {
		weekString = CurrWeekString
	} else {
		weekString = fmt.Sprintf("%v", week)
	}
	actualName := fmt.Sprintf("%s/%s_%v.txt", testOutDir, pos, weekString)
	newFile, err := os.Create(actualName)
	fantasyfootball.HandleError(err)
	defer newFile.Close()
	file.Seek(0, 0)
	_, err = io.Copy(newFile, file)
	fantasyfootball.HandleError(err)

	return week
}

func parse(in io.Reader, out io.Writer) (week int, pos string) {
	content, err := ioutil.ReadAll(in)
	fantasyfootball.HandleError(err)

	//extract position from page
	posMatch := posRegex.FindSubmatch(content)
	if posMatch == nil {
		fantasyfootball.HandleError(errors.New(fmt.Sprintf("unable to parse position from file: %s", in)))
	}
	pos = string(posMatch[1])

	//extract week from page
	weekMatch := weekRegex.FindSubmatch(content)
	if weekMatch == nil {
		fantasyfootball.HandleError(errors.New(fmt.Sprintf("unable to parse week from file: %s", in)))
	}
	var weekString string
	if parseWeekMatch := parseWeek.FindSubmatch(weekMatch[1]); parseWeekMatch != nil {
		weekString = fmt.Sprintf("%s", parseWeekMatch[1])
	} else {
		weekString = CurrWeekString
	}

	//prepare csv output
	csvWriter := csv.NewWriter(out)
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
	weekInt, err := strconv.ParseInt(weekString, 10, 0)
	if err != nil {
		weekInt = -1
	}
	return int(weekInt), pos
}
