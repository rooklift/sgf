package sgf

import (
	"fmt"
	"strconv"
)

func (self *Node) Parent() *Node {
	return self.parent
}

func (self *Node) Children() []*Node {
	var ret []*Node
	for _, child := range self.children {
		ret = append(ret, child)
	}
	return ret
}

func (self *Node) MainChild() *Node {
	if len(self.children) == 0 {
		return nil
	}
	return self.children[0]
}

func (self *Node) SetParent(new_parent *Node) {

	if self.parent != nil {
		for i := len(self.parent.children) - 1; i >= 0; i-- {
			if self.parent.children[i] == self {
				self.parent.children = append(self.parent.children[:i], self.parent.children[i+1:]...)
			}
		}
	}

	self.parent = new_parent

	if self.parent != nil {
		self.parent.children = append(self.parent.children, self)
	}

	self.cyclic_attachment_detection()
	self.clear_board_cache_recursive()		// IMPORTANT!
}

func (self *Node) cyclic_attachment_detection() {
	node := self
	for {
		if node.parent == nil {
			return
		}
		if node.parent == self {
			panic("Cyclic structure created.")
		}
		node = node.parent
	}
}

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

func (self *Node) GetEnd() *Node {

	// The end of the line we're on only. Use GetRoot().GetEnd() for the mainline end.

	node := self
	for {
		if len(node.children) > 0 {
			node = node.children[0]
		} else {
			return node
		}
	}
}

func (self *Node) GetLine() []*Node {		// The line of nodes from root to here

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

func (self *Node) SubtreeSize() int {
	count := 1
	for _, child := range self.children {
		count += child.SubtreeSize()
	}
	return count
}

func (self *Node) TreeSize() int {
	return self.GetRoot().SubtreeSize()
}

func (self *Node) SubtreeNodes() []*Node {
	ret := []*Node{self}
	for _, child := range self.children {
		ret = append(ret, child.SubtreeNodes()...)
	}
	return ret
}

func (self *Node) TreeNodes() []*Node {
	return self.GetRoot().SubtreeNodes()
}

func (self *Node) RootBoardSize() int {			// Fairly expensive, callers should save the result if needed again.
	root := self.GetRoot()
	sz_string, _ := root.GetValue("SZ")
	sz, _ := strconv.Atoi(sz_string)
	if sz < 1  { return 19 }
	if sz > 52 { return 52 }					// SGF limit
	return sz
}

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
