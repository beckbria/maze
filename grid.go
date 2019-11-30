package maze

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	debug = false
)

// Grid represents a maze grid
type Grid interface {
	RowCount() int64
	ColumnCount() int64
	At(row, column int64) *Cell
	AllRows() <-chan []*Cell
	AllCells() <-chan *Cell
	RandomCell() *Cell
	Size() int64
	ToString() string
}

// RectangleGrid represents a rectangular maze grid
type RectangleGrid struct {
	// Rows and Columns indicate the size of the grid
	rows, columns int64
	// The cells in the grid
	grid [][]*RectangleCell
}

// NewGrid creates a new rectangular grid with all cells connected to their neighbors
func NewGrid(rows, columns int64) RectangleGrid {
	if rows < 0 || columns < 0 {
		log.Fatalf("Grid dimensions invalid: [%d, %d]", rows, columns)
	}
	g := RectangleGrid{
		rows:    rows,
		columns: columns,
		grid:    make([][]*RectangleCell, rows)}
	g.prepareGrid()
	g.configureCells()
	return g
}

// RowCount returns the number of rows in the grid
func (g *RectangleGrid) RowCount() int64 {
	return g.rows
}

// ColumnCount returns the number of columns in the grid
func (g *RectangleGrid) ColumnCount() int64 {
	return g.columns
}

// At accesses a cell from the grid
func (g *RectangleGrid) At(row, column int64) *RectangleCell {
	if row < 0 || column < 0 || row >= g.rows || column >= g.columns {
		return nil
	}
	return g.grid[row][column]
}

// prepareGrid creates the cells in the grid
func (g *RectangleGrid) prepareGrid() {
	for r := int64(0); r < g.rows; r++ {
		g.grid[r] = make([]*RectangleCell, g.columns)
		for c := int64(0); c < g.columns; c++ {
			cell := NewRectangleCell(r, c)
			g.grid[r][c] = &cell
		}
	}
}

// configureCells establishes links between cells and their neighbors
func (g *RectangleGrid) configureCells() {
	for cell := range g.AllCells() {
		cell.North = g.At(cell.Row()-1, cell.Column())
		cell.South = g.At(cell.Row()+1, cell.Column())
		cell.West = g.At(cell.Row(), cell.Column()-1)
		cell.East = g.At(cell.Row(), cell.Column()+1)
	}
}

// AllRows returns a row of cells in the grid at a time
func (g *RectangleGrid) AllRows() <-chan []*RectangleCell {
	c := make(chan []*RectangleCell)
	go func() {
		for _, row := range g.grid {
			c <- row
		}
		close(c)
	}()
	return c
}

// AllCells iterates over all of the cells in the grid
func (g *RectangleGrid) AllCells() <-chan *RectangleCell {
	c := make(chan *RectangleCell)
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
func (g *RectangleGrid) RandomCell() *RectangleCell {
	return g.At(rand.Int63n(g.rows), rand.Int63n(g.columns))
}

// Size returns the number of cells in the grid
func (g *RectangleGrid) Size() int64 {
	return g.rows * g.columns
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
func (g *RectangleGrid) ToString() string {
	return g.toString(3, 1)
}

// toString creates a textual representation of the maze grid
func (g *RectangleGrid) toString(horizontalSize, verticalSize int) string {
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
	for r := int64(0); r <= g.rows; r++ {
		// Generate the representation of this row
		topEdge := "" // The horizontal lines between cells
		area := ""    // The contents of the cells
		// Loop inclusive of the column count to get the right edge
		for c := int64(0); c <= g.columns; c++ {
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
		if r < g.rows {
			for i := 0; i < verticalSize; i++ {
				output = output + area + "\n"
			}
		}
	}

	return output
}

// upperLeftCornerGlyph returns the glyph which should be shown at the
// upper-left corner of a cell
func (g *RectangleGrid) upperLeftCornerGlyph(row, column int64) rune {
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