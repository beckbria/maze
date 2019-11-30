package maze

// RectangleCell represents a cell in a rectangular maze
type RectangleCell struct {
	// The immediate neighbors of this cell
	North, South, East, West *RectangleCell

	base CellBase
}

// NewRectangleCell creates a new RectangleCell
func NewRectangleCell(row, column int64) RectangleCell {
	c := RectangleCell{base: newCellBase(row, column)}
	return c
}

// Neighbors returns the list of direct neighbors of this cell
func (c *RectangleCell) Neighbors() []Cell {
	ret := []Cell{}
	for _, n := range []*RectangleCell{c.North, c.South, c.East, c.West} {
		if n != nil {
			ret = append(ret, n)
		}
	}
	return ret
}

// Row returns the index of the row where the cell is located
func (c *RectangleCell) Row() int64 {
	return c.base.Row()
}

// Column returns the index of the column where the cell is located
func (c *RectangleCell) Column() int64 {
	return c.base.Column()
}

// LinkOneWay links one cell to another unidirectionally
func (c *RectangleCell) LinkOneWay(neighbor Cell) {
	c.base.LinkOneWay(neighbor)
}

// Link links one cell to another bidirectionally
func (c *RectangleCell) Link(neighbor Cell) {
	c.base.Link(neighbor)
}

// Unlink removes the bidirectional link between two cells
func (c *RectangleCell) Unlink(neighbor Cell) {
	c.base.Unlink(neighbor)
}

// UnlinkOneWay removes the unidirectional link between a cell and its neighbor
func (c *RectangleCell) UnlinkOneWay(neighbor Cell) {
	c.base.UnlinkOneWay(neighbor)
}

// Linked returns true if a cell is linked to another
func (c *RectangleCell) Linked(neighbor Cell) bool {
	return c.base.Linked(neighbor)
}
