package sgf

import (
	"fmt"
	"strconv"
)

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

// MakeMainLine adjusts the tree structure so that the main line leads to this
// node.
func (self *Node) MakeMainLine() {
	node := self

	for node.parent != nil {

		for i, sibling := range node.parent.children {
			if sibling == node {
				node.parent.children[i] = node.parent.children[0]
				node.parent.children[0] = node
				break
			}
		}

		node = node.parent
	}
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

// SubTreeKeyValueCount returns the number of keys and values in a node's
// subtree, including itself.
func (self *Node) SubTreeKeyValueCount() (int, int) {
	keys := self.KeyCount()
	vals := 0
	for _, key := range self.AllKeys() {
		vals += self.ValueCount(key)
	}
	for _, child := range self.children {
		k, v := child.SubTreeKeyValueCount()
		keys += k; vals += v
	}
	return keys, vals
}

// TreeKeyValueCount returns the number of keys and values in the whole tree.
func (self *Node) TreeKeyValueCount() (int, int) {
	return self.GetRoot().SubTreeKeyValueCount()
}

// RootBoardSize travels up the tree to the root, and then finds the board size,
// which it returns as an integer. If no SZ property is present, it returns 19.
func (self *Node) RootBoardSize() int {
	root := self.GetRoot()
	sz_string, _ := root.GetValue("SZ")
	sz, _ := strconv.Atoi(sz_string)
	if sz < 1  { return 19 }
	if sz > 52 { return 52 }					// SGF limit
	return sz
}

// RootKomi travels up the tree to the root, and then finds the komi, which it
// returns as a float64. If no KM property is present, it returns 0.
func (self *Node) RootKomi() float64 {
	root := self.GetRoot()
	km_string, _ := root.GetValue("KM")
	km, _ := strconv.ParseFloat(km_string, 64)
	return km
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
