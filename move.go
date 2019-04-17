package sgf

import (
	"fmt"
)

func (self *Node) PlayMove(p Point) (*Node, error) {

	// Uses board info to determine colour.
	// Returns the new node on success, or self on failure.

	board := self.Board()

	if p.X < 0 || p.Y < 0 || p.X >= board.Size || p.Y >= board.Size {
		return self, fmt.Errorf("Node.PlayMove(): offboard coordinates %d,%d", p.X, p.Y)
	}

	if board.State[p.X][p.Y] != EMPTY {
		return self, fmt.Errorf("Node.PlayMove(): point %d,%d was not empty", p.X, p.Y)
	}

	if board.HasKo() && board.GetKo() == p {
		return self, fmt.Errorf("Node.PlayMove(): ko recapture forbidden")
	}

	key := "B"; if board.Player == WHITE { key = "W" }
	val := SGFFromPoint(p)		// e.g. gets "cd" or whatever

	// Return the already-extant child if there is such a thing...

	for _, child := range self.Children {
		mv, _ := child.GetValue(key)
		if mv == val {
			return child, nil
		}
	}

	proposed_node := NewNode(self, map[string][]string{key: []string{val}})		// Note: already appends child to self
	proposed_board := proposed_node.Board()

	if proposed_board.State[p.X][p.Y] == EMPTY {								// Because of suicide
		self.RemoveChild(proposed_node)											// Delete child (see above)
		return self, fmt.Errorf("Node.PlayMove(): suicide forbidden")
	}

	return proposed_node, nil
}

func (self *Node) Pass() *Node {

	// Uses board info to determine colour.

	board := self.Board()

	key := "B"; if board.Player == WHITE { key = "W" }

	// Return the already-extant child if there is such a thing...

	for _, child := range self.Children {

		val, ok := child.GetValue(key)

		if ok {		// i.e. there is a value (possibly empty string) for the key

			_, ok := PointFromSGF(val, board.Size)

			if ok == false {		// move is a pass
				return child
			}
		}
	}

	new_node := NewNode(self, map[string][]string{key: []string{""}})

	return new_node
}
