package maze

import (
	"math/rand"
)

// BinaryTree uses the binary tree maze creation algorithm to create a maze in a
// rectangular grid
func BinaryTree(g *Grid) {
	for cell := range(g.AllCells()) {
		neighbors := []*Cell{}
		// Each cell should be randomly linked to either its east or north neighbor
		if (cell.North != nil) {
			neighbors = append(neighbors, cell.North)
		}

		if (cell.East != nil) {
			neighbors = append(neighbors, cell.East)
		}

		if (len(neighbors) > 0) {
			cell.Link(neighbors[rand.Intn(len(neighbors))])
		}
	}
}