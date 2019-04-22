package sgf

import (
	"fmt"
	"strconv"
)

// Parent retuns the parent of a node. This will be nil if the node is the root
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
				self.parent.children = append(self.parent.children[:i], self.parent.children[i+1:]...)
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
	for {
		if node.parent == nil {
			break
		}
		if node.parent == self {
			panic("Cyclic structure created.")
		}
		node = node.parent
	}

	// Clear the board cache (and that of all descendents) because it's invalid now.

	self.clear_board_cache_recursive()
}

// GetRoot travels up the tree, examining each node's parent until it finds the
// root node, which it returns.
func (self *Node) GetRoot() *Node {
	node := self
	for {
		if node.parent != nil {
			node = node.parent
		} else {
			return node
		}
	}
}

// GetEnd travels down the tree from the node, until it reaches a node with zero
// children. It returns that node. Note that, if GetEnd is called on a node that
// is not on the main line, the result will not be on the main line either, but
// will instead be the end of the current branch.
func (self *Node) GetEnd() *Node {

	node := self
	for {
		if len(node.children) > 0 {
			node = node.children[0]
		} else {
			return node
		}
	}
}

// GetLine returns a list of all nodes between the root and the node, inclusive.
func (self *Node) GetLine() []*Node {

	var ret []*Node

	node := self

	for {
		ret = append(ret, node)
		node = node.parent
		if node == nil {
			break
		}
	}

	// Reverse the slice...

	for left, right := 0, len(ret) - 1; left < right; left, right = left + 1, right - 1 {
		ret[left], ret[right] = ret[right], ret[left]
	}

	return ret
}

// SubtreeSize returns the number of nodes in a node's subtree, including
// itself.
func (self *Node) SubtreeSize() int {
	count := 1
	for _, child := range self.children {
		count += child.SubtreeSize()
	}
	return count
}

// TreeSize returns the number of nodes in the whole tree.
func (self *Node) TreeSize() int {
	return self.GetRoot().SubtreeSize()
}

// SubtreeNodes returns a slice of every node in a node's subtree, including
// itself.
func (self *Node) SubtreeNodes() []*Node {
	ret := []*Node{self}
	for _, child := range self.children {
		ret = append(ret, child.SubtreeNodes()...)
	}
	return ret
}

// TreeNodes returns a slice of every node in the whole tree.
func (self *Node) TreeNodes() []*Node {
	return self.GetRoot().SubtreeNodes()
}

// RootBoardSize travels up the tree to the root, and then finds the board size,
// which it returns. If no SZ property is present, it returns 19.
func (self *Node) RootBoardSize() int {
	root := self.GetRoot()
	sz_string, _ := root.GetValue("SZ")
	sz, _ := strconv.Atoi(sz_string)
	if sz < 1  { return 19 }
	if sz > 52 { return 52 }					// SGF limit
	return sz
}

// Dyer returns the Dyer Signature of the entire tree.
func (self *Node) Dyer() string {

	vals := map[int]string{20: "??", 40: "??", 60: "??", 31: "??", 51: "??", 71: "??"}

	move_count := 0

	node := self.GetRoot()
	size := node.RootBoardSize()

	for {

		for _, key := range []string{"B", "W"} {

			mv, ok := node.GetValue(key)		// Assuming just 1, as per SGF specs.

			if ok {

				move_count++

				if move_count == 20 || move_count == 40 || move_count == 60 ||
				   move_count == 31 || move_count == 51 || move_count == 71 {

					if ValidPoint(mv, size) {
						vals[move_count] = mv
					}
				}
			}
		}

		node = node.MainChild()

		if node == nil || move_count > 71 {
			break
		}
	}

	return fmt.Sprintf("%s%s%s%s%s%s", vals[20], vals[40], vals[60], vals[31], vals[51], vals[71])
}
