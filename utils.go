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
	elems []interface{}
}

func NewStack() *Stack {
	return &Stack{
		elems: make([]interface{}, 0),
	}
}

func (s *Stack) IsEmpty() bool {
	return len(s.elems) == 0
}

func (s *Stack) Push(x interface{}) {
	s.elems = append(s.elems, x)
}

func (s *Stack) Pop() interface{} {
	lastIndex := len(s.elems) - 1
	val := s.elems[lastIndex]
	s.elems = s.elems[:lastIndex]
	return val
}

func (s *Stack) Peek() interface{} {
	return s.elems[len(s.elems)-1]
}

func (s *Stack) Remove(x interface{}) {
	var i int
	for i = len(s.elems) - 1; s.elems[i] != x; i-- {
	}
	switch(i) {
	case 0:
		s.elems = s.elems[1:]
	case len(s.elems)-1:
		s.elems = s.elems[:i]
	default:
		s.elems = append(s.elems[:i], s.elems[i+1:])
	}
}
