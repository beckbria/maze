package maze

import (
	"math/rand"
)

// Sidewinder uses the sidewinder maze creation algorithm to create a maze in a
// rectangular grid
func Sidewinder(g *RectangleGrid) {
	for row := range(g.AllRows()) {
		run := []*RectangleCell{}

		for _, cell := range(row) {
			run = append(run, cell)
			
			atEasternBoundary := (cell.East == nil)
			atNorthernBoundary := (cell.North == nil)

			// Chose randomly between going east or closing out this run
			shouldCloseOut := atEasternBoundary || (!atNorthernBoundary && rand.Intn(2) == 0)
			if shouldCloseOut {
				// Select a random cell from this run to connect North
				member := run[rand.Intn(len(run))]
				if (member.North != nil) {
					member.Link(member.North)
				}
			} else {
				cell.Link(cell.East)
			}
		}
	}
}