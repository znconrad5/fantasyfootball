package fantasyfootball

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

type DataSource interface {
	AllPlayers() map[string]*FootballPlayer
	DefenseSpecialTeams() []*FootballPlayer
	Kickers() []*FootballPlayer
	Quarterbacks() []*FootballPlayer
	RunningBacks() []*FootballPlayer
	TightEnds() []*FootballPlayer
	WideReceivers() []*FootballPlayer
}

type NormalizedDataSource struct {
	DataSource
	defaultDst  *FootballPlayer
	defaultK    *FootballPlayer
	defaultQb   *FootballPlayer
	defaultRb   *FootballPlayer
	defaultTe   *FootballPlayer
	defaultWr   *FootballPlayer
	defaultFlex   *FootballPlayer
}

func NewNormalizedDataSource(dataSource DataSource) *NormalizedDataSource {
	normalizedDataSource := &NormalizedDataSource{
		dataSource,
		&FootballPlayer{Position:DST},
		&FootballPlayer{Position:K},
		&FootballPlayer{Position:QB},
		&FootballPlayer{Position:RB},
		&FootballPlayer{Position:TE},
	    &FootballPlayer{Position:WR},
		new(FootballPlayer),
	}
	
	var waitGroup sync.WaitGroup
	waitGroup.Add(6)
	go func() {
		normalizedDataSource.defaultDst = normalizePlayers(9*4/3, dataSource.DefenseSpecialTeams())
		waitGroup.Done()
	}()
	go func() {
		normalizedDataSource.defaultK = normalizePlayers(9*4/3, dataSource.Kickers())
		waitGroup.Done()
	}()
	go func() {
		normalizedDataSource.defaultQb = normalizePlayers(9*2, dataSource.Quarterbacks())
		waitGroup.Done()
	}()
	go func() {
		normalizedDataSource.defaultRb = normalizePlayers(int(math.Ceil(9*4.5)), dataSource.RunningBacks())
		waitGroup.Done()
	}()
	go func() {
		normalizedDataSource.defaultTe = normalizePlayers(9*4/3, dataSource.TightEnds())
		waitGroup.Done()
	}()
	go func() {
		normalizedDataSource.defaultWr = normalizePlayers(int(math.Ceil(9*4.5)), dataSource.WideReceivers())
		waitGroup.Done()
	}()
	waitGroup.Wait()
	if normalizedDataSource.defaultRb.TotalPoints() > normalizedDataSource.defaultWr.TotalPoints() {
		normalizedDataSource.defaultFlex = normalizedDataSource.defaultRb
	} else {
		normalizedDataSource.defaultFlex = normalizedDataSource.defaultWr
	}
	return normalizedDataSource
}

func normalizePlayers(offset int, players []*FootballPlayer) *FootballPlayer {
	defaultPlayer := &FootballPlayer{
		Name:     "default",
	}
	for week := 1; week <= SEASON_LENGTH; week++ {
		sort.Sort(&ByWeekPointsDesc{players, week})
		defaultPlayer.Points[week-1] = players[offset].Points[week-1]
	}
	// associate a name with the "default" player, for funsies only since the point values are taken from the weekly nth best player, not the season's nth best player
	sort.Sort(&ByTotalPointsDesc{players})
	defaultPlayer.Team = fmt.Sprintf("~%s", players[offset].Name)
	// normalize each player to the default player
	for _, p := range players {
		normalizePlayer(defaultPlayer, p)
	}
	return defaultPlayer
}

func normalizePlayer(defaultPlayer *FootballPlayer, player *FootballPlayer) {
	for week := 1; week <= SEASON_LENGTH; week++ {
		player.Points[week-1] -= defaultPlayer.Points[week-1];
		player.totalPoints_ = 0
	}
}

type FileDataSource struct {
	dir       string
	startWeek int
	endWeek   int

	allPlayers  map[string]*FootballPlayer
	dsts        []*FootballPlayer // defenses/special teams
	ks          []*FootballPlayer // kickers
	qbs         []*FootballPlayer // quarterbacks
	rbs         []*FootballPlayer // running backs
	tes         []*FootballPlayer // tight ends
	wrs         []*FootballPlayer // wide receivers
}

func NewFileDataSource(dir string, startWeek, endWeek int) *FileDataSource {
	fileDataSource := &FileDataSource{
		dir:        dir,
		startWeek:  startWeek,
		endWeek:    endWeek,
		allPlayers: make(map[string]*FootballPlayer),
	}
	fileDataSource.loadFiles()
	return fileDataSource
}

func (fds *FileDataSource) AllPlayers() map[string]*FootballPlayer {
	return fds.allPlayers
}

func (fds *FileDataSource) DefenseSpecialTeams() []*FootballPlayer {
	return fds.dsts
}

func (fds *FileDataSource) Kickers() []*FootballPlayer {
	return fds.ks
}

func (fds *FileDataSource) Quarterbacks() []*FootballPlayer {
	return fds.qbs
}

func (fds *FileDataSource) RunningBacks() []*FootballPlayer {
	return fds.rbs
}

func (fds *FileDataSource) TightEnds() []*FootballPlayer {
	return fds.tes
}

func (fds *FileDataSource) WideReceivers() []*FootballPlayer {
	return fds.wrs
}

func (fds *FileDataSource) loadFiles() {
	c := make(chan []*FootballPlayer, 6)
	go func() { c <- fds.loadDsts() }()
	go func() { c <- fds.loadKs() }()
	go func() { c <- fds.loadQbs() }()
	go func() { c <- fds.loadRbs() }()
	go func() { c <- fds.loadTes() }()
	go func() { c <- fds.loadWrs() }()
	for i := 0; i < 6; i++ {
		for _, p := range <-c {
			fds.allPlayers[fmt.Sprintf("%s (%s)", p.Name, p.Team)] = p
		}
	}
}

func (fds *FileDataSource) loadDsts() []*FootballPlayer {
	parser := newDstParser()
	fds.dsts = fds.load(parser, DST)
	return fds.dsts
}

func (fds *FileDataSource) loadKs() []*FootballPlayer {
	parser := newKParser()
	fds.ks = fds.load(parser, K)
	return fds.ks
}

func (fds *FileDataSource) loadQbs() []*FootballPlayer {
	parser := newQbParser()
	fds.qbs = fds.load(parser, QB)
	return fds.qbs
}

func (fds *FileDataSource) loadRbs() []*FootballPlayer {
	parser := newRbParser()
	fds.rbs = fds.load(parser, RB)
	return fds.rbs
}

func (fds *FileDataSource) loadTes() []*FootballPlayer {
	parser := newTeParser()
	fds.tes = fds.load(parser, TE)
	return fds.tes
}

func (fds *FileDataSource) loadWrs() []*FootballPlayer {
	parser := newWrParser()
	fds.wrs = fds.load(parser, WR)
	return fds.wrs
}

func (fds *FileDataSource) load(parser *Parser, position Position) []*FootballPlayer {
	var fileName string
	switch position {
	case DST:
		fileName = "def-st"
	case K:
		fileName = "k"
	case QB:
		fileName = "qb"
	case RB:
		fileName = "rb"
	case TE:
		fileName = "te"
	case WR:
		fileName = "wr"
	}
	for week := fds.startWeek; week <= fds.endWeek; week++ {
		parser.parseFile(fmt.Sprintf("%s/%s_%d.txt", fds.dir, fileName, week), week)
	}
	players := make([]*FootballPlayer, len(parser.players))
	i := 0
	for _, v := range parser.players {
		players[i] = v
		i++
	}
	return players
}
