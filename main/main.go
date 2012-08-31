package main

import (
	"fantasyfootball"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	// "runtime"
	"runtime/pprof"
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
		move, val := draft.Alphabeta(7, math.MinInt32, math.MaxInt32)
		fmt.Printf("Suggested draft for %v: %v %v\n", i, move, val)
		draft.Draft(move)
	}
}
