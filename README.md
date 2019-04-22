Golang library for manipulation of SGF trees (i.e. Go / Weiqi / Baduk kifu). An auto-generated [list of methods etc](https://godoc.org/github.com/fohristiwhirl/sgf) is at GoDoc.

# Technical notes

* Nodes are based on `map[string][]string`, which is what SGF nodes really are.
* Nodes also have a parent node, and zero or more child nodes.
* A tree is just a bunch of nodes connected together.
* Boards are generated only as needed, and cached.
* **NOTE TO SELF:** if a board cache becomes invalid, internally we **must** call `clear_board_cache_recursive()`.
* Nodes are generally created by playing a move at an existing node.
* Functions that want a point expect it to be an SGF-string e.g. `"dd"` is the top-left hoshi.
* Such strings can be produced with `sgf.Point(3,3)` - the numbers are zeroth based.
* Escaping of `]` and `\` characters is handled invisibly to the user.
* Behind the scenes, properties are stored in an escaped state.

# Limitations

* For weird encodings (e.g. not utf-8), some potential problems if a character contains a `]` or `\` byte.
* Assumes an SGF file has one game (the normal case) and doesn't handle collections.

# Example

```golang
package main

import (
	"fmt"
	sgf ".."
)

func main() {

	// Start a new game tree and get the root node...

	root := sgf.NewTree(19)
	node := root

	// Here we create the ancient Chinese pattern...

	node.AddValue("AB", sgf.Point(3, 3))
	node.AddValue("AB", sgf.Point(15, 15))
	node.AddValue("AW", sgf.Point(15, 3))
	node.AddValue("AW", sgf.Point(3, 15))

	// The normal way to create new nodes is by playing moves.
	// If successful, PlayMove() returns the new node.

	node, err := node.PlayMove("cf")					// "cf" is SGF-speak
	fmt.Printf("%v\n", err)								// Prints nil (no error)

	// We can get an SGF coordinate (e.g. "cf") by calling Point().
	// Note that the coordinate system is zeroth-based, from the top left.

	node, err = node.PlayMove(sgf.Point(2, 5))
	fmt.Printf("%v\n", err)								// Already filled

	// Illegal moves (including suicide and basic ko) will return an error.
	// As a convenience, PlayMove() returns the original node in this case.
	// You may still wish to check for errors...

	node, err = node.PlayMove(sgf.Point(19, 19))
	fmt.Printf("%v\n", err)								// Off-board

	// We can create variations from any node.

	node = node.Parent()
	node.PlayMove(sgf.Point(13, 2))						// Create variation 1
	node.PlayMove(sgf.Point(16, 5))						// Create variation 2

	// Colours are determined intelligently, but we can always force a colour.

	node.PlayMoveColour(sgf.Point(2, 5), sgf.WHITE)		// Create variation 3

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
	bar := node.Pass()									// Does not create a new node
	node = node.Pass()									// Does not create a new node

	fmt.Printf("%v, %v\n", foo == bar, bar == node)		// true, true

	// We can directly manipulate SGF properties...
	// We can also examine the board.

	node.SetValue("C", "White passed. Lets highlight all white stones for some reason...")

	board := node.Board()								// Note that this is a deep copy

	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			if board.State[x][y] == sgf.WHITE {
				node.AddValue("TR", sgf.Point(x, y))
			}
		}
	}

	// It is also possible to directly manage node creation and properties,
	// though this is not really recommended...

	node = sgf.NewNode(node)							// Specify the parent
	node.AddValue("B", "dj")

	// It is possible to edit board-altering properties even if a node has
	// children. All cached boards in descendent nodes will be cleared, and
	// remade as needed.

	root.AddValue("AB", "jj")							// Editing the root...
	board = node.Board()								// but looking at current node

	fmt.Printf("%v\n", board.GetState("jj") == sgf.BLACK)

	// Calling Save() will save the entire tree, regardless of node position.

	node.Save("foo.sgf")

	// We can also load files.

	node, err = sgf.Load("foo.sgf")
}
```
