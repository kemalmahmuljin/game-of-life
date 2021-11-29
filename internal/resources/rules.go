package resources

type Rules struct {
	aMin, aMax, dA 	int
}

func (r *Rules) golLogic(cell Cell) int {
	if cell.state {
		if cell.neigs >= r.aMax || cell.neigs <= r.aMin {
			return -1
		}
	} else {
		if cell.neigs == r.dA {
			return 1
		}
	}
	return 0
}
