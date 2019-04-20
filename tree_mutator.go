package sgf

type mutFunc func(original *Node) map[string][]string

// MutateTree() takes as its argument a function which examines a node and generates
// a new property map for the mutated version of the node in the new tree.

func (self *Node) MutateTree(mutator mutFunc) *Node {

	if self == nil { panic("Node.MutateTree(): called on nil node") }

	// We mutate the entire tree but we want to return the node that's equivalent to self.
	// To accomplish this, mutate_recursive() gets a pointer to a pointer which it can set
	// when it sees that it is mutating self, which we also have to send. This is a bit more
	// complicated that one would wish.

	var final_return *Node

	mutate_recursive(self.GetRoot(), mutator, self, &final_return)

	if final_return == nil {
		panic("Node.MutateTree(): failed to set equivalent node, this is normally impossible")
	}

	return final_return
}

func mutate_recursive(node *Node, mutator mutFunc, initial_caller *Node, final_return **Node) *Node {

	// We call NewNode() with a nil parent so that we can handle parent/child relationships manually.
	// We could pass the parent as an argument to mutate_recursive() and so on, but the code is less clear.

	mutant := NewNode(nil)

	// initial_caller is the node whose equivalent we ultimately want to return at the top level.
	// See note in MutateTree().

	if node == initial_caller {
		*final_return = mutant
	}

	new_props := mutator(node)

	for key, list := range new_props {
		for _, val := range list {
			mutant.AddValue(key, val)
		}
	}

	for _, child := range(node.children) {
		mutant_child := mutate_recursive(child, mutator, initial_caller, final_return)
		mutant_child.parent = mutant
		mutant.children = append(mutant.children, mutant_child)
	}

	return mutant
}
