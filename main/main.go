package main

import (
	"fantasyfootball"
	"flag"
	"fmt"
	"log"
	"os"
	// "runtime"
	"runtime/pprof"
	"time"
)

var dataDir = flag.String("dataDir", "/Users/zachconrad/Documents/go/src/fantasyfootball/data", "The directory to look in for player statistics.")
var startWeek = flag.Int("startWeek", 1, "The week to start player statistic gathering.")
var endWeek = flag.Int("endWeek", 3, "The week to end player statistic gathering, inclusive.")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// Set up
	// runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	data := fantasyfootball.NewDataSource(*dataDir, *startWeek, *endWeek)
	data.LoadAll()
	draft := fantasyfootball.NewFantasyDraft([...]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, "1", data)
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		        if err != nil {
		            log.Fatal(err)
		        }
		        pprof.StartCPUProfile(f)
		        defer pprof.StopCPUProfile()
	}
	for i:=0; i<20; i++ {
		stopper := make(chan bool);
		time.AfterFunc(5*time.Second, func() { close(stopper) })
		moves := draft.IterativeAlphabeta(stopper)
		fmt.Printf("Suggested draft for %v\n", i)
		var lastMove fantasyfootball.Move
		for move := range moves {
			fmt.Printf("\t%v %v\n", move.Evaluation, move.Player)
			lastMove = move
		}
		draft.Draft(lastMove.Player)
	}
}
