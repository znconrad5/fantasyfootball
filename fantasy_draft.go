package fantasyfootball

import (
	"math"
	"runtime"
	"sort"
	"strings"
)

const (
	START_DEPTH = 4
)

type FantasyDraft struct {
	players        []*FantasyPlayer
	maxPlayer      *FantasyPlayer
	playersDrafted int
	dsts           *Stack
	ks             *Stack
	qbs            *Stack
	rbs            *Stack
	tes            *Stack
	wrs            *Stack
}

type Move struct {
	Player     *FootballPlayer
	Evaluation int
}

func NewFantasyDraft(names []string, maxPlayerIdx int, data *NormalizedDataSource) *FantasyDraft {
	fd := &FantasyDraft{playersDrafted: 0}
	for i, name := range names {
		fd.players = append(fd.players, newFantasyPlayer(name, data))
		if i == maxPlayerIdx {
			fd.maxPlayer = fd.players[i]
		}
	}

	// the resulting stacks will will be in ascending order, so Pop() will return the best remaining player for that position
	fd.dsts = NewStack()
	sort.Sort(&ByTotalPointsAsc{data.DefenseSpecialTeams()})
	for _, v := range data.DefenseSpecialTeams() {
		fd.dsts.Push(v)
	}
	fd.ks = NewStack()
	sort.Sort(&ByTotalPointsAsc{data.Kickers()})
	for _, v := range data.Kickers() {
		fd.ks.Push(v)
	}
	fd.qbs = NewStack()
	sort.Sort(&ByTotalPointsAsc{data.Quarterbacks()})
	for _, v := range data.Quarterbacks() {
		fd.qbs.Push(v)
	}
	fd.rbs = NewStack()
	sort.Sort(&ByTotalPointsAsc{data.RunningBacks()})
	for _, v := range data.RunningBacks() {
		fd.rbs.Push(v)
	}
	fd.tes = NewStack()
	sort.Sort(&ByTotalPointsAsc{data.TightEnds()})
	for _, v := range data.TightEnds() {
		fd.tes.Push(v)
	}
	fd.wrs = NewStack()
	sort.Sort(&ByTotalPointsAsc{data.WideReceivers()})
	for _, v := range data.WideReceivers() {
		fd.wrs.Push(v)
	}
	return fd
}

func (fd *FantasyDraft) CurrentPlayer() *FantasyPlayer {
	// draft goes 1, 2, 3, ..., 10, 10, 9, 8, ..., 1
	round := fd.playersDrafted / len(fd.players)
	offset := fd.playersDrafted % len(fd.players)
	if round%2 == 0 {
		return fd.players[offset]
	}
	return fd.players[len(fd.players)-1-offset]
}

func (fd *FantasyDraft) length() int {
	return len(fd.players) * 15
}

func (fd *FantasyDraft) Draft(draftee *FootballPlayer) {
	fd.removeFootballPlayer(draftee)
	fd.CurrentPlayer().draft(draftee)
	fd.playersDrafted++
}

func (fd *FantasyDraft) removeFootballPlayer(player *FootballPlayer) {
	var pool *Stack
	switch player.Position {
	case DST:
		pool = fd.dsts
	case K:
		pool = fd.ks
	case QB:
		pool = fd.qbs
	case RB:
		pool = fd.rbs
	case TE:
		pool = fd.tes
	case WR:
		pool = fd.wrs
	}
	pool.Remove(player)
}

func (fd *FantasyDraft) IterativeAlphabeta(stop <-chan bool) <-chan Move {
	moves := make(chan Move)
	go func() {
		defer close(moves)
		remaining := fd.length() - fd.playersDrafted
		cache := make([]*FootballPlayer, fd.length())
		for depth := min(START_DEPTH, remaining); depth <= remaining; depth++ {
			// for depth:=min(START_DEPTH, remaining); depth<=10; depth++ {
			move, val, ok := fd.Alphabeta(depth, math.MinInt32, math.MaxInt32, cache, stop)
			if ok {
				moves <- Move{move, val}
			} else {
				return
			}
		}
	}()
	return moves
}

func (fd *FantasyDraft) Alphabeta(depth, alpha, beta int, cache []*FootballPlayer, stop <-chan bool) (*FootballPlayer, int, bool) {
	runtime.Gosched()
	select {
	case <-stop:
		return nil, 0, false
	default:
		var move *FootballPlayer
		if depth == 0 || fd.IsOver() {
			var value int
			if fd.playersDrafted > 7*len(fd.players) {
				value = fd.evaluate()
			} else {
				value = fd.estimate()
			}
			return move, value, true
		}
		currentPlayer := fd.CurrentPlayer()
		if fd.maxPlayer == currentPlayer {
			for _, v := range fd.generateMoves(cache) {
				draftee := v.Pop()
				currentPlayer.draft(draftee)
				fd.playersDrafted++
				_, moveValue, ok := fd.Alphabeta(depth-1, alpha, beta, cache, stop)
				fd.playersDrafted--
				currentPlayer.undraft(draftee)
				v.Push(draftee)
				if !ok {
					return nil, 0, false
				}
				if moveValue > alpha {
					alpha = moveValue
					move = draftee
					cache[fd.playersDrafted] = draftee
				}
				if beta <= alpha {
					break
				}
			}
			return move, alpha, true
		} else {
			for _, v := range fd.generateMoves(cache) {
				draftee := v.Pop()
				currentPlayer.draft(draftee)
				fd.playersDrafted++
				_, moveValue, ok := fd.Alphabeta(depth-1, alpha, beta, cache, stop)
				fd.playersDrafted--
				currentPlayer.undraft(draftee)
				v.Push(draftee)
				if !ok {
					return nil, 0, false
				}
				if moveValue < beta {
					beta = moveValue
					move = draftee
					cache[fd.playersDrafted] = draftee
				}
				if beta <= alpha {
					break
				}
			}
			return move, beta, true
		}
	}
	panic("alphabeta() unexpectedly returned")
}

func (fd *FantasyDraft) generateMoves(cache []*FootballPlayer) []*Stack {
	if cache[fd.playersDrafted] != nil {
		switch cache[fd.playersDrafted].Position {
		case DST:
			return []*Stack{fd.dsts, fd.rbs, fd.wrs, fd.qbs, fd.tes, fd.ks}
		case K:
			return []*Stack{fd.ks, fd.rbs, fd.wrs, fd.qbs, fd.tes, fd.dsts}
		case QB:
			return []*Stack{fd.qbs, fd.rbs, fd.wrs, fd.tes, fd.ks, fd.dsts}
		case RB:
			return []*Stack{fd.rbs, fd.wrs, fd.qbs, fd.tes, fd.ks, fd.dsts}
		case TE:
			return []*Stack{fd.tes, fd.rbs, fd.wrs, fd.qbs, fd.ks, fd.dsts}
		case WR:
			return []*Stack{fd.wrs, fd.rbs, fd.qbs, fd.tes, fd.ks, fd.dsts}
		}
	}
	byBestMove := &ByBestLikelyMove{fd.CurrentPlayer(), []*Stack{fd.rbs, fd.wrs, fd.qbs, fd.tes, fd.ks, fd.dsts}}
	sort.Sort(byBestMove)
	return byBestMove.stacks
}

type ByBestLikelyMove struct {
	currentPlayer *FantasyPlayer
	stacks        []*Stack
}

func (s *ByBestLikelyMove) Len() int { return len(s.stacks) }

func (s *ByBestLikelyMove) Less(i, j int) bool {
	iDraftee := s.stacks[i].Peek()
	s.currentPlayer.draft(iDraftee)
	iTotal := s.currentPlayer.estimateTotalPoints()
	s.currentPlayer.undraft(iDraftee)

	jDraftee := s.stacks[j].Peek()
	s.currentPlayer.draft(jDraftee)
	jTotal := s.currentPlayer.estimateTotalPoints()
	s.currentPlayer.undraft(jDraftee)

	return iTotal > jTotal
}

func (s *ByBestLikelyMove) Swap(i, j int) {
	s.stacks[i], s.stacks[j] = s.stacks[j], s.stacks[i]
}

func (fd *FantasyDraft) evaluate() int {
	value := 0
	for _, player := range fd.players {
		if player == fd.maxPlayer {
			value += (len(fd.players) - 1) * player.totalPoints()
		} else {
			value -= player.totalPoints()
		}
	}
	return value / (len(fd.players) - 1)
}

func (fd *FantasyDraft) estimate() int {
	value := 0
	for _, player := range fd.players {
		if player == fd.maxPlayer {
			value += (len(fd.players) - 1) * player.estimateTotalPoints()
		} else {
			value -= player.estimateTotalPoints()
		}
	}
	return value / (len(fd.players) - 1)
}

func (fd *FantasyDraft) IsOver() bool {
	return fd.playersDrafted == len(fd.players)*15
}

func (fd *FantasyDraft) String() string {
	return strings.Join([]string{
		fd.players[0].String(),
		fd.players[1].String(),
		fd.players[2].String(),
		fd.players[3].String(),
		fd.players[4].String(),
		fd.players[5].String(),
		fd.players[6].String(),
		fd.players[7].String(),
		fd.players[8].String(),
		fd.players[9].String(),
	}, "\n\n")
}
