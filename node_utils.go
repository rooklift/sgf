package sgf

func (self *Node) Parent() *Node {
	if self == nil { panic("Node.Parent(): called on nil node") }
	return self.parent
}

func (self *Node) Children() []*Node {
	if self == nil { panic("Node.Children(): called on nil node") }
	var ret []*Node
	for _, child := range self.children {
		ret = append(ret, child)
	}
	return ret
}

func (self *Node) MainChild() *Node {
	if self == nil { panic("Node.MainChild(): called on nil node") }
	if len(self.children) == 0 {
		return nil
	}
	return self.children[0]
}

func (self *Node) Destroy() {

	// Note: the "destroyed" node and its own descendents should not be reused.
	// They may have cached boards which are now obsolete. The name was chosen
	// to emphasise this point.

	if self == nil || self.parent == nil {
		return
	}

	for i := len(self.parent.children) - 1; i >= 0; i-- {
		if self.parent.children[i] == self {
			self.parent.children = append(self.parent.children[:i], self.parent.children[i+1:]...)
		}
	}

	self.parent = nil
}

func (self *Node) GetRoot() *Node {
	if self == nil { panic("Node.GetRoot(): called on nil node") }
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
	if self == nil { panic("Node.GetEnd(): called on nil node") }
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

	if self == nil { panic("Node.GetLine(): called on nil node") }

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

func (self *Node) GetLineIndices() []int {	// The child indices from the root to get to here

	if self == nil { panic("Node.GetLineIndices(): called on nil node") }

	var ret []int

	node := self

	for {
		if node.parent == nil {
			break
		}
		for n, sibling := range node.parent.children {
			if sibling == node {
				ret = append(ret, n)
				break
			}
		}
		node = node.parent
	}

	// Reverse the slice...

	for left, right := 0, len(ret) - 1; left < right; left, right = left + 1, right - 1 {
		ret[left], ret[right] = ret[right], ret[left]
	}

	return ret
}
