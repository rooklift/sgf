package main

import (
	"fmt"
	sgf ".."
)

func main() {

	// To create a plain new Tree, you can generally use:
	//      node := sgf.NewTree(size)
	//
	// If you want handicap or other stones, you can use:
	//		node := sgf.NewSetup(size, black_stones, white_stones, next_player)
	//
	// Here we create the ancient Chinese starting position, and specify that
	// Black plays first...

	black_stones := []string{sgf.Point(3, 3), sgf.Point(15, 15)}
	white_stones := []string{sgf.Point(15, 3), sgf.Point(3, 15)}

	node := sgf.NewSetup(19, black_stones, white_stones, sgf.BLACK)

	// We can now make moves.
	// If successful, PlayMove() returns the new node.

	node, err := node.PlayMove(sgf.Point(2, 5))
	fmt.Printf("%v\n", err)								// Prints nil

	// Illegal moves (including suicide and basic ko) will return an error.
	// As a convenience, PlayMove() returns the original node in this case.
	// You may still wish to check for errors...

	node, err = node.PlayMove(sgf.Point(2, 5))
	fmt.Printf("%v\n", err)
	node, err = node.PlayMove(sgf.Point(19, 19))
	fmt.Printf("%v\n", err)

	// We can create variations from any node.

	node = node.Parent()
	node.PlayMove(sgf.Point(13, 2))						// Create variation 1
	node.PlayMove(sgf.Point(16, 5))						// Create variation 2

	// Colours are determined intelligently, but we can always force a colour.

	node.PlayMoveColour(sgf.Point(2, 5), sgf.WHITE)		// Create variation 3

	// By the way, what are these mysterious sgf.Points, anyway?

	fmt.Printf("%v\n", sgf.Point(0, 0))					// Prints "aa"

	// We can iterate through a node's children.

	for i, child := range node.Children() {
		child.SetValue("C", fmt.Sprintf("Comment %d", i))
	}

	// And we can go down those variations if we wish.
	// (Errors ignored here for simplicity.)

	node, _ = node.PlayMove(sgf.Point(5, 16))			// Create variation 4 and go down it
	node, _ = node.PlayMove(sgf.Point(2, 12))			// ...continue going down it
	node, _ = node.PlayMove(sgf.Point(3, 17))			// ...continue going down it

	// Passes are a thing.
	// Doing the same action on the same node many times just returns the first-created child each time.

	foo := node.Pass()
	bar := node.Pass()									// Does not create a new node.
	node = node.Pass()									// Does not create a new node.

	fmt.Printf("%v, %v\n", foo == bar, bar == node)		// true, true

	// We can directly manipulate SGF properties, EXCEPT board-altering properties...

	node.AddValue("TR", sgf.Point(3, 3))				// "dd"
	node.AddValue("TR", sgf.Point(3, 15))				// "dp"
	node.AddValue("TR", sgf.Point(15, 3))				// "pd"
	node.AddValue("TR", "pp")							// We can always name the point directly

	node.DeleteValue("TR", "dp")

	// Calling Save() will save the entire tree, regardless of node position.

	node.Save("foo.sgf")

	// We can also load files.

	node, err = sgf.Load("foo.sgf")
}
