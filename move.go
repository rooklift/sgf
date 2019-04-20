package sgf

import (
	"fmt"
)

func (self *Node) PlayMove(p string) (*Node, error) {							// Uses board info to determine colour.
	if self == nil { panic("Node.PlayMove(): called on nil node") }
	return self.PlayMoveColour(p, self.Board().Player)
}

func (self *Node) PlayMoveColour(p string, colour Colour) (*Node, error) {		// Returns new node on success; self on failure.

	if self == nil { panic("Node.PlayMoveColour(): called on nil node") }

	if colour != BLACK && colour != WHITE {
		panic("Node.PlayMoveColour(): no colour specified")						// This is a programming error, so panic, not error.
	}

	board := self.Board()

	x, y, onboard := ParsePoint(p, board.Size)

	if onboard == false {
		return self, fmt.Errorf("Node.PlayMoveColour(): invalid or off-board string \"%v\"", p)
	}

	if board.GetState(p) != EMPTY {
		return self, fmt.Errorf("Node.PlayMoveColour(): point \"%v\" (%v,%v) was not empty", p, x, y)
	}

	if board.HasKo() && board.Ko == p {
		if colour == board.Player {					// i.e. we've not forced a move by the wrong colour.
			return self, fmt.Errorf("Node.PlayMoveColour(): ko recapture forbidden")
		}
	}

	// Return the already-extant child if there is such a thing...

	key := "B"; if colour == WHITE { key = "W" }

	for _, child := range self.children {
		if child.ValueCount(key) == 1 {				// Ignore any illegal nodes with 2 or more...
			mv, _ := child.GetValue(key)
			if mv == p {
				return child, nil
			}
		}
	}

	proposed_node := NewNode(self)					// Note: already appends child to self.
	proposed_node.SetValue(key, p)

	if proposed_node.Board().GetState(p) == EMPTY {								// Because of suicide
		self.RemoveChild(proposed_node)											// Delete child (see above)
		return self, fmt.Errorf("Node.PlayMoveColour(): suicide forbidden")
	}

	return proposed_node, nil
}

func (self *Node) Pass() *Node {												// Uses board info to determine colour.
	if self == nil { panic("Node.Pass(): called on nil node") }
	return self.PassColour(self.Board().Player)
}

func (self *Node) PassColour(colour Colour) *Node {

	if self == nil { panic("Node.PassColour(): called on nil node") }

	if colour != BLACK && colour != WHITE {
		panic("Node.PassColour(): no colour specified")							// This is a programming error, so panic, not error.
	}

	board := self.Board()

	key := "B"; if colour == WHITE { key = "W" }

	// Return the already-extant child if there is such a thing...

	for _, child := range self.children {
		if child.ValueCount(key) == 1 {											// Ignore any illegal nodes with 2 or more...
			mv, _ := child.GetValue(key)
			if Onboard(mv, board.Size) == false {
				return child
			}
		}
	}

	new_node := NewNode(self)
	new_node.SetValue(key, "")

	return new_node
}
