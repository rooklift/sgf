package sgf

func (self *Node) MutateTree(mutator func(props map[string][]string) map[string][]string) *Node {
	if self == nil { panic("Node.MutateTree(): called on nil node") }
	root := self.GetRoot()
	mutant_root := mutate_recursive(root, mutator)
	return mutant_root
}

func mutate_recursive(node *Node, mutator func(props map[string][]string) map[string][]string) *Node {

	mutant := make_mutant(node, mutator)

	for _, child := range(node.Children) {
		mutant_child := mutate_recursive(child, mutator)
		mutant_child.Parent = mutant
		mutant.Children = append(mutant.Children, mutant_child)
	}

	return mutant
}

func make_mutant(node *Node, mutator func(props map[string][]string) map[string][]string) *Node {

	props := node.AllProperties()		// This returns a deep copy, so is safe to modify.

	new_props := mutator(props)

	// We call NewNode() with a nil parent so that we can handle parent/child relationships manually.
	// We could in fact pass the parent as an argument to make_mutant() and so on but it is less clean.

	mutant := NewNode(nil, new_props)

	return mutant
}
