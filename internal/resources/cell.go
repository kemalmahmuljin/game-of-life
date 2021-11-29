package resources


type Cell struct {
	state  bool
	neigs  int
}

func (c *Cell) AddNeigs(neigs int) () {
	c.neigs += neigs
}

func (c *Cell) GetState() bool {
	return c.state
}

func (c *Cell) SetState(state bool) {
	c.state = state
}