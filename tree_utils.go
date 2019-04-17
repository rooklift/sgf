package sgf

func (self *Node) MainChild() *Node {
	if len(self.Children) == 0 {
		return nil
	}
	return self.Children[0]
}

func (self *Node) RemoveChild(child *Node) {
	for i := len(self.Children) - 1; i >= 0; i-- {
		if self.Children[i] == child {
			self.Children = append(self.Children[:i], self.Children[i+1:]...)
		}
	}
}

func (self *Node) GetRoot() *Node {
	node := self
	for {
		if node.Parent != nil {
			node = node.Parent
		} else {
			return node
		}
	}
}

func (self *Node) GetEnd() *Node {
	node := self
	for {
		if len(node.Children) > 0 {
			node = node.Children[0]
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
		node = node.Parent
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
