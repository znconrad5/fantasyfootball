package fantasyfootball

import (
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
}

func newFantasyPlayer(name string, data *DataSource) *FantasyPlayer {
	return &FantasyPlayer{
		name: name,
		dsts: []*FootballPlayer{data.defaultDst},
		ks: []*FootballPlayer{data.defaultK},
		qbs: []*FootballPlayer{data.defaultQb},
		rbs: []*FootballPlayer{data.defaultRb, data.defaultRb, data.defaultRb},
		tes: []*FootballPlayer{data.defaultTe},
		wrs: []*FootballPlayer{data.defaultWr, data.defaultWr, data.defaultWr},
	}
}

func (fp *FantasyPlayer) draft(player *FootballPlayer) {
	switch(player.position) {
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
	lastIndex := len(slice)-1
	for i=lastIndex; slice[i] != x; i-- {}
	if i!=lastIndex {
		slice[i] = slice[lastIndex]
	}
	return slice[:lastIndex]
}

func (fp *FantasyPlayer) undraft(player *FootballPlayer) {
	switch(player.position) {
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
	for _, v := range [...][]*FootballPlayer{fp.dsts, fp.ks, fp.qbs, fp.tes} {
		positionMax := v[0].points[week-1]
		for i:=1; i<len(v); i++ {
			positionMax = max(positionMax, v[i].points[week-1])
		}
		points += positionMax
	}
	for _, v := range [...][]*FootballPlayer{fp.rbs, fp.wrs} {
		byWeekPoints := &ByWeekPointsAsc{v, week}
		if !sort.IsSorted(byWeekPoints) {
			sort.Sort(byWeekPoints)
		}
		points += v[len(v)-1].points[week-1]
	}
	points += fp.rbs[len(fp.rbs)-2].points[week-1]
	points += fp.wrs[len(fp.wrs)-2].points[week-1]
	points += max(fp.rbs[len(fp.rbs)-3].points[week-1], fp.wrs[len(fp.wrs)-3].points[week-1])
	return points
}

func (fp *FantasyPlayer) totalPoints() int {
	totalPoints := 0
	for week := 1; week <= SEASON_LENGTH; week++ {
		totalPoints += fp.points(week)
	}
	return totalPoints
}

func (fp *FantasyPlayer) estimateTotalPoints() int {
	totalPoints := 0
	for _, v := range [...][]*FootballPlayer{fp.dsts, fp.ks, fp.qbs, fp.tes} {
		positionMax := v[0].totalPoints()
		for i:=1; i<len(v); i++ {
			positionMax = max(positionMax, v[i].totalPoints())
		}
		totalPoints += positionMax
	}
	for _, v := range [...][]*FootballPlayer{fp.rbs, fp.wrs} {
		byTotalPoints := &ByTotalPointsAsc{v}
		if !sort.IsSorted(byTotalPoints) {
			sort.Sort(byTotalPoints)
		}
		totalPoints += v[len(v)-1].totalPoints()
	}
	totalPoints += fp.rbs[len(fp.rbs)-2].totalPoints()
	totalPoints += fp.wrs[len(fp.wrs)-2].totalPoints()
	totalPoints += max(fp.rbs[len(fp.rbs)-3].totalPoints(), fp.wrs[len(fp.wrs)-2].totalPoints())
	return totalPoints
}
