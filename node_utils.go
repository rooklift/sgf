package sgf

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

	// Nothing will stop you creating a cyclic structure, which will then hang immediately.

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

	self.clear_board_cache_recursive()		// IMPORTANT!
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
