package main

import (
	"fmt"
	sgf ".."
)

func main() {

	// To create a plain new Tree, you can generally use:
	//
	//      node := k.NewTree(size)
	//
	// But if you want handicap or other stones, one must pass
	// some actual properties and use k.NewNode().
	//
	// In this example, we create the ancient Chinese pattern.

	properties := make(map[string][]string)
	properties["AB"] = []string{sgf.SGFFromPoint(sgf.Point{3, 3}), sgf.SGFFromPoint(sgf.Point{15, 15})}		// ["dd", "pp"]
	properties["AW"] = []string{sgf.SGFFromPoint(sgf.Point{15, 3}), sgf.SGFFromPoint(sgf.Point{3, 15})}		// ["pd", "dp"]
	properties["SZ"] = []string{"19"}

	node := sgf.NewNode(nil, properties)			// nil means this node has no parent (it's the root)

	// We can now make moves.
	// If successful, PlayMove() returns the new node.

	node, err := node.PlayMove(sgf.Point{2, 5})
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Illegal moves (including suicide and basic ko) will return an error.
	// As a convenience, PlayMove() returns the original node in this case.
	// You may still wish to check for errors...

	node, err = node.PlayMove(sgf.Point{2, 5})
	if err != nil {
		fmt.Printf("%v\n", err)						// Will complain about the occupied point
	}

	// We can create variations from any node.

	node = node.Parent
	node.PlayMove(sgf.Point{13, 2})					// Create variation 1
	node.PlayMove(sgf.Point{16, 5})					// Create variation 2

	// We can iterate through a node's children.

	for i, child := range node.Children {
		child.SetValue("C", fmt.Sprintf("Comment %d", i))
	}

	// And we can go down those variations if we wish.
	// (Errors ignored here for simplicity.)

	node, _ = node.PlayMove(sgf.Point{5, 16})		// Create variation 3 and go down it
	node, _ = node.PlayMove(sgf.Point{2, 12})		// ...continue going down it
	node, _ = node.PlayMove(sgf.Point{3, 17})		// ...continue going down it

	// We can add properties, EXCEPT board-altering properties...

	val := sgf.SGFFromPoint(sgf.Point{3, 17})		// The string "pr" - corresponds to 15,17
	node.AddValue("TR", val)

	// Calling Save() will save the entire tree, regardless of node position.

	node.Save("foo.sgf")
}
