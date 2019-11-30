package maze

import (
	"log"
	"math/rand"
	"time"
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
	rand.Seed(time.Now().UnixNano())
	return g.At(rand.Int63n(g.Rows), rand.Int63n(g.Columns))
}

// Size returns the number of cells in the grid
func (g *Grid) Size() int64 {
	return g.Rows * g.Columns
}
