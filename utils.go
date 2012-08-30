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
