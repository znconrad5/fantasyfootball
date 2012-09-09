package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var testFile = "../html/DB_2.html"
var testOutDir = "../testout"

var rowRegex = regexp.MustCompile("(?s)<tr[^>]*>.*?</tr>")
var dataRegex = regexp.MustCompile(">[\\s\\r\\n]*([^<>]*?\\w+[^<>]*?)[\\s\\r\\n]*<")
var posRegex = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"pos\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
var weekRegex = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"split\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
var parseWeek = regexp.MustCompile("^Week-(\\d+)$")

func main() {
	parseFile(testFile, testOutDir)
}

func parseFile(in string, out string) {
	content, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Printf("encountered error opening file for read: %v", err)
		return
	}

	//extract position from page
	posMatch := posRegex.FindSubmatch(content)
	if posMatch == nil {
		fmt.Printf("unable to find position")
		return
	}
	pos := posMatch[1]

	//extract week from page
	weekMatch := weekRegex.FindSubmatch(content)
	if weekMatch == nil {
		fmt.Printf("unable to find week")
		return
	}
	var week string
	if parseWeekMatch := parseWeek.FindSubmatch(weekMatch[1]); parseWeekMatch != nil {
		week = fmt.Sprintf("%s", parseWeekMatch[1])
	} else {
		week = "curr"
	}

	//prepare csv file
	file, err := os.Create(fmt.Sprintf("%s/%s_%v.csv", testOutDir, pos, week))
	if err != nil {
		fmt.Printf("encountered error opening file: %v", err)
		return
	}
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
}
