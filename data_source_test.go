package fantasyfootball

import (
	"fmt"
	"os"
	"sort"
	"testing"
)

var loader = NewFileDataSource(os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed"), 5, 14)

func TestLoadDsts(t *testing.T) {
	sort.Sort(&ByTotalPointsDesc{loader.DefenseSpecialTeams()})
	for i, v := range loader.DefenseSpecialTeams() {
		if i >= 32 {
			break
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Println()
}

func TestLoadKs(t *testing.T) {
	sort.Sort(&ByTotalPointsDesc{loader.Kickers()})
	for i, v := range loader.Kickers() {
		if i >= 32 {
			break
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Println()
}

func TestLoadQbs(t *testing.T) {
	sort.Sort(&ByTotalPointsDesc{loader.Quarterbacks()})
	for i, v := range loader.Quarterbacks() {
		if i >= 32 {
			break
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Println()
}

func TestLoadRbs(t *testing.T) {
	sort.Sort(&ByTotalPointsDesc{loader.RunningBacks()})
	for i, v := range loader.RunningBacks() {
		if i >= 100 {
			break
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Println()
}

func TestLoadTes(t *testing.T) {
	sort.Sort(&ByTotalPointsDesc{loader.TightEnds()})
	for i, v := range loader.TightEnds() {
		if i >= 32 {
			break
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Println()
}

func TestLoadWrs(t *testing.T) {
	sort.Sort(&ByTotalPointsDesc{loader.WideReceivers()})
	for i, v := range loader.WideReceivers() {
		if i >= 100 {
			break
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Println()
}

func TestTimeLoadAll(t *testing.T) {
	allPlayers := make([]*FootballPlayer, 0)
	for _, v := range loader.AllPlayers() {
		allPlayers = append(allPlayers, v)
	}
	sort.Sort(&ByTotalPointsDesc{allPlayers})
	for i, v := range allPlayers {
		if i >= 150 {
			break
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
}
