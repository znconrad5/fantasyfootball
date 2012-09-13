package fantasyfootball

import (
	"fmt"
	"os"
	"testing"
)

var dataSourceTestDir = os.ExpandEnv("$GOPATH/src/github.com/znconrad5/fantasyfootball/parsed")
var dataSourceTestStartWeek = 2
var dataSourceTestEndWeek = 14

func TestLoadDsts(t *testing.T) {
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadDsts()
	for i, v := range players {
		if i>=32 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Println()
}

func TestLoadKs(t *testing.T) {
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadKs()
	for i, v := range players {
		if i>=32 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Println()
}

func TestLoadQbs(t *testing.T) {
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadQbs()
	for i, v := range players {
		if i>=32 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Println()
}

func TestLoadRbs(t *testing.T) {
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadRbs()
	for i, v := range players {
		if i>=100 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Println()
}

func TestLoadTes(t *testing.T) {
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadTes()
	for i, v := range players {
		if i>=32 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Println()
}

func TestLoadWrs(t *testing.T) {
	loader := NewFileDataSource(dataSourceTestDir, dataSourceTestStartWeek, dataSourceTestEndWeek)
	players := loader.loadWrs()
	for i, v := range players {
		if i>=100 {
			break;
		}
		fmt.Printf("%s (%s) %d %d\n", v.name, v.team, v.totalPoints(), v.points)
	}
	fmt.Println()
}
