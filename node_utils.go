package sgf

import (
	"fmt"
	"sort"
)

// Parent returns the parent of a node. This will be nil if the node is the root
// of the tree.
func (self *Node) Parent() *Node {
	return self.parent
}

// Children returns a new slice of pointers to all the node's children.
func (self *Node) Children() []*Node {
	var ret []*Node
	for _, child := range self.children {
		ret = append(ret, child)
	}
	return ret
}

// MainChild returns the first child a node has. If the node has zero children,
// nil is returned.
func (self *Node) MainChild() *Node {
	if len(self.children) == 0 {
		return nil
	}
	return self.children[0]
}

// SetParent sets a node's parent. The node is also removed from the original
// parent's list of children, and added to the new parent's list. SetParent
// panics if a cyclic tree is created.
func (self *Node) SetParent(new_parent *Node) {

	// Delete from parent's list of children...

	if self.parent != nil {
		for i := len(self.parent.children) - 1; i >= 0; i-- {
			if self.parent.children[i] == self {
				self.parent.children = append(self.parent.children[:i], self.parent.children[i + 1:]...)
			}
		}
	}

	// Attach to parent (at both ends)...

	self.parent = new_parent

	if self.parent != nil {
		self.parent.children = append(self.parent.children, self)
	}

	// Check no cyclic structure was created...

	node := self

	for node.parent != nil {
		if node.parent == self {
			panic("Cyclic structure created.")
		}
		node = node.parent
	}

	// Clear the board cache (and that of all descendents) because it's invalid now.

	self.clear_board_cache_recursive()
}

// Detach sets a node's parent to nil, turning the node into a root. The node is
// also removed from the parent's list of children.
func (self *Node) Detach() {
	self.SetParent(nil)
}

// String returns some human-readable info about the node.
func (self *Node) String() string {
	if self == nil {
		return "<nil>"
	}

	noun := "children" ; if len(self.children) == 1 { noun = "child" }

	keys := self.AllKeys()
	sort.Strings(keys)

	return fmt.Sprintf("Node %p: depth %d, %d %s, subtree size %d, keys %v",
				self, len(self.GetLine()) - 1, len(self.children), noun, self.SubtreeSize(), keys)
}

// Validate checks a node for obvious problems; it returns the first problem
// found as an error, otherwise it returns nil.
func (self *Node) Validate() error {

	all_b := self.AllValues("B")
	all_w := self.AllValues("W")
	all_ab := self.AllValues("AB")
	all_aw := self.AllValues("AW")
	all_ae := self.AllValues("AE")

	if len(all_b) + len(all_w) > 1 {
		return fmt.Errorf("Too many moves (B or W values)")
	}

	if len(all_b) + len(all_w) > 0 && len(all_ab) + len(all_aw) + len(all_ae) > 0 {
		return fmt.Errorf("Mix of move and setup properties")
	}

	if self.parent != nil {

		board := self.parent.Board()

		if len(all_b) > 0 {
			mv := all_b[0]
			if ValidPoint(mv, board.Size) {
				legal, err := board.LegalColour(mv, BLACK)
				if legal == false {
					return err
				}
			} else if mv != "" && mv != "tt" {
				return fmt.Errorf("Got B value: %s", mv)
			}
		}

		if len(all_w) > 0 {
			mv := all_w[0]
			if ValidPoint(mv, board.Size) {
				legal, err := board.LegalColour(mv, WHITE)
				if legal == false {
					return err
				}
			} else if mv != "" && mv != "tt" {
				return fmt.Errorf("Got W value: %s", mv)
			}
		}
	}

	return nil
}
