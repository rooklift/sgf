Golang library for manipulation of SGF trees (i.e. Go / Weiqi / Baduk kifu).

Architecture notes:

* A tree is just a bunch of nodes connected together.
* Nodes do not contain any board representation.
* Boards are generated as needed and cached.
* Therefore, properties (B, W, AB, AW, AE) cannot be altered after node creation.
* Nodes are generally created by playing a move at an existing node.
* Functions that want a point expect it to be an SGF-string e.g. "dd" is the top-left hoshi.
* Such strings can be produced with sgf.Point(3,3) - the numbers are zeroth based.
* Escaping of ] and \ characters is handled invisibly to the user.

```golang
package main

import (
	"fmt"
	
  sgf "github.com/fohristiwhirl/sgf"
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
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Illegal moves (including suicide and basic ko) will return an error.
	// As a convenience, PlayMove() returns the original node in this case.
	// You may still wish to check for errors...

	node, err = node.PlayMove(sgf.Point(2, 5))
	fmt.Printf("%v\n", err)
	node, err = node.PlayMove(sgf.Point(19, 19))
	fmt.Printf("%v\n", err)

	// We can create variations from any node.

	node = node.Parent
	node.PlayMove(sgf.Point(13, 2))						// Create variation 1
	node.PlayMove(sgf.Point(16, 5))						// Create variation 2

	// By the way, what are these mysterious sgf.Points, anyway?

	fmt.Printf("%v\n", sgf.Point(0, 0))					// Prints "aa"

	// We can iterate through a node's children.

	for i, child := range node.Children {
		child.SetValue("C", fmt.Sprintf("Comment %d", i))
	}

	// And we can go down those variations if we wish.
	// (Errors ignored here for simplicity.)

	node, _ = node.PlayMove(sgf.Point(5, 16))			// Create variation 3 and go down it
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
}
```
