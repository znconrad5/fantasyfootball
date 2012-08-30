package fantasyfootball

import (
	"container/stack"
	"sort"
)

type FantasyDraft struct {
	players [10]*FantasyPlayer
	maxPlayer *FantasyPlayer
	playersDrafted int
	dsts *stack.Stack	// +7
	ks   *stack.Stack	// +8
	qbs  *stack.Stack	// +30
	rbs  *stack.Stack	// +27
	tes  *stack.Stack	// +22
	wrs  *stack.Stack	// +16
}

func NewFantasyDraft(names [10]string, maxName string, data *DataSource) *FantasyDraft {
	fd := &FantasyDraft{playersDrafted:0}
	for i, name := range names {
		fd.players[i] = newFantasyPlayer(name, data)
		if name == maxName {
			fd.maxPlayer = fd.players[i]
		}
	}
	
	fd.dsts = stack.New()
	sort.Sort(&ByTotalPointsAsc{data.dsts})
	for _, v := range data.dsts {
		fd.dsts.Push(v)
	}
	fd.ks = stack.New()
	sort.Sort(&ByTotalPointsAsc{data.ks})
	for _, v := range data.ks {
		fd.ks.Push(v)
	}
	fd.qbs = stack.New()
	sort.Sort(&ByTotalPointsAsc{data.qbs})
	for _, v := range data.qbs {
		fd.qbs.Push(v)
	}
	fd.rbs = stack.New()
	sort.Sort(&ByTotalPointsAsc{data.rbs})
	for _, v := range data.rbs {
		fd.rbs.Push(v)
	}
	fd.tes = stack.New()
	sort.Sort(&ByTotalPointsAsc{data.tes})
	for _, v := range data.tes {
		fd.tes.Push(v)
	}
	fd.wrs = stack.New()
	sort.Sort(&ByTotalPointsAsc{data.wrs})
	for _, v := range data.wrs {
		fd.wrs.Push(v)
	}
	return fd
}

func (fd *FantasyDraft) currentPlayer() *FantasyPlayer {
	round := fd.playersDrafted / 10
	offset := fd.playersDrafted % 10
	if round % 2 == 0 {
		return fd.players[offset]
	}
	return fd.players[len(fd.players)-1-offset]
}

func (fd *FantasyDraft) Draft(draftee *FootballPlayer) {
	fd.removeFootballPlayer(draftee)
	fd.currentPlayer().draft(draftee)
	fd.playersDrafted++
}

func (fd *FantasyDraft) removeFootballPlayer(player *FootballPlayer) {
	var pool *stack.Stack
	switch(player.position) {
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

func (fd *FantasyDraft) Alphabeta(depth, alpha, beta int) (*FootballPlayer, int) {
	var move *FootballPlayer
	if depth == 0 || fd.isDraftOver() {
		var value int
		if fd.playersDrafted > 90 {
			value = fd.evaluate()
		} else {
			value = fd.estimate()
		}
		return move, value
	}
	currentPlayer := fd.currentPlayer()
	s := &ByBestLikelyMove{currentPlayer, [...]*stack.Stack{fd.qbs, fd.rbs, fd.tes, fd.wrs, fd.ks, fd.dsts}}
	sort.Sort(s)
	if (fd.maxPlayer == currentPlayer) {
		for _, v := range s.stacks {
			draftee := v.Pop().(*FootballPlayer)
			currentPlayer.draft(draftee)
			fd.playersDrafted++
			_, moveValue := fd.Alphabeta(depth-1, alpha, beta)
			if (moveValue > alpha) {
				alpha = moveValue
				move = draftee
			}
			fd.playersDrafted--
			currentPlayer.undraft(draftee)
			v.Push(draftee)
			if beta <= alpha {
				break
			}
		}
		return move, alpha
	} else {
		for _, v := range s.stacks {
			draftee := v.Pop().(*FootballPlayer)
			currentPlayer.draft(draftee)
			fd.playersDrafted++
			_, moveValue := fd.Alphabeta(depth-1, alpha, beta)
			if (moveValue < beta) {
				beta = moveValue
				move = draftee
			}
			fd.playersDrafted--
			currentPlayer.undraft(draftee)
			v.Push(draftee)
			if beta <= alpha {
				break
			}
		}
		return move, beta
	}
	panic("alphabeta() unexpectedly returned")
}

type ByBestLikelyMove struct {
	currentPlayer *FantasyPlayer
	stacks [6]*stack.Stack
}

func (s *ByBestLikelyMove) Len() int { return len(s.stacks) }

func (s *ByBestLikelyMove) Less(i, j int) bool {
	iDraftee := s.stacks[i].Peek().(*FootballPlayer)
	s.currentPlayer.draft(iDraftee)
	iTotal := s.currentPlayer.estimateTotalPoints()
	s.currentPlayer.undraft(iDraftee)
	
	jDraftee := s.stacks[j].Peek().(*FootballPlayer)
	s.currentPlayer.draft(jDraftee)
	jTotal := s.currentPlayer.estimateTotalPoints()
	s.currentPlayer.undraft(jDraftee)
	
	return iTotal > jTotal
}

func (sort *ByBestLikelyMove) Swap(i, j int) {
	sort.stacks[i], sort.stacks[j] = sort.stacks[j], sort.stacks[i]
}

func (fd *FantasyDraft) evaluate() int {
	opponentTotal := 0
	for _, player := range fd.players {
		opponentTotal += player.totalPoints()
	}
	return fd.maxPlayer.totalPoints() * len(fd.players) - opponentTotal
}

func (fd *FantasyDraft) estimate() int {
	opponentTotal := 0
	for _, player := range fd.players {
		opponentTotal += player.estimateTotalPoints()
	}
	return fd.maxPlayer.estimateTotalPoints() * len(fd.players) - opponentTotal
}

func (fd *FantasyDraft) isDraftOver() bool {
	return fd.playersDrafted == 150
}