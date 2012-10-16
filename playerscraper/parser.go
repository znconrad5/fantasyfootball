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
	CurrWeekString      = "curr"
	CurrWeekPlaceholder = -1
)

var (
	inputDir         = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/html")
	testOutDir       = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed")
	headerRegex      = regexp.MustCompile("(?s)<thead>[\\r\\n\\s]*<tr[^>]*>.*?</tr>[\\r\\n\\s]*</thead>")
	rowRegex         = regexp.MustCompile("(?s)<tr[^>]*class=\"[^\"]*\"[^>]*>.*?</tr>")
	dataRegex        = regexp.MustCompile(">[\\s\\r\\n]*([^<>]*?\\w+[^<>]*?)[\\s\\r\\n]*<")
	posRegex         = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"pos\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
	weekRegex        = regexp.MustCompile("(?s)<select\\s+[^>]*name=\"split\"[^>]*>.*?<option\\s+selected\\s+value=\"([^\"]+)\">.*?</select>")
	parseWeek        = regexp.MustCompile("^Week-(\\d+)$")
	currWeekRegex    = regexp.MustCompile(CurrWeekString)
	placeholderRegex = regexp.MustCompile(fmt.Sprintf("%v", CurrWeekPlaceholder))
	name             = "PLAYER"
	team             = "TEAM"
)

type PlayerRecord struct {
	Week     int
	Name     string
	Team     string
	Position fantasyfootball.Position
	Stats    map[string]float64
}

type Parser interface {
	Parse(content <-chan io.ReadCloser) <-chan PlayerRecord
}

type AccuscoreParser struct {
	row            *regexp.Regexp
	header         *regexp.Regexp
	data           *regexp.Regexp
	position       *regexp.Regexp
	week           *regexp.Regexp
	parseWeek      *regexp.Regexp
	currWeekString string
	nameString     string
	teamString     string
}

func NewAccuscoreParser() *AccuscoreParser {
	return &AccuscoreParser{
		row:            rowRegex,
		header:         headerRegex,
		data:           dataRegex,
		position:       posRegex,
		week:           weekRegex,
		parseWeek:      parseWeek,
		currWeekString: CurrWeekString,
		nameString:     name,
		teamString:     team,
	}
}

/**func main() {
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
}**/

type ParseOutput struct {
	writer  *csv.Writer
	headers []string
}

func main() {
	//create channel of content from directory
	content := make(chan io.ReadCloser)
	go func(newContent chan<- io.ReadCloser) {
		files, err := ioutil.ReadDir(inputDir)
		fantasyfootball.HandleError(err)
		var channelContentWg sync.WaitGroup
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			channelContentWg.Add(1)
			go func(fileName string, addContent chan<- io.ReadCloser, wg *sync.WaitGroup) {
				fileStream, err := os.Open(fileName)
				fantasyfootball.HandleError(err)
				addContent <- fileStream
				wg.Done()
			}(fmt.Sprintf("%v/%v", inputDir, file.Name()), newContent, &channelContentWg)
		}
		channelContentWg.Wait()
		close(newContent)
	}(content)

	parser := NewAccuscoreParser()
	playerRecordsChannel := parser.Parse(content)

	fileMap := make(map[string]*ParseOutput)
	minWeekResult :=
		func(playerRecords <-chan PlayerRecord) int {
			var minWeek int = math.MaxInt32
			for record := range playerRecords {
				if record.Week > 0 && record.Week < minWeek {
					minWeek = record.Week
				}
				fileName := fmt.Sprintf("%v_%v.txt", record.Position.String(), record.Week)
				playerFile := fileMap[fileName]
				if playerFile == nil {
					file, err := os.Create(fmt.Sprintf("%v/%v", testOutDir, fileName))
					fantasyfootball.HandleError(err)
					csvWriter := csv.NewWriter(file)
					defer file.Close()
					defer csvWriter.Flush()
					csvWriter.Comma = '\t'
					headers := make([]string, len(record.Stats))
					csvHeaders := make([]string, len(record.Stats)+2)
					csvHeaders[0] = name
					csvHeaders[1] = team
					idx := 0
					for key := range record.Stats {
						headers[idx] = key
						csvHeaders[idx+2] = key
						idx++
					}
					csvWriter.Write(csvHeaders)
					playerFile = &ParseOutput{
						writer:  csvWriter,
						headers: headers,
					}
					fileMap[fileName] = playerFile
				}
				line := make([]string, len(record.Stats)+2)
				line[0] = record.Name
				line[1] = record.Team
				for i, v := range playerFile.headers {
					line[i+2] = fmt.Sprintf("%v", record.Stats[v])
				}
				playerFile.writer.Write(line)
			}
			return minWeek
		}(playerRecordsChannel)
	//rename current week file
	newMinWeek := fmt.Sprintf("%v", minWeekResult-1)
	for fName := range fileMap {
		newFilename := placeholderRegex.ReplaceAllString(fName, newMinWeek)
		if newFilename != fName {
			fullOldFilename := fmt.Sprintf("%v/%v", testOutDir, fName)
			oldFile, err := os.Open(fullOldFilename)
			fantasyfootball.HandleError(err)
			newFile, err := os.Create(fmt.Sprintf("%v/%v", testOutDir, newFilename))
			defer newFile.Close()
			fantasyfootball.HandleError(err)
			_, err = io.Copy(newFile, oldFile)
			fantasyfootball.HandleError(err)
			err = oldFile.Close()
			fantasyfootball.HandleError(err)
			err = os.Remove(fullOldFilename)
			fantasyfootball.HandleError(err)
		}
	}
}

func (parser *AccuscoreParser) Parse(content <-chan io.ReadCloser) <-chan PlayerRecord {
	records := make(chan PlayerRecord)
	go func(in <-chan io.ReadCloser, out chan<- PlayerRecord) {
		var recordWg sync.WaitGroup
		for cntnt := range in {
			recordWg.Add(1)
			go func(stream io.Reader, output chan<- PlayerRecord, wg *sync.WaitGroup) {
				parser.parseContent(stream, output)
				wg.Done()
			}(cntnt, out, &recordWg)
		}
		recordWg.Wait()
		close(out)
	}(content, records)
	return records
}

func (parser *AccuscoreParser) parseContent(in io.Reader, out chan<- PlayerRecord) {
	content, err := ioutil.ReadAll(in)
	fantasyfootball.HandleError(err)

	//extract position from page
	posMatch := parser.position.FindSubmatch(content)
	if posMatch == nil {
		fantasyfootball.HandleError(errors.New(fmt.Sprintf("unable to parse position from file: %s", in)))
	}
	pos := string(posMatch[1])

	//extract week from page
	weekMatch := parser.week.FindSubmatch(content)
	if weekMatch == nil {
		fantasyfootball.HandleError(errors.New(fmt.Sprintf("unable to parse week from file: %s", in)))
	}
	var weekString string
	if parseWeekMatch := parser.parseWeek.FindSubmatch(weekMatch[1]); parseWeekMatch != nil {
		weekString = fmt.Sprintf("%s", parseWeekMatch[1])
	} else {
		weekString = parser.currWeekString
	}
	weekInt, err := strconv.ParseInt(weekString, 10, 0)
	if err != nil {
		weekInt = CurrWeekPlaceholder
	}

	//extract header
	header := parser.header.Find(content)
	headerColumns := parser.data.FindAllSubmatch(header, -1)
	headerMap := make(map[int]string)
	for i, headMatch := range headerColumns {
		head := headMatch[1]
		headerMap[i] = fmt.Sprintf("%s", head)
	}

	//extract player info from page
	rows := parser.row.FindAll(content, -1)
	for _, row := range rows {
		columns := parser.data.FindAllSubmatch(row, -1)
		dataMap := make(map[string]string)
		for i, column := range columns {
			columnData := fmt.Sprintf("%s", column[1])
			dataHeader := headerMap[i]
			dataMap[dataHeader] = columnData
		}
		footbalPosition, err := fantasyfootball.ParsePosition(pos)
		fantasyfootball.HandleError(err)
		record := PlayerRecord{
			Week:     int(weekInt),
			Position: footbalPosition,
			Name:     dataMap[parser.nameString],
			Team:     dataMap[parser.teamString],
			Stats:    make(map[string]float64),
		}
		delete(dataMap, parser.nameString)
		delete(dataMap, parser.teamString)
		for key, value := range dataMap {
			parsedFloat, err := strconv.ParseFloat(value, 64)
			fantasyfootball.HandleError(err)
			record.Stats[key] = parsedFloat
		}
		out <- record
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
	if week == CurrWeekPlaceholder {
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

	header := headerRegex.Find(content)
	headerColumn := dataRegex.FindAllSubmatch(header, -1)
	csvHeader := make([]string, len(headerColumn))
	for i, head := range headerColumn {
		csvHeader[i] = fmt.Sprintf("%s", head[1])
	}
	csvWriter.Write(csvHeader)
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
		weekInt = CurrWeekPlaceholder
	}
	return int(weekInt), pos
}
