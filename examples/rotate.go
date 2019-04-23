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

// The mutator function is shown the original node and must return a new node with
// no parent or children.

func rotate_clockwise(original *sgf.Node, boardsize int) *sgf.Node {

	node := original.Copy()

	for _, key := range []string{"AB", "AW", "AE", "B", "CR", "MA", "SL", "SQ", "TR", "W"} {

		all_values := node.AllValues(key)

		if len(all_values) > 0 {

			node.DeleteKey(key)

			for _, val := range all_values {
				x, y, onboard := sgf.ParsePoint(val, boardsize)
				if onboard {
					new_x := boardsize - 1 - y
					new_y := x
					node.AddValue(key, sgf.Point(new_x, new_y))
				} else {
					node.AddValue(key, val)
				}
			}
		}
	}

	return node
}
