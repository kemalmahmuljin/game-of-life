package resources

type GroupBoundary struct {
	N, E, S, W 		[]int
	NW, NE, SE, SW 	int
}

type CellGroup struct {
	cells                  [][]Cell
	NBar, EBar, SBar, WBar chan []int
	NW, NE, SE, SW         chan int
}

func (cg *CellGroup) computeInternal(globalSlice [][]Cell, R Rules) {
	var temp int
	// Apply GoL rules to internal part of cellgroup
	for r, v := range globalSlice {
		if r > 0 && r < len(globalSlice)-1 {
			for c, w := range v {
				if c > 0 && c < len(v)-1 {
					// GoL logic
					temp = R.golLogic(w) // -1 -> Alive to dead | +1 -> Dead to alive
					if temp != 0 {
						cg.cells[r][c].SetState(temp==1)
						cg.updateNeighbours(r, c, temp)
					}
				}
			}
		}
	}
}

func (cg *CellGroup) updateNeighbours(r, c, update int) {
	for _, neig := range directions {
		cg.cells[r + neig[0]][c +neig[1]].AddNeigs(update)
	}
}

func (cg *CellGroup) updateBoundaries() {
	S := <-cg.SBar
	for r := range cg.cells[len(cg.cells)-1] {
		if S[r] != 0 {
			cg.cells[len(cg.cells)-1][r].AddNeigs(S[r])
		}
	}
	W := <-cg.WBar
	for c := range cg.cells {
		if W[c] != 0 {
			cg.cells[c][0].AddNeigs(W[c])
		}
	}
	N := <-cg.NBar
	for r := range cg.cells[0] {
		if N[r] != 0 {
			cg.cells[0][r].AddNeigs(N[r])
		}
	}
	E := <-cg.EBar
	for c := range cg.cells {
		if E[c] != 0 {
			cg.cells[c][len(cg.cells[0])-1].AddNeigs(E[c])
		}
	}

	NE := <-cg.NE
	SE := <-cg.SE
	SW := <-cg.SW
	NW := <-cg.NW

	cg.cells[0][len(cg.cells[0])-1].AddNeigs(NE)
	cg.cells[len(cg.cells)-1][len(cg.cells[0])-1].AddNeigs(SE)
	cg.cells[len(cg.cells)-1][0].AddNeigs(SW)
	cg.cells[0][0].AddNeigs(NW)
}

func (cg *CellGroup) computeBoundaries(globalSlice [][]Cell, R Rules) (b *GroupBoundary) {
	var temp int
	b = &GroupBoundary{
		N: make([]int, len(globalSlice[0])),
		E: make([]int, len(globalSlice)),
		S: make([]int, len(globalSlice[0])),
		W: make([]int, len(globalSlice)),
	}
	// N (NW + NE)
	for c, w := range globalSlice[0] {
		temp = R.golLogic(w) // -1 -> Alive to dead | +1 -> Dead to alive
		if temp != 0 {
			cg.cells[0][c].SetState(temp==1)
			cg.updateBoundaryNeigs(0, c, temp, b)
		}
	}
	// S (SW + SE)
	for c, w := range globalSlice[len(globalSlice)-1] {
		temp = R.golLogic(w)
		if temp != 0 {
			cg.cells[len(globalSlice)-1][c].SetState(temp==1)
			cg.updateBoundaryNeigs(len(globalSlice)-1, c, temp, b)
		}
	}

	for r, w := range globalSlice {
		if r > 0 && r < len(globalSlice)-1 {
			// E
			temp = R.golLogic(w[len(w)-1])
			if temp != 0 {
				cg.cells[r][len(w)-1].SetState(temp==1)
				cg.updateBoundaryNeigs(r, len(w)-1, temp, b)
			}
			// W
			temp = R.golLogic(w[0])
			if temp != 0 {
				cg.cells[r][0].SetState(temp==1)
				cg.updateBoundaryNeigs(r, 0, temp, b)
			}
		}
	}
	return b
}


func (cg *CellGroup) updateBoundaryNeigs(r, c int, update int, b *GroupBoundary) {
	var row, col int
	for _, neig := range directions {
		row = r + neig[0]
		col = c + neig[1]
		if row >= 0 && row < len(cg.cells) && col >= 0 && col < len(cg.cells) { 	// internal
			cg.cells[row][col].AddNeigs(update)
		} else if row == -1 && col == -1 { 											// NW
			b.NW += update
		} else if row == -1 && col >= 0 && col < len(cg.cells) { 					// N
			b.N[col] += update
		} else if row == -1 && col == len(cg.cells) { 								// NE
			b.NE += update
		} else if row >= 0 && row < len(cg.cells) && col == len(cg.cells) { 		// E
			b.E[row] += update
		} else if row == len(cg.cells) && col == len(cg.cells) { 					// SE
			b.SE += update
		} else if row == len(cg.cells) && col >= 0 && col < len(cg.cells) { 		// S
			b.S[col] += update
		} else if row == len(cg.cells) && col == -1 { 								// SW
			b.SW += update
		} else if row >= 0 && row < len(cg.cells) && col == -1 { 					// W
			b.W[row] += update
		}
	}
}
