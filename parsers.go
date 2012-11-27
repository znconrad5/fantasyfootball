package fantasyfootball

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func newDstParser() *Parser {
	return &Parser{
		numHeaders:  12,
		players:     make(map[string]*FootballPlayer),
		pointParser: new(DstPointParser),
		position:    DST,
	}
}

func newKParser() *Parser {
	return &Parser{
		numHeaders:  10,
		players:     make(map[string]*FootballPlayer),
		pointParser: new(KPointParser),
		position:    K,
	}
}

func newQbParser() *Parser {
	return &Parser{
		numHeaders:  12,
		players:     make(map[string]*FootballPlayer),
		pointParser: new(QbPointParser),
		position:    QB,
	}
}

func newRbParser() *Parser {
	return &Parser{
		numHeaders:  11,
		players:     make(map[string]*FootballPlayer),
		pointParser: new(RbPointParser),
		position:    RB,
	}
}

func newTeParser() *Parser {
	return &Parser{
		numHeaders:  8,
		players:     make(map[string]*FootballPlayer),
		pointParser: new(RecPointParser),
		position:    TE,
	}
}

func newWrParser() *Parser {
	return &Parser{
		numHeaders:  8,
		players:     make(map[string]*FootballPlayer),
		pointParser: new(RecPointParser),
		position:    WR,
	}
}

type PointParser interface {
	parsePoints(header map[string]int, statsLine []string) int
}

type Parser struct {
	numHeaders  int
	players     map[string]*FootballPlayer
	pointParser PointParser
	position    Position
}

func (p *Parser) parseFile(fileName string, week int) {
	// open file
	file, err := os.Open(fileName)
	HandleError(err)
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	// read in header lines
	header := make(map[string]int)
	line, err := reader.Read()
	HandleError(err)
	for i, h := range line {
		header[h] = i
	}
	// read in player/stats
	for statsLine, err := reader.Read(); err != io.EOF; statsLine, err = reader.Read() {
		HandleError(err)
		playerName := statsLine[header["PLAYER"]]
		team := statsLine[header["TEAM"]]
		playerKey := fmt.Sprintf("%s (%s)", playerName, team)
		player, ok := p.players[playerKey]
		if !ok {
			player = &FootballPlayer{
				Name:     playerName,
				Team:     team,
				Position: p.position,
			}
			p.players[playerKey] = player
		}
		points := p.pointParser.parsePoints(header, statsLine)
		player.Points[week-1] = points
	}
}

type DstPointParser struct{}

func (p *DstPointParser) parsePoints(header map[string]int, statsLine []string) int {
	// Each sack:	1 point
	points := parsePointValue(1, statsLine[header["SACK"]])
	// Each interception:	2 points
	// Each fumble recovery:	2 points
	points += parsePointValue(2, statsLine[header["TO"]])
	// Each TD:	6 points
	points += parsePointValue(6, statsLine[header["INTTD"]])
	// Kickoff and Punt Return Touchdowns:	6 points
	points += parsePointValue(6, statsLine[header["KRETTD"]])
	points += parsePointValue(6, statsLine[header["PRETTD"]])
	// Each safety:	2 points
	// Each blocked kick:	2 points
	pointsAgainst, err := strconv.ParseFloat(statsLine[header["PA"]], 32)
	HandleError(err)
	// Shutout:	10 points
	// 1-6 points allowed:	8 points
	// 7-13 points allowed:	6 points
	// 14-20 points allowed:	3 point
	// 21-27 points allowed:	1 points
	// 28-34 points allowed:	-1 points
	// 35+ points allowed:	-3 points
	if pointsAgainst == 0 {
		// if points is currently 0 this is a bye week
		if points != 0 {
			points += 10 * 100
		}
	} else if pointsAgainst < 1 {
		points += int((10 - 2*pointsAgainst) * 100)
	} else if pointsAgainst < 7 {
		points += int((8 - 2*(pointsAgainst-1)/6) * 100)
	} else if pointsAgainst < 14 {
		points += int((6 - 3*(pointsAgainst-7)/7) * 100)
	} else if pointsAgainst < 21 {
		points += int((3 - 2*(pointsAgainst-14)/7) * 100)
	} else if pointsAgainst < 28 {
		points += int((1 - 2 * (pointsAgainst-21)/7) * 100)
	} else if pointsAgainst < 35 {
		points += int((-1 - 2 * (pointsAgainst-28)/7) * 100)
	} else {
		points += -3 * 100
	}
	return points
}

type KPointParser struct{}

func (p *KPointParser) parsePoints(header map[string]int, statsLine []string) int {
	// Field goal 0-19 yards:	3 points
	// Field goal 20-29 yards:	3 points
	// Field goal 30-39 yards:	3 points
	points := parsePointValue(3, statsLine[header["FG0-39"]])
	// Field goal 40-49 yards:	4 points
	points += parsePointValue(4, statsLine[header["FG40-49"]])
	// Field goal 50+ yards:	5 points
	points += parsePointValue(5, statsLine[header["FG50"]])
	// Each extra point:	1 point
	points += parsePointValue(1, statsLine[header["XPM"]])
	return points
}

type QbPointParser struct{}

func (p *QbPointParser) parsePoints(header map[string]int, statsLine []string) int {
	// Every 25 passing yards:	1 point
	points := parsePointValue(1, statsLine[header["PASSYD"]]) / 25
	// Each passing TD:	6 points
	points += parsePointValue(6, statsLine[header["PASSTD"]])
	// Each interception thrown:	-1 points
	points += parsePointValue(-1, statsLine[header["INT"]])
	// Every 10 rushing yards:	1 point
	points += parsePointValue(1, statsLine[header["RUSHYD"]]) / 10
	// Each rushing TD:	6 points
	points += parsePointValue(6, statsLine[header["RUSHTD"]])
	// Each fumble lost:	-1 points
	points += parsePointValue(-2, statsLine[header["FUML"]])
	return points
}

type RbPointParser struct{}

func (p *RbPointParser) parsePoints(header map[string]int, statsLine []string) int {
	// Every 10 rushing yards:	1 point
	points := parsePointValue(1, statsLine[header["RUSHYD"]]) / 10
	// Each rushing TD:	6 points
	points += parsePointValue(6, statsLine[header["RUSHTD"]])
	// Every 10 receiving yards:	1 point
	points += parsePointValue(1, statsLine[header["RECYD"]]) / 10
	// Each receiving TD:	6 points
	points += parsePointValue(6, statsLine[header["RECTD"]])
	// Each fumble lost:	-1 points
	points += parsePointValue(-1, statsLine[header["FUML"]])
	return points
}

type RecPointParser struct{}

func (p *RecPointParser) parsePoints(header map[string]int, statsLine []string) int {
	// Every 10 receiving yards:	1 point
	points := parsePointValue(1, statsLine[header["RECYD"]]) / 10
	// Each receiving TD:	6 points
	points += parsePointValue(6, statsLine[header["RECTD"]])
	// Each fumble lost:	-1 points
	points += parsePointValue(-1, statsLine[header["FUML"]])
	return points
}

func validateSingleValue(line []string) {
	if len(line) != 1 {
		log.Fatalf("Expected a single value, got: %s", line)
	}
}

func parsePointValue(weight int, s string) int {
	f, err := strconv.ParseFloat(s, 32)
	HandleError(err)
	return weight * int(f*100)
}
