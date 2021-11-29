package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"game-of-life/internal/resources"
	"game-of-life/internal/tools"

	"github.com/gdamore/tcell/v2"
)
var (
	seed int64 	= 14
	aMin       	= 1
	aMax 		= 4
	dA   		= 3
)

func main() {
	filePtr := flag.String("filename", " ", "Optional file to read in.")
	rPtr    := flag.Int("rows", 64, "Number of rows to generate. Should be divisible by sqrt(P).")
	cPtr    := flag.Int("cols", 64, "Number of rows to generate. Should be divisible by sqrt(P).")
	nGoPtr  := flag.Int("sqrtP", 1, "Square root of number of goroutines to use.")
	ItPtr   := flag.Int("nits", 150, "Number of iterations to simulate.")
	plotPtr := flag.Bool("plot", true, "True if separate iterations should be plotted.")
	delayPtr:= flag.Int("delay", 30, "Sleep time between two subsequent iterations.")
	dPtr 	:= flag.Float64("dLevel", 0.3, "Set density level for generated board float between (0,1).")

	flag.Parse()

	var density = *dPtr
	var filename = *filePtr
	var numberOfIts = *ItPtr
	var sqrtNThreads = *nGoPtr
	var drawIterations = *plotPtr

	// Verify that the number of goroutines is valid for the world dimensions
	if !(*rPtr%sqrtNThreads == 0 && *cPtr%sqrtNThreads == 0) {
		err := errors.New("number of columns and rows should be a multiple of sqrt(P)")
		log.Fatal(err)
	}

	// Randomly generate board when no specific board is requested
	if filename == " " {
		filename = tools.GenerateRandomBoard(*rPtr, *cPtr, density, seed)
	}

	// Setup screen
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	s := tools.InitializeScreen()
	go tools.CancelRoutine(s)

	var world resources.World
	world.ReadWorld(filename, sqrtNThreads, aMin, aMax, dA)

	if drawIterations { // Only draw if specified
		tools.UpdateScreen(s, &world, boxStyle)
		time.Sleep(time.Duration(*delayPtr) * time.Millisecond)
	}

	total := time.Duration(0)
	var start time.Time

	// Iterate over world
	for i:=0; i < numberOfIts; i++ {
		start = time.Now()
		world.DoStep()
		total += time.Now().Sub(start)
		if drawIterations { // Only draw if specified
			tools.UpdateScreen(s, &world, boxStyle)
			time.Sleep(time.Duration(*delayPtr) * time.Millisecond)
		}
	}

	// Terminate screen
	s.Fini()

	// Print performance results
	fmt.Printf("Finished %d iterations in %s\n", numberOfIts, total)
	fmt.Printf("Average is %0.3f ms / It\n", (float64(total)/float64(numberOfIts))/1000000.0)
}
