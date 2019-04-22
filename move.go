package sgf

import (
	"fmt"
)

// PlayMove attempts to play the specified move at the node. The argument should
// be an SGF-formatted coordinate, e.g. "dd". The colour is determined
// intelligently. If successful, a new node is created, and attached as a child.
// That child is then returned and the error is nil. However, if the specified
// move already existed in a child, that child is returned instead and no new
// node is created; the error is still nil. On failure, the original node is
// returned, along with an error. Failure indicates the move was illegal.
func (self *Node) PlayMove(p string) (*Node, error) {							// Uses board info to determine colour.
	return self.PlayMoveColour(p, self.Board().Player)
}

// PlayMoveColour is like PlayMove, except the colour is specified rather than
// being automatically determined.
func (self *Node) PlayMoveColour(p string, colour Colour) (*Node, error) {		// Returns new node on success; self on failure.

	if colour != BLACK && colour != WHITE {
		panic("Node.PlayMoveColour(): no colour specified")						// This is a programming error, so panic, not error.
	}

	board := self.Board()

	x, y, onboard := ParsePoint(p, board.Size)

	if onboard == false {
		return self, fmt.Errorf("Node.PlayMoveColour(): invalid or off-board string %q", p)
	}

	if board.GetState(p) != EMPTY {
		return self, fmt.Errorf("Node.PlayMoveColour(): point %q (%v,%v) was not empty", p, x, y)
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

	if proposed_node.Board().GetState(p) == EMPTY {								// Because of suicide.
		proposed_node.SetParent(nil)											// Unlink the child from self.
		return self, fmt.Errorf("Node.PlayMoveColour(): suicide forbidden")
	}

	return proposed_node, nil
}

// Pass passes. The colour is determined intelligently. Normally, a new node is
// created, and attached as a child. However, if the specified pass already
// existed in a child, that child is returned instead and no new node is
// created.
func (self *Node) Pass() *Node {												// Uses board info to determine colour.
	return self.PassColour(self.Board().Player)
}

// PassColour is like Pass, except the colour is specified rather than being
// automatically determined.
func (self *Node) PassColour(colour Colour) *Node {

	if colour != BLACK && colour != WHITE {
		panic("Node.PassColour(): no colour specified")							// This is a programming error, so panic, not error.
	}

	board := self.Board()

	key := "B"; if colour == WHITE { key = "W" }

	// Return the already-extant child if there is such a thing...

	for _, child := range self.children {
		if child.ValueCount(key) == 1 {											// Ignore any illegal nodes with 2 or more...
			mv, _ := child.GetValue(key)
			if ValidPoint(mv, board.Size) == false {
				return child
			}
		}
	}

	new_node := NewNode(self)
	new_node.SetValue(key, "")

	return new_node
}
