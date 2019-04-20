package main

// Example of using MutateTree() with a function argument. Rotates the whole tree.

import (
	"fmt"
	sgf ".."
)

func main() {
	original, _ := sgf.Load("example.sgf", true)
	original = original.GetEnd()
	mutated := original.MutateTree(rotate_clockwise)
	original.Board().DumpBoard()						// Unharmed
	fmt.Printf("\n")
	mutated.Board().DumpBoard()				// We could also save with mutated.Save()
}

// This function will be called for every node in the original tree, with the properties and
// the board position. It must return the properties it wants the mutated version to have.
// It is safe for it to modify the "props" map it is supplied with.

func rotate_clockwise(props map[string][]string, board *sgf.Board) map[string][]string {
	for _, key := range []string{"AB", "AW", "AE", "B", "CR", "MA", "SL", "SQ", "TR", "W"} {
		for i, s := range props[key] {
			if len(s) == 2 {
				x, y, onboard := sgf.ParsePoint(s, board.Size)
				if onboard {
					new_x := board.Size - 1 - y
					new_y := x
					props[key][i] = sgf.Point(new_x, new_y)
				}
			}
		}
	}
	return props
}
