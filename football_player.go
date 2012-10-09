package fantasyfootball

import (
	"errors"
	"fmt"
)

type Position uint8

const (
	DST Position = iota
	K
	QB
	RB
	TE
	WR
)

func (pos Position) ToString() string {
	switch pos {
	case DST:
		return "DEF-ST"
	case K:
		return "K"
	case QB:
		return "QB"
	case RB:
		return "RB"
	case TE:
		return "TE"
	case WR:
		return "WR"
	}
	return ""
}

func ParsePosition(positionString string) (pos Position, err error) {
	switch positionString {
	case "DEF-ST":
		return DST, nil
	case "K":
		return K, nil
	case "QB":
		return QB, nil
	case "RB":
		return RB, nil
	case "TE":
		return TE, nil
	case "WR":
		return WR, nil
	default:
		return DST, errors.New(fmt.Sprintf("Unable to parse as a position: %v", positionString))
	}
	return DST, errors.New(fmt.Sprintf("Unable to parse as a position: %v", positionString))
}

type FootballPlayer struct {
	Name         string
	Team         string
	Position     Position
	Points       [SEASON_LENGTH]int
	totalPoints_ int
}

func (player *FootballPlayer) TotalPoints() int {
	if player.totalPoints_ == 0 {
		sum := 0
		for _, v := range player.Points {
			sum += v
		}
		player.totalPoints_ = sum
	}
	return player.totalPoints_
}

type ByTotalPointsAsc struct {
	players []*FootballPlayer
}

func (s *ByTotalPointsAsc) Len() int { return len(s.players) }

func (s *ByTotalPointsAsc) Less(i, j int) bool {
	return s.players[i].TotalPoints() < s.players[j].TotalPoints()
}

func (s *ByTotalPointsAsc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}

type ByTotalPointsDesc struct {
	players []*FootballPlayer
}

func (s *ByTotalPointsDesc) Len() int { return len(s.players) }

func (s *ByTotalPointsDesc) Less(i, j int) bool {
	return s.players[i].TotalPoints() > s.players[j].TotalPoints()
}

func (s *ByTotalPointsDesc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}

type ByWeekPointsAsc struct {
	players []*FootballPlayer
	week    int
}

func (s *ByWeekPointsAsc) Len() int { return len(s.players) }

func (s *ByWeekPointsAsc) Less(i, j int) bool {
	return s.players[i].Points[s.week-1] < s.players[j].Points[s.week-1]
}

func (s *ByWeekPointsAsc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}

type ByWeekPointsDesc struct {
	players []*FootballPlayer
	week    int
}

func (s *ByWeekPointsDesc) Len() int { return len(s.players) }

func (s *ByWeekPointsDesc) Less(i, j int) bool {
	return s.players[i].Points[s.week-1] > s.players[j].Points[s.week-1]
}

func (s *ByWeekPointsDesc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}
