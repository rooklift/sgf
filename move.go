package sgf

import (
	"fmt"
)

func (self *Node) PlayMove(p string) (*Node, error) {

	// Uses board info to determine colour.
	// Returns the new node on success, or self on failure.

	board := self.Board()

	x, y, onboard := XYFromSGF(p, board.size)

	if onboard == false {
		return self, fmt.Errorf("Node.PlayMove(): invalid or off-board string \"%v\"", p)
	}

	if board.GetState(p) != EMPTY {
		return self, fmt.Errorf("Node.PlayMove(): point \"%v\" (%v,%v) was not empty", p, x, y)
	}

	if board.HasKo() && board.GetKo() == p {
		return self, fmt.Errorf("Node.PlayMove(): ko recapture forbidden")
	}

	key := "B"; if board.player == WHITE { key = "W" }

	// Return the already-extant child if there is such a thing...

	for _, child := range self.Children {
		mv, ok := child.GetValue(key)
		if ok {
			if mv == p {
				return child, nil
			}
		}
	}

	proposed_node := NewNode(self, map[string][]string{key: []string{p}})		// Note: already appends child to self
	proposed_board := proposed_node.Board()

	if proposed_board.GetState(p) == EMPTY {									// Because of suicide
		self.RemoveChild(proposed_node)											// Delete child (see above)
		return self, fmt.Errorf("Node.PlayMove(): suicide forbidden")
	}

	return proposed_node, nil
}

func (self *Node) Pass() *Node {

	// Uses board info to determine colour.

	board := self.Board()

	key := "B"; if board.player == WHITE { key = "W" }

	// Return the already-extant child if there is such a thing...

	for _, child := range self.Children {
		mv, ok := child.GetValue(key)
		if ok {
			if Onboard(mv, board.size) == false {
				return child
			}
		}
	}

	new_node := NewNode(self, map[string][]string{key: []string{""}})

	return new_node
}
