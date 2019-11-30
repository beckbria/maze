package maze

// Cell represents a cell in a maze grid
type Cell interface {
	// Row returns the index of the row where the cell is located
	Row() int64
	// Column returns the index of the column where the cell is located
	Column() int64
	// Link links one cell to another bidirectionally
	Link(neighbor Cell)
	// LinkOneWay links one cell to another unidirectionally
	LinkOneWay(neighbor Cell)
	// Unlink removes the bidirectional link between two cells
	Unlink(neighbor Cell)
	// UnlinkOneWay removes the unidirectional link between a cell and its neighbor
	UnlinkOneWay(neighbor Cell)
	// Linked returns true if a cell is linked to another
	Linked(neighbor Cell) bool
	// Neighbors returns the list of direct neighbors of this cell
	Neighbors() []Cell
}

// CellBase implements the linking behavior common to all cells
type CellBase struct {
	// The location of this cell in the Grid
	row, column int64
	// A set of cells direclty linked to this cell
	links map[Cell]bool
}

func newCellBase(row, column int64) CellBase {
	c := CellBase{
		row:    row,
		column: column,
		links:  make(map[Cell]bool)}
	return c
}

// Row returns the index of the row where the cell is located
func (c *CellBase) Row() int64 {
	return c.row
}

// Column returns the index of the column where the cell is located
func (c *CellBase) Column() int64 {
	return c.column
}

// LinkOneWay links one cell to another unidirectionally
func (c *CellBase) LinkOneWay(neighbor Cell) {
	c.links[neighbor] = true
}

// Link links one cell to another bidirectionally
func (c *CellBase) Link(neighbor Cell) {
	c.LinkOneWay(neighbor)
	neighbor.LinkOneWay(c)
}

// Unlink removes the bidirectional link between two cells
func (c *CellBase) Unlink(neighbor Cell) {
	c.UnlinkOneWay(neighbor)
	neighbor.UnlinkOneWay(c)
}

// UnlinkOneWay removes the unidirectional link between a cell and its neighbor
func (c *CellBase) UnlinkOneWay(neighbor Cell) {
	delete(c.links, neighbor)
}

// Linked returns true if a cell is linked to another
func (c *CellBase) Linked(neighbor Cell) bool {
	if neighbor == nil {
		return false
	}
	linked, ok := c.links[neighbor]
	return ok && linked
}

// Neighbors returns the list of direct neighbors of this cell
func (c *CellBase) Neighbors() []Cell {
	return []Cell{}
}
