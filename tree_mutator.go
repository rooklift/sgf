package sgf

func (self *Node) MutateTree(mutator func(props map[string][]string, board *Board) map[string][]string) *Node {
	if self == nil { panic("Node.MutateTree(): called on nil node") }
	root := self.GetRoot()
	mutant_root := mutate_recursive(root, mutator)
	return mutant_root
}

func mutate_recursive(node *Node, mutator func(props map[string][]string, board *Board) map[string][]string) *Node {

	// Note that the mutator function only receives copies of stuff as its arguments, so it can do whatever.

	new_props := mutator(node.AllProperties(), node.Board())

	// We call NewNode() with a nil parent so that we can handle parent/child relationships manually.
	// We could pass the parent as an argument to mutate_recursive() and so on, but the code is less clear.

	mutant := NewNode(nil, new_props)

	for _, child := range(node.Children) {
		mutant_child := mutate_recursive(child, mutator)
		mutant_child.Parent = mutant
		mutant.Children = append(mutant.Children, mutant_child)
	}

	return mutant
}
