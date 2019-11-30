package maze

// Cell represents a cell in a rectangular maze
type Cell struct {
	// The location of this cell in the Grid
	Row, Column int64
	// The immediate neighbors of this cell
	North, South, East, West *Cell
	// A set of cells direclty linked to this cell
	links map[*Cell]bool
}

func NewCell(row, column int64) Cell {
	c := Cell{
		Row:    row,
		Column: column,
		links:  make(map[*Cell]bool)}
	return c
}

// LinkOneWay links one cell to another unidirectionally
func (c *Cell) LinkOneWay(neighbor *Cell) {
	c.links[neighbor] = true
}

// Link links one cell to another bidirectionally
func (c *Cell) Link(neighbor *Cell) {
	c.LinkOneWay(neighbor)
	neighbor.LinkOneWay(c)
}

// Unlink removes the bidirectional link between two cells
func (c *Cell) Unlink(neighbor *Cell) {
	c.UnlinkOneWay(neighbor)
	neighbor.UnlinkOneWay(c)
}

// UnlinkOneWay removes the unidirectional link between a cell and its neighbor
func (c *Cell) UnlinkOneWay(neighbor *Cell) {
	delete(c.links, neighbor)
}

// Linked returns true if a cell is linked to another
func (c *Cell) Linked(neighbor *Cell) bool {
	if neighbor == nil {
		return false
	}
	linked, ok := c.links[neighbor]
	return ok && linked
}

// Neighbors returns the list of direct neighbors of this cell
func (c *Cell) Neighbors() []*Cell {
	ret := []*Cell{}
	for _, n := range []*Cell{c.North, c.South, c.East, c.West} {
		if n != nil {
			ret = append(ret, n)
		}
	}
	return ret
}
