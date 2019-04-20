package sgf

// MutateTree() takes as its argument a function which examines a node and generates
// a new property map for the mutated version of the node in the new tree.

func (self *Node) MutateTree(mutator func(original *Node) map[string][]string) *Node {

	if self == nil { panic("Node.MutateTree(): called on nil node") }

	root := self.GetRoot()
	mutant_root := mutate_recursive(root, mutator)

	// Return the node in the new tree which is equivalent to self...

	ret := mutant_root
	for _, index := range self.GetLineIndices() {
		ret = ret.children[index]
	}
	return ret
}

func mutate_recursive(node *Node, mutator func(original *Node) map[string][]string) *Node {

	new_props := mutator(node)

	// We call NewNode() with a nil parent so that we can handle parent/child relationships manually.
	// We could pass the parent as an argument to mutate_recursive() and so on, but the code is less clear.

	mutant := NewNode(nil)

	for key, list := range new_props {
		for _, val := range list {
			mutant.AddValue(key, val)
		}
	}

	for _, child := range(node.children) {
		mutant_child := mutate_recursive(child, mutator)
		mutant_child.parent = mutant
		mutant.children = append(mutant.children, mutant_child)
	}

	return mutant
}
