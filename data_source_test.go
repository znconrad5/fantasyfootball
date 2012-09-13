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
<<<<<<< HEAD
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadDsts()
	for i, v := range players {
		if i>=32 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
=======
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadDsts()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
>>>>>>> 7c391b1a565999696b1eedbc8ecf0d1932e014a9
	fmt.Println()
}

func TestLoadKs(t *testing.T) {
<<<<<<< HEAD
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadKs()
	for i, v := range players {
		if i>=32 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
=======
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadKs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
>>>>>>> 7c391b1a565999696b1eedbc8ecf0d1932e014a9
	fmt.Println()
}

func TestLoadQbs(t *testing.T) {
<<<<<<< HEAD
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadQbs()
	for i, v := range players {
		if i>=32 {
			break;
		}
			fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
=======
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadQbs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
>>>>>>> 7c391b1a565999696b1eedbc8ecf0d1932e014a9
	fmt.Println()
}

func TestLoadRbs(t *testing.T) {
<<<<<<< HEAD
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadRbs()
	for i, v := range players {
		if i>=100 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
=======
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadRbs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
>>>>>>> 7c391b1a565999696b1eedbc8ecf0d1932e014a9
	fmt.Println()
}

func TestLoadTes(t *testing.T) {
<<<<<<< HEAD
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadTes()
	for i, v := range players {
		if i>=32 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
=======
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadTes()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
>>>>>>> 7c391b1a565999696b1eedbc8ecf0d1932e014a9
	fmt.Println()
}

func TestLoadWrs(t *testing.T) {
<<<<<<< HEAD
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadWrs()
	for i, v := range players {
		if i>=100 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
=======
	loader := NewDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players, defaultPlayer := loader.loadWrs()
	for _, v := range players {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
	fmt.Printf("%s (%s) %d %d\n", defaultPlayer.Name, defaultPlayer.Team, defaultPlayer.TotalPoints(), defaultPlayer.Points)
>>>>>>> 7c391b1a565999696b1eedbc8ecf0d1932e014a9
	fmt.Println()
}

func TestTimeLoadAll(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	allPlayers := make([]*FootballPlayer, 0)
	for _, v := range loader.AllPlayers {
		allPlayers = append(allPlayers, v)
	}
	sort.Sort(&ByTotalPointsDesc{allPlayers})
	for _, v := range allPlayers {
		fmt.Printf("%s (%s) %d %d\n", v.Name, v.Team, v.TotalPoints(), v.Points)
	}
}
