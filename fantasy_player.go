package fantasyfootball

import (
	"fmt"
	"sort"
)

type FantasyPlayer struct {
	name string
	// team
	dsts []*FootballPlayer
	ks   []*FootballPlayer
	qbs  []*FootballPlayer
	rbs  []*FootballPlayer
	tes  []*FootballPlayer
	wrs  []*FootballPlayer
	// needed for calculations
	defaultRb   *FootballPlayer
	defaultWr   *FootballPlayer
	defaultFlex *FootballPlayer
}

func newFantasyPlayer(name string, data DataSource) *FantasyPlayer {
	return &FantasyPlayer{
		name:        name,
		defaultRb:   data.defaultRb,
		defaultWr:   data.defaultWr,
		defaultFlex: data.defaultFlex,
	}
}

func (fp *FantasyPlayer) draft(player *FootballPlayer) {
	switch player.position {
	case DST:
		fp.dsts = append(fp.dsts, player)
	case K:
		fp.ks = append(fp.ks, player)
	case QB:
		fp.qbs = append(fp.qbs, player)
	case RB:
		fp.rbs = append(fp.rbs, player)
	case TE:
		fp.tes = append(fp.tes, player)
	case WR:
		fp.wrs = append(fp.wrs, player)
	}
}

func remove(slice []*FootballPlayer, x *FootballPlayer) []*FootballPlayer {
	var i int
	lastIndex := len(slice) - 1
	for i = lastIndex; slice[i] != x; i-- {
	}
	if i != lastIndex {
		slice[i] = slice[lastIndex]
	}
	return slice[:lastIndex]
}

func (fp *FantasyPlayer) undraft(player *FootballPlayer) {
	switch player.position {
	case DST:
		fp.dsts = remove(fp.dsts, player)
	case K:
		fp.ks = remove(fp.ks, player)
	case QB:
		fp.qbs = remove(fp.qbs, player)
	case RB:
		fp.rbs = remove(fp.rbs, player)
	case TE:
		fp.tes = remove(fp.tes, player)
	case WR:
		fp.wrs = remove(fp.wrs, player)
	}
}

func (fp *FantasyPlayer) points(week int) int {
	points := 0
	for _, players := range [...][]*FootballPlayer{fp.dsts, fp.ks, fp.qbs, fp.tes} {
		positionMax := 0
		for _, player := range players {
			positionMax = max(positionMax, player.points[week-1])
		}
		points += positionMax
	}
	for _, players := range [...][]*FootballPlayer{fp.rbs, fp.wrs} {
		switch len(players) {
		case 2:
			points += max(0, players[1].points[week-1])
			fallthrough
		case 1:
			points += max(0, players[0].points[week-1])
		case 0:
			continue
		default:
			byWeekPoints := &ByWeekPointsDesc{players, week}
			if !sort.IsSorted(byWeekPoints) {
				sort.Sort(byWeekPoints)
			}
			points += max(0, players[0].points[week-1]) + max(0, players[1].points[week-1])
		}
	}
	flexRb := 0
	if len(fp.rbs) >= 3 {
		flexRb = max(0, fp.rbs[2].points[week-1]+fp.defaultRb.points[week-1]-fp.defaultFlex.points[week-1])
	}
	flexWr := 0
	if len(fp.wrs) >= 3 {
		flexWr = max(0, fp.wrs[2].points[week-1]+fp.defaultWr.points[week-1]-fp.defaultFlex.points[week-1])
	}
	points += max(flexRb, flexWr)
	return points
}

/*
	A slow estimate of the season point total, based on the assumption lineup changes can be made each week.
*/
func (fp *FantasyPlayer) totalPoints() int {
	totalPoints := 0
	for week := 1; week <= SEASON_LENGTH; week++ {
		totalPoints += fp.points(week)
	}
	return totalPoints
}

/*
	A quick estimate of the season point total, based on the assumption no lineup changes are made each week.
*/
func (fp *FantasyPlayer) estimateTotalPoints() int {
	totalPoints := 0
	for _, players := range [...][]*FootballPlayer{fp.dsts, fp.ks, fp.qbs, fp.tes} {
		positionMax := 0
		for _, player := range players {
			positionMax = max(positionMax, player.totalPoints())
		}
		totalPoints += positionMax
	}
	for _, players := range [...][]*FootballPlayer{fp.rbs, fp.wrs} {
		switch len(players) {
		case 2:
			totalPoints += max(0, players[1].totalPoints())
			fallthrough
		case 1:
			totalPoints += max(0, players[0].totalPoints())
		case 0:
			continue
		default:
			byTotalPoints := &ByTotalPointsDesc{players}
			if !sort.IsSorted(byTotalPoints) {
				sort.Sort(byTotalPoints)
			}
			totalPoints += max(0, players[0].totalPoints()) + max(0, players[1].totalPoints())
		}
	}
	flexRb := 0
	if len(fp.rbs) >= 3 {
		flexRb = max(0, fp.rbs[2].totalPoints()+fp.defaultRb.totalPoints()-fp.defaultFlex.totalPoints())
	}
	flexWr := 0
	if len(fp.wrs) >= 3 {
		flexWr = max(0, fp.wrs[2].totalPoints()+fp.defaultWr.totalPoints()-fp.defaultFlex.totalPoints())
	}
	totalPoints += max(flexRb, flexWr)
	return totalPoints
}

func (fp *FantasyPlayer) String() string {
	str := fmt.Sprintf("%v\nQBS:\n", fp.name)
	for _, v := range fp.qbs {
		str += fmt.Sprintf("\t%v\n", v)
	}
	str += "RBS:\n"
	for _, v := range fp.rbs {
		str += fmt.Sprintf("\t%v\n", v)
	}
	str += "WRS:\n"
	for _, v := range fp.wrs {
		str += fmt.Sprintf("\t%v\n", v)
	}
	str += "TES:\n"
	for _, v := range fp.tes {
		str += fmt.Sprintf("\t%v\n", v)
	}
	str += "KS:\n"
	for _, v := range fp.ks {
		str += fmt.Sprintf("\t%v\n", v)
	}
	str += "DSTS:\n"
	for _, v := range fp.dsts {
		str += fmt.Sprintf("\t%v\n", v)
	}
	return str
}
