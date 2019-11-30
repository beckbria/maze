package maze

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	debug = false
)

// Grid represents a rectangular maze grid
type Grid struct {
	// Rows and Columns indicate the size of the grid
	Rows, Columns int64
	// The cells in the grid
	grid [][]*Cell
}

// NewGrid creates a new rectangular grid with all cells connected to their neighbors
func NewGrid(rows, columns int64) Grid {
	if rows < 0 || columns < 0 {
		log.Fatalf("Grid dimensions invalid: [%d, %d]", rows, columns)
	}
	g := Grid{
		Rows:    rows,
		Columns: columns,
		grid:    make([][]*Cell, rows)}
	g.prepareGrid()
	g.configureCells()
	return g
}

// At accesses a cell from the grid
func (g *Grid) At(row, column int64) *Cell {
	if row < 0 || column < 0 || row >= g.Rows || column >= g.Columns {
		return nil
	}
	return g.grid[row][column]
}

// prepareGrid creates the cells in the grid
func (g *Grid) prepareGrid() {
	for r := int64(0); r < g.Rows; r++ {
		g.grid[r] = make([]*Cell, g.Columns)
		for c := int64(0); c < g.Columns; c++ {
			cell := NewCell(r, c)
			g.grid[r][c] = &cell
		}
	}
}

// configureCells establishes links between cells and their neighbors
func (g *Grid) configureCells() {
	for cell := range g.AllCells() {
		cell.North = g.At(cell.Row-1, cell.Column)
		cell.South = g.At(cell.Row+1, cell.Column)
		cell.West = g.At(cell.Row, cell.Column-1)
		cell.East = g.At(cell.Row, cell.Column+1)
	}
}

// AllRows returns a row of cells in the grid at a time
func (g *Grid) AllRows() <-chan []*Cell {
	c := make(chan []*Cell)
	go func() {
		for _, row := range g.grid {
			c <- row
		}
		close(c)
	}()
	return c
}

// AllCells iterates over all of the cells in the grid
func (g *Grid) AllCells() <-chan *Cell {
	c := make(chan *Cell)
	go func() {
		for _, row := range g.grid {
			for _, cell := range row {
				c <- cell
			}
		}
		close(c)
	}()
	return c
}

// RandomCell returns a random cell from the grid
func (g *Grid) RandomCell() *Cell {
	return g.At(rand.Int63n(g.Rows), rand.Int63n(g.Columns))
}

// Size returns the number of cells in the grid
func (g *Grid) Size() int64 {
	return g.Rows * g.Columns
}

// Unicode light box drawing characters
const (
	horizontal      = '\u2500' // ─
	vertical        = '\u2502' // │
	cornerDownRight = '\u250c' // ┌
	cornerDownLeft  = '\u2510' // ┐
	cornerUpRight   = '\u2514' // └
	cornerUpLeft    = '\u2518' // ┘
	verticalRight   = '\u251c' // ├
	verticalLeft    = '\u2524' // ┤
	horizontalDown  = '\u252c' // ┬
	horizontalUp    = '\u2534' // ┴
	intersection    = '\u253c' // ┼
)

// ToString creates a textual representation of the maze grid
func (g *Grid) ToString() string {
	return g.toString(3, 1)
}

// toString creates a textual representation of the maze grid
func (g *Grid) toString(horizontalSize, verticalSize int) string {
	if (horizontalSize < 1) || (verticalSize < 1) {
		log.Fatalf("Invalid grid size for toString: [%d, %d]", horizontalSize, verticalSize)
	}

	output := ""

	// When drawing a horizontal line across cells, we use several horizontal glyphs in a row
	horizontalLine := ""
	horizontalSpace := ""
	for i := 0; i < horizontalSize; i++ {
		horizontalLine += string(horizontal)
		horizontalSpace += " "
	}

	// Render all of the cells.  Loop inclusive of the column count to get the bottom edge
	for r := int64(0); r <= g.Rows; r++ {
		// Generate the representation of this row
		topEdge := "" // The horizontal lines between cells
		area := ""    // The contents of the cells
		// Loop inclusive of the column count to get the right edge
		for c := int64(0); c <= g.Columns; c++ {
			cell := g.At(r, c)
			if debug {
				fmt.Printf("\tC[%d,%d]: {", r, c)
				if cell == nil {
					fmt.Print("nil")
				} else {
					if cell.Linked(cell.North) {
						fmt.Print("N")
					}
					if cell.Linked(cell.East) {
						fmt.Print("E")
					}
					if cell.Linked(cell.West) {
						fmt.Print("W")
					}
					if cell.Linked(cell.South) {
						fmt.Print("S")
					}
				}
				fmt.Print("}")
			}

			ul := g.upperLeftCornerGlyph(r, c)
			topEdge += string(ul)
			// Ensure that this cell is actually linked to the cell above it.  Otherwise the following shape
			// would detect a horizontal corner glyph and incorrectly draw the line indicated by XXX
			// +-------+
			// |       |
			// +----XXX|
			// |       |
			// +-------+
			if pointsRight(ul) && (cell == nil || !cell.Linked(cell.North)) {
				topEdge += horizontalLine
			} else {
				topEdge += horizontalSpace
			}

			if pointsDown(ul) && (cell == nil || !cell.Linked(cell.West)) {
				area += string(vertical)
			} else {
				area += " "
			}
			area += horizontalSpace
		}
		if debug {
			fmt.Print("\n")
		}

		// Append this row to the complete output
		output = output + topEdge + "\n"
		if r < g.Rows {
			for i := 0; i < verticalSize; i++ {
				output = output + area + "\n"
			}
		}
	}

	return output
}

// upperLeftCornerGlyph returns the glyph which should be shown at the
// upper-left corner of a cell
func (g *Grid) upperLeftCornerGlyph(row, column int64) rune {
	// We care about four cells.  Those cells are located (relative to this glyph)
	// to the upper-left, upper-right, lower-left, and lower-right.  The row and
	// column parameters correspond to the cell to the lower-right of the glyph
	ul := g.At(row-1, column-1)
	ur := g.At(row-1, column)
	ll := g.At(row, column-1)
	lr := g.At(row, column)

	// Determine which directions this glyph needs to face
	up := false
	down := false
	left := false
	right := false

	if ul != nil {
		up = (ur == nil) || !ul.Linked(ur)
		left = (ll == nil) || !ul.Linked(ll)
	} else {
		up = (ur != nil)   // If so, we're at the left edge of the grid
		left = (ll != nil) // If so, we're at the top edge of the grid
	}

	if lr != nil {
		down = (ll == nil) || !lr.Linked(ll)
		right = (ur == nil) || !lr.Linked(ur)
	} else {
		down = (ll != nil)
		right = (ur != nil)
	}

	return cornerGlyph(up, left, down, right)
}

// cornerGlyph returns the glyph appropriate for drawing at a corner
// given the directions it extends into
func cornerGlyph(up, left, down, right bool) rune {
	// Select the proper glyph given the directions
	idx := 0
	if up {
		idx |= 1
	}
	if down {
		idx |= 2
	}
	if left {
		idx |= 4
	}
	if right {
		idx |= 8
	}
	glyphs := [16]rune{
		' ',             // Nothing
		vertical,        // Up
		vertical,        // Down
		vertical,        // Down | Up
		horizontal,      // Left
		cornerUpLeft,    // Left | Up
		cornerDownLeft,  // Left | Down
		verticalLeft,    // Left | Down | Up
		horizontal,      // Right
		cornerUpRight,   // Right | Up
		cornerDownRight, // Right | Down
		verticalRight,   // Right | Down | Up
		horizontal,      // Right | Left
		horizontalUp,    // Right | Left | Up
		horizontalDown,  // Right | Left | Down
		intersection,    // Right | Left | Down | Up
	}
	return glyphs[idx]
}

// pointsDown returns true if the provided rune is a box drawing glyph
// which points down
func pointsDown(r rune) bool {
	switch r {
	case vertical, cornerDownRight, cornerDownLeft, verticalRight, verticalLeft, intersection:
		return true
	}
	return false
}

// pointsDown returns true if the provided rune is a box drawing glyph
// which points down
func pointsRight(r rune) bool {
	switch r {
	case horizontal, cornerDownRight, cornerUpRight, verticalRight, horizontalDown, horizontalUp, intersection:
		return true
	}
	return false
}
