package main

// Example of using MutateTree() with a function argument. Rotates the whole tree.

import (
	"fmt"
	"os"

	sgf ".."
)

func main() {
	original := sgf.LoadArgOrQuit(1)					// Equivalent to sgf.Load(os.Args[1])
	mutated := original.MutateTree(rotate_clockwise)
	mutated.Save(os.Args[1] + ".rotated.sgf")
	fmt.Printf("Saved. %d nodes in original, %d nodes in mutated.\n", original.TreeSize(), mutated.TreeSize())
}

// The mutator function is shown the original node and must return the properties
// that it wants the mutated node to have...

func rotate_clockwise(original *sgf.Node, boardsize int) map[string][]string {

	props := original.AllProperties()		// Fetches a copy, which is safe to edit.

	for _, key := range []string{"AB", "AW", "AE", "B", "CR", "MA", "SL", "SQ", "TR", "W"} {
		for i, s := range props[key] {
			if len(s) == 2 {
				x, y, onboard := sgf.ParsePoint(s, boardsize)
				if onboard {
					new_x := boardsize - 1 - y
					new_y := x
					props[key][i] = sgf.Point(new_x, new_y)
				}
			}
		}
	}

	return props
}
