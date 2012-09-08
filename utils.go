package fantasyfootball

import (
	"log"
)

const (
	SEASON_LENGTH = 14
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type Stack struct {
	elems []*FootballPlayer
}

func NewStack() *Stack {
	return &Stack{
		elems: make([]*FootballPlayer, 0),
	}
}

func (s *Stack) IsEmpty() bool {
	return len(s.elems) == 0
}

func (s *Stack) Push(x *FootballPlayer) {
	s.elems = append(s.elems, x)
}

func (s *Stack) Pop() *FootballPlayer {
	lastIndex := len(s.elems) - 1
	val := s.elems[lastIndex]
	s.elems = s.elems[:lastIndex]
	return val
}

func (s *Stack) Peek() *FootballPlayer {
	return s.elems[len(s.elems)-1]
}

func (s *Stack) Remove(x *FootballPlayer) {
	var i int
	for i = len(s.elems) - 1; s.elems[i] != x; i-- {
	}
	switch i {
	case 0:
		s.elems = s.elems[1:]
	case len(s.elems) - 1:
		s.elems = s.elems[:i]
	default:
		newElems := s.elems[:i]
		for _, v := range s.elems[i+1:] {
			newElems = append(newElems, v)
		}
		s.elems = newElems
	}
}
