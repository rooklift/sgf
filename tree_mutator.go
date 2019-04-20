package sgf

type mutFunc func(original *Node) map[string][]string

// MutateTree() takes as its argument a function which examines a node and generates
// a new property map for the mutated version of the node in the new tree.

func (self *Node) MutateTree(mutator mutFunc) *Node {

	if self == nil { panic("Node.MutateTree(): called on nil node") }

	// We mutate the entire tree but we want to return the node that's equivalent to self.
	// To accomplish this, mutate_recursive() gets a pointer to a pointer which it can set
	// when it sees that it is mutating self, which is the initial value of that pointer.

	foo := self

	mutate_recursive(self.GetRoot(), mutator, &foo)

	if foo == self {
		panic("Node.MutateTree(): failed to set equivalent node, this is normally impossible")
	}

	return foo
}

func mutate_recursive(node *Node, mutator mutFunc, foo **Node) *Node {

	// We call NewNode() with a nil parent so that we can handle parent/child relationships manually.
	// We could pass the parent as an argument to mutate_recursive() and so on, but the code is less clear.

	mutant := NewNode(nil)

	// foo starts off as the node whose equivalent we ultimately want to return at the top level.
	// When we actually see that node, we set foo to be the mutant. See note in MutateTree().
	// This is a slightly-too-cute way of doing it.

	if node == *foo {
		*foo = mutant
	}

	new_props := mutator(node)

	for key, list := range new_props {
		for _, val := range list {
			mutant.AddValue(key, val)
		}
	}

	for _, child := range(node.children) {
		mutant_child := mutate_recursive(child, mutator, foo)
		mutant_child.parent = mutant
		mutant.children = append(mutant.children, mutant_child)
	}

	return mutant
}
