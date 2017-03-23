package main

// Cell - This is used by the solver to hold information gathered about a
//        specific cell in the playing field to help determine if it is a
//        safe move or a mine
type Cell struct {
	/* weight - Holds probability weight for risky moves
	   isMine - -1 = unknown, 0 = not mine, 1 = is mine */
	x, y, val, covNeighbors, nearbyMines, weight, isMine int
}

// NewCell - constructor
func NewCell(value int) *Cell {
	c := new(Cell)
	c.isMine = -1
	c.val = value

	return c
}

// SetValue - Method to change the value of a cell and reset weight/isMine
func (c *Cell) SetValue(value int) {
	c.val = value
	// Reset weight and isMine state in the process
	c.weight = 0
	c.isMine = -1
}
