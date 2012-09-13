package fantasyfootball

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"testing"
)

var dataSourceTestDir = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed")
var dataSourceTestStartWeek = 2
var dataSourceTestEndWeek = 14

func TestLoadDsts(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadDsts()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
	fmt.Println()
}

func TestLoadKs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadKs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
	fmt.Println()
}

func TestLoadQbs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadQbs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
	fmt.Println()
}

func TestLoadRbs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadRbs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
	fmt.Println()
}

func TestLoadTes(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadTes()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
	fmt.Println()
}

func TestLoadWrs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadWrs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
	fmt.Println()
}

func TestTimeLoadAll(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	loader.LoadAll()
	allPlayers := make([]*FootballPlayer, 0)
	for _, v := range loader.AllPlayers {
		allPlayers = append(allPlayers, v)
	}
	sort.Sort(&ByTotalPointsDesc{allPlayers})
	for _, v := range allPlayers {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
}
