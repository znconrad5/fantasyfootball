package fantasyfootball

type Position uint8

const (
	DST Position = iota
	K
	QB
	RB
	TE
	WR
)

type FootballPlayer struct {
	name         string
	team         string
	position     Position
	points       [SEASON_LENGTH]int
	totalPoints_ int
}

func (player *FootballPlayer) totalPoints() int {
	if player.totalPoints_ == 0 {
		sum := 0
		for _, v := range player.points {
			sum += v
		}
		player.totalPoints_ = sum
	}
	return player.totalPoints_
}

type ByTotalPointsAsc struct {
	players []*FootballPlayer
}

func (s *ByTotalPointsAsc) Len() int { return len(s.players) }

func (s *ByTotalPointsAsc) Less(i, j int) bool {
	return s.players[i].totalPoints() < s.players[j].totalPoints()
}

func (s *ByTotalPointsAsc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}

type ByTotalPointsDesc struct {
	players []*FootballPlayer
}

func (s *ByTotalPointsDesc) Len() int { return len(s.players) }

func (s *ByTotalPointsDesc) Less(i, j int) bool {
	return s.players[i].totalPoints() > s.players[j].totalPoints()
}

func (s *ByTotalPointsDesc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}

type ByWeekPointsAsc struct {
	players []*FootballPlayer
	week    int
}

func (s *ByWeekPointsAsc) Len() int { return len(s.players) }

func (s *ByWeekPointsAsc) Less(i, j int) bool {
	return s.players[i].points[s.week-1] < s.players[j].points[s.week-1]
}

func (s *ByWeekPointsAsc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}

type ByWeekPointsDesc struct {
	players []*FootballPlayer
	week    int
}

func (s *ByWeekPointsDesc) Len() int { return len(s.players) }

func (s *ByWeekPointsDesc) Less(i, j int) bool {
	return s.players[i].points[s.week-1] > s.players[j].points[s.week-1]
}

func (s *ByWeekPointsDesc) Swap(i, j int) {
	s.players[i], s.players[j] = s.players[j], s.players[i]
}
