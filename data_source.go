package fantasyfootball

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

type DataSource interface {
	DefenseSpecialTeams() []*FootballPlayer
	Kickers() []*FootballPlayer
	Quarterbacks() []*FootballPlayer
	RunningBacks() []*FootballPlayer
	TightEnds() []*FootballPlayer
	WideReceivers() []*FootballPlayer
}

type NormalizedDataSource struct {
	defaultDst  *FootballPlayer
	defaultK    *FootballPlayer
	defaultQb   *FootballPlayer
	defaultRb   *FootballPlayer
	defaultTe   *FootballPlayer
	defaultWr   *FootballPlayer
}

func NewNormalizedDataSource(dataSource DataSource) *NormalizedDataSource {
	normalizedDataSource := &NormalizedDataSource{
		defaultDst: &FootballPlayer{position:DST},
		defaultK: &FootballPlayer{position:K},
		defaultQb: &FootballPlayer{position:QB},
		defaultRb: &FootballPlayer{position:RB},
		defaultTe: &FootballPlayer{position:TE},
		defaultWr: &FootballPlayer{position:WR},
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
	return normalizedDataSource
}

func normalizePlayers(offset int, players []*FootballPlayer) *FootballPlayer {
	defaultPlayer := &FootballPlayer{
		name:     "default",
	}
	for week := 1; week <= SEASON_LENGTH; week++ {
		sort.Sort(&ByWeekPointsDesc{players, week})
		defaultPlayer.points[week-1] = players[offset].points[week-1]
	}
	// associate a name with the "default" player, for funsies only since the point values are taken from the weekly nth best player, not the season's nth best player
	sort.Sort(&ByTotalPointsDesc{players})
	defaultPlayer.team = fmt.Sprintf("~%s", players[offset].name)
	// normalize each player to the default player
	for _, p := range players {
		normalizePlayer(defaultPlayer, p)
	}
	return defaultPlayer
}

func normalizePlayer(defaultPlayer *FootballPlayer, player *FootballPlayer) {
	for week := 1; week <= SEASON_LENGTH; week++ {
		player.points[week-1] -= defaultPlayer.points[week-1];
		player.totalPoints_ = 0
	}
}

type FileDataSource struct {
	dir       string
	startWeek int
	endWeek   int

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
	}
	fileDataSource.loadFiles()
	return fileDataSource
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
	var waitGroup sync.WaitGroup
	waitGroup.Add(6)
	go func() { fds.loadDsts(); waitGroup.Done() }()
	go func() { fds.loadKs(); waitGroup.Done() }()
	go func() { fds.loadQbs(); waitGroup.Done() }()
	go func() { fds.loadRbs(); waitGroup.Done() }()
	go func() { fds.loadTes(); waitGroup.Done() }()
	go func() { fds.loadWrs(); waitGroup.Done() }()
	waitGroup.Wait()
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
	var offset int
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
