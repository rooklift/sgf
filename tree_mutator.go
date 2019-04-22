package sgf

type mutFunc func(original *Node, boardsize int) map[string][]string

// MutateTree creates a new tree that is isomorphic to the original tree. The only argument
// is a function which examines each node and returns a map of the keys and values which
// should exist in the isomorphic node. The node returned by MutateTree is the one in the
// new tree which is equivalent to the node on which MutateTree was called.
func (self *Node) MutateTree(mutator mutFunc) *Node {

	root := self.GetRoot()
	boardsize := root.RootBoardSize()

	// We mutate the entire tree but we want to return the node that's equivalent to self.
	// To accomplish this, mutate_recursive() gets a pointer to a pointer which it can set
	// when it sees that it is mutating self, which is the initial value of that pointer.

	foo := self

	mutate_recursive(root, boardsize, mutator, &foo)

	if foo == self {
		panic("Node.MutateTree(): failed to set equivalent node, this is normally impossible")
	}

	return foo
}

func mutate_recursive(node *Node, boardsize int, mutator mutFunc, foo **Node) *Node {

	// We call NewNode() with a nil parent so that we can handle parent/child relationships manually.
	// Although not essential, the code is clearer this way.

	mutant := NewNode(nil)

	// foo starts off as the node whose mutant we ultimately want to return at the top level.
	// When we actually see that node, we set foo to be the mutant. See note in MutateTree().
	// This is a slightly-too-cute way of doing it.

	if node == *foo {
		*foo = mutant
	}

	new_props := mutator(node, boardsize)

	for key, list := range new_props {
		for _, val := range list {
			mutant.AddValue(key, val)
		}
	}

	for _, child := range(node.children) {
		mutant_child := mutate_recursive(child, boardsize, mutator, foo)
		mutant_child.parent = mutant
		mutant.children = append(mutant.children, mutant_child)
	}

	return mutant
}
