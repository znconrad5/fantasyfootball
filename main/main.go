package main

import (
	"bufio"
	"fantasyfootball"
	"flag"
	"fmt"
	"log"
	// "math"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	// "time"
)

var dataDir = flag.String("dataDir", "/Users/zachconrad/Documents/go/src/fantasyfootball/data", "The directory to look in for player statistics.")
var startWeek = flag.Int("startWeek", 1, "The week to start player statistic gathering.")
var endWeek = flag.Int("endWeek", 14, "The week to end player statistic gathering, inclusive.")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func closeChannelOnStdin(channel chan bool) {
	var temp string
	fmt.Scanf("%s", &temp)
	close(channel)
}

func main() {
	// Set up
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	data := fantasyfootball.NewDataSource(*dataDir, *startWeek, *endWeek)
	data.LoadAll()
	draft := fantasyfootball.NewFantasyDraft([]string{"Brewsers", "Legacy Losers", "Smokin Weeden", "Benson's next BWI", "Stafford Infection", "The Croakin Krogans", "Kirkland Bearhawks", "Russel Wilson Jay Cutler", "We're Going Streaking", "amit is late"}, 7, data)
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		        if err != nil {
		            log.Fatal(err)
		        }
		        pprof.StartCPUProfile(f)
		        defer pprof.StopCPUProfile()
	}
	for i:=1; !draft.IsOver(); i++ {
		stopper := make(chan bool);
		go closeChannelOnStdin(stopper)
		fmt.Printf("Suggested draft for %v\n", draft.CurrentPlayer())
		moves := draft.IterativeAlphabeta(stopper)
		for move := range moves {
			fmt.Printf("\t%v %v\n", move.Evaluation, move.Player)
		}
		var player *fantasyfootball.FootballPlayer
		var ok bool
		for ; !ok; {
			fmt.Printf("Actual Draft: ")
			playerName, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			player, ok = data.Get(strings.TrimRight(playerName, "\n"))
			fmt.Println(player)
		}	
		draft.Draft(player)
	}
	// for i:=1; i<2; i++ {
	// 	stopper := make(chan bool);
	// 	time.AfterFunc(60*time.Second, func() { close(stopper) })
	// 	moves := draft.IterativeAlphabeta(stopper)
	// 	fmt.Printf("Suggested draft for %v\n", i)
	// 	var lastMove fantasyfootball.Move
	// 	for move := range moves {
	// 		fmt.Printf("\t%v %v\n", move.Evaluation, move.Player)
	// 		lastMove = move
	// 	}
	// 	draft.Draft(lastMove.Player)
	// }
	// fmt.Println(draft.String())
	// To test profiling
	// for i:=0; i<60; i++ {
	// 	move, val, _ := draft.Alphabeta(7, math.MinInt32, math.MaxInt32, nil)
	// 	fmt.Printf("Suggested draft for %v: %v %v\n", i, val, move)
	// 	draft.Draft(move)
	// }
}
