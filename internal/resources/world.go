package resources

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var directions = [8][2]int{ {-1,-1}, {-1,0}, {-1,1}, {0,-1}, {0,1}, {1,-1}, {1,0}, {1,1} }

type World struct {
	cells 			[][]Cell
	cellgroups  	[]CellGroup
	cWidth, cHeight int 			// cell group dimensions
	sqrtNThreads 	int
	R            	Rules
}

func (w *World) DoStep() {
	var nThreads = w.sqrtNThreads * w.sqrtNThreads

	finished := make(chan int, nThreads)

	// Start of concurrent threads
	for id := 0; id < nThreads; id++ {
		go w.parallelUpdate(id, finished)
	}

	// Write world sections to world board as goroutines finish
	// Wait for all goroutines to finish
	for id := 0; id < nThreads; id++ {
		w.copyFromCellGroup(<- finished)
	}
}

func (w *World) Init() {
	var row, col, r1, r2, c1, c2 int
	for id := 0; id < w.sqrtNThreads*w.sqrtNThreads; id++ {
		row, col = id/w.sqrtNThreads, id%w.sqrtNThreads
		r1 = row*w.cHeight
		r2 = (row+1)*w.cHeight
		c1 = col*w.cWidth
		c2 = (col+1)*w.cWidth

		for i := 0; i < r2-r1; i++ {
			for j := 0; j < c2-c1; j++ {
				w.cellgroups[id].cells[i][j] = w.cells[r1+i][c1+j]
			}
		}
	}
}

func (w *World) GetCellState(r, c int) bool {
	return w.cells[r][c].GetState()
}

func (w *World) GetWidth() int {
	return w.cWidth*w.sqrtNThreads
}

func (w *World) GetHeight() int {
	return w.cHeight*w.sqrtNThreads
}

func (w *World) ReadWorld(filename string, sqrtNThreads, aMin, aMax, dA int) {
	// Read in a world and process line by line
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")

	var nThreads = sqrtNThreads * sqrtNThreads
	var temp []string
	var i1, i2 int

	w.sqrtNThreads = sqrtNThreads
	w.cellgroups = make([]CellGroup, nThreads)
	for row, cell := range lines {
		temp = strings.Split(cell, ",")
		if len(temp) != 2 {
			break
		}
		i1, _ = strconv.Atoi(temp[0])
		i2, _ = strconv.Atoi(temp[1])

		if row == 0 {
			// Set dimensions of world
			w.cHeight = i1 / sqrtNThreads
			w.cWidth = i2 / sqrtNThreads

			w.cells = make([][]Cell, i1)
			for i:= range w.cells {
				w.cells[i] = make([]Cell, i2)
			}
		} else {
			// Set alive cells of world
			w.cells[i1][i2].SetState(true)
			w.addToNeighbour(i1, i2)
		}
	}

	// Allocate and initialize cellgroups
	for i := 0; i < nThreads; i++ {
		tmp := new(CellGroup)
		tmp.cells = make([][]Cell, w.cHeight)
		for j := range tmp.cells {
			tmp.cells[j] = make([]Cell, w.cWidth)
		}

		tmp.NBar = make(chan []int, 1)
		tmp.EBar = make(chan []int, 1)
		tmp.SBar = make(chan []int, 1)
		tmp.WBar = make(chan []int, 1)

		tmp.NW = make(chan int, 1)
		tmp.NE = make(chan int, 1)
		tmp.SE = make(chan int, 1)
		tmp.SW = make(chan int, 1)
		w.cellgroups[i] = *tmp
	}
	w.Init()

	// Set world rules
	w.R.aMin = aMin
	w.R.aMax = aMax
	w.R.dA = dA
}

func (w *World) copyFromCellGroup(id int) {
	r := (id / w.sqrtNThreads) * w.cHeight
	c := (id % w.sqrtNThreads) * w.cWidth

	for i := 0; i < w.cHeight; i++ {
		for j:=0; j < w.cWidth; j++ {
			w.cells[r+i][c+j] = w.cellgroups[id].cells[i][j]
		}
	}
}

func (w* World) parallelUpdate(id int, finished chan<- int) {
	r := id / w.sqrtNThreads
	c := id % w.sqrtNThreads

	worldSection := w.getSlice(r*w.cHeight, (r+1)*w.cHeight, c*w.cWidth, (c+1)*w.cWidth)

	// Compute world section border and info to communicate
	boundary := w.cellgroups[id].computeBoundaries(worldSection, w.R)

	// Send boundary information of world section to neighbouring world sections
	w.sendToNeighbours(id, boundary)

	// Perform world section internal computation
	w.cellgroups[id].computeInternal(worldSection, w.R)

	// Update world section border information using neighbours information
	w.cellgroups[id].updateBoundaries()

	// Fill channel to signal task end
	finished <- id
}

func (w *World) getSlice(r1, r2, c1, c2 int) (slice [][]Cell) {
	slice = make([][]Cell, r2-r1, r2-r1)
	for i := 0; i < r2-r1; i++ {
		slice[i] = make([]Cell, c2-c1, c2-c1)
		slice[i] = w.cells[r1+i][c1:c2]
	}
	return slice
}

func (w *World) sendToNeighbours(id int, b *GroupBoundary) {

	// Retrieve channels to send to depending on own id
	P := w.sqrtNThreads
	P2 := P*P

	east := (id + P + 1)%P + (id/P)*P
	west := (id + P - 1)%P + (id/P)*P

	N := w.cellgroups[(id + P2 - P)%P2].SBar
	E := w.cellgroups[east].WBar
	S := w.cellgroups[(id + P)%P2].NBar
	W := w.cellgroups[west].EBar

	NW := w.cellgroups[(west + P2 - P)%P2].SE
	NE := w.cellgroups[(east + P2 - P)%P2].SW
	SE := w.cellgroups[(east + P)%P2].NW
	SW := w.cellgroups[(west + P)%P2].NE

	// Send boundary information to neighbours
	N <- b.N
	E <- b.E
	S <- b.S
	W <- b.W

	NW <- b.NW
	NE <- b.NE
	SE <- b.SE
	SW <- b.SW
}

func (w *World) addToNeighbour(r, c int) {
	var row, col int
	height := len(w.cells)
	width := len(w.cells[0])
	for _, dir := range directions {
		row = (dir[0] + r + height) % height
		col = (dir[1] + c + width) % width
		w.cells[row][col].AddNeigs(1)
	}
}
