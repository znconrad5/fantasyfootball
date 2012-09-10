package fantasyfootball

import (
	"fmt"
	"runtime"
	"sort"
	"testing"
)

var dataSourceTestDir = "C:/Users/Dustin/Documents/golibs/src/github.com/znconrad5/fantasyfootball/parsed"

func TestLoadDsts(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, 1, 14)
	players, defaultPlayer := loader.loadDsts()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.name, defaultPlayer.team, defaultPlayer.totalPoints(), defaultPlayer.points)
	fmt.Println()
}

func TestLoadKs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, 1, 14)
	players, defaultPlayer := loader.loadKs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.name, defaultPlayer.team, defaultPlayer.totalPoints(), defaultPlayer.points)
	fmt.Println()
}

func TestLoadQbs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, 1, 14)
	players, defaultPlayer := loader.loadQbs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.name, defaultPlayer.team, defaultPlayer.totalPoints(), defaultPlayer.points)
	fmt.Println()
}

func TestLoadRbs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, 1, 14)
	players, defaultPlayer := loader.loadRbs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.name, defaultPlayer.team, defaultPlayer.totalPoints(), defaultPlayer.points)
	fmt.Println()
}

func TestLoadTes(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, 1, 14)
	players, defaultPlayer := loader.loadTes()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.name, defaultPlayer.team, defaultPlayer.totalPoints(), defaultPlayer.points)
	fmt.Println()
}

func TestLoadWrs(t *testing.T) {
	loader := NewDataSource(dataSourceTestDir, 1, 14)
	players, defaultPlayer := loader.loadWrs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.name, defaultPlayer.team, defaultPlayer.totalPoints(), defaultPlayer.points)
	fmt.Println()
}

func TestTimeLoadAll(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	loader := NewDataSource(dataSourceTestDir, 1, 14)
	loader.LoadAll()
	allPlayers := make([]*FootballPlayer, 0)
	for _, v := range loader.allPlayers {
		allPlayers = append(allPlayers, v)
	}
	sort.Sort(&ByTotalPointsDesc{allPlayers})
	for _, v := range allPlayers {
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
}