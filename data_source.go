package fantasyfootball

import (
	"fmt"
	"math"
	"sort"
)

type DataSource struct {
	dir       string
	startWeek int
	endWeek   int

	dsts       []*FootballPlayer // defenses/special teams
	defaultDst *FootballPlayer
	ks         []*FootballPlayer // kickers
	defaultK   *FootballPlayer
	qbs        []*FootballPlayer // quarterbacks
	defaultQb  *FootballPlayer
	rbs        []*FootballPlayer // running backs
	defaultRb  *FootballPlayer
	tes        []*FootballPlayer // tight ends
	defaultTe  *FootballPlayer
	wrs        []*FootballPlayer // wide receivers
	defaultWr  *FootballPlayer
}

func NewDataSource(dir string, startWeek, endWeek int) *DataSource {
	return &DataSource{
		dir:       dir,
		startWeek: startWeek,
		endWeek:   endWeek,
	}
}

func (loader *DataSource) LoadAll() {
	c := make(chan []*FootballPlayer, 6)
	go func() { c <- depair(loader.loadDsts()) }()
	go func() { c <- depair(loader.loadKs()) }()
	go func() { c <- depair(loader.loadQbs()) }()
	go func() { c <- depair(loader.loadRbs()) }()
	go func() { c <- depair(loader.loadTes()) }()
	go func() { c <- depair(loader.loadWrs()) }()
	for i := 0; i < 6; i++ {
		<-c
	}
}

func (loader *DataSource) loadDsts() ([]*FootballPlayer, *FootballPlayer) {
	parser := newDstParser()
	loader.dsts, loader.defaultDst = loader.load(parser, DST)
	return loader.dsts, loader.defaultDst
}

func (loader *DataSource) loadKs() ([]*FootballPlayer, *FootballPlayer) {
	parser := newKParser()
	loader.ks, loader.defaultK = loader.load(parser, K)
	return loader.ks, loader.defaultK
}

func (loader *DataSource) loadQbs() ([]*FootballPlayer, *FootballPlayer) {
	parser := newQbParser()
	loader.qbs, loader.defaultQb = loader.load(parser, QB)
	return loader.qbs, loader.defaultQb
}

func (loader *DataSource) loadRbs() ([]*FootballPlayer, *FootballPlayer) {
	parser := newRbParser()
	loader.rbs, loader.defaultRb = loader.load(parser, RB)
	return loader.rbs, loader.defaultRb
}

func (loader *DataSource) loadTes() ([]*FootballPlayer, *FootballPlayer) {
	parser := newTeParser()
	loader.tes, loader.defaultTe = loader.load(parser, TE)
	return loader.tes, loader.defaultTe
}

func (loader *DataSource) loadWrs() ([]*FootballPlayer, *FootballPlayer) {
	parser := newWrParser()
	loader.wrs, loader.defaultWr = loader.load(parser, WR)
	return loader.wrs, loader.defaultWr
}

func (loader *DataSource) load(parser *Parser, position Position) ([]*FootballPlayer, *FootballPlayer) {
	var fileName string
	var offset int
	switch position {
	case DST:
		fileName = "dst"
		offset = 2 * 9 // assume each player drafts 2
	case K:
		fileName = "k"
		offset = 2 * 9 // assume each player drafts 2
	case QB:
		fileName = "qb"
		offset = 2 * 9 // assume each player drafts 2
	case RB:
		fileName = "rb"
		offset = int(math.Ceil(3.5 * 9)) // assume each player drafts 3.5
	case TE:
		fileName = "te"
		offset = 2 * 9 // assume each player drafts 2
	case WR:
		fileName = "wr"
		offset = int(math.Ceil(3.5 * 9)) // assume each player drafts 3.5
	}
	for week := loader.startWeek; week <= loader.endWeek; week++ {
		parser.parseFile(fmt.Sprintf("%s/%s_%d.txt", loader.dir, fileName, week), week)
	}
	players := make([]*FootballPlayer, len(parser.players))
	i := 0
	for _, v := range parser.players {
		players[i] = v
		i++
	}
	defaultPlayer := &FootballPlayer{
		name:     "default",
		position: position,
	}
	// the "default" player is a guess of the best undrafted player for a position each week
	for week := loader.startWeek; week <= loader.endWeek; week++ {
		sort.Sort(&ByWeekPointsDesc{players, week})
		defaultPlayer.points[week-1] = players[offset].points[week-1]
	}
	// associate a name with the "default" player, for funsies only since the point values are taken from the weekly nth best player, not the season's nth best player
	sort.Sort(&ByTotalPointsDesc{players})
	defaultPlayer.team = fmt.Sprintf("~%s", players[offset].name)
	return players, defaultPlayer
}

func depair(players []*FootballPlayer, defaultPlayer *FootballPlayer) []*FootballPlayer {
	return players
}
