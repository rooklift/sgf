package main

// Example of mutating an entire game tree.

import (
	"fmt"
	"os"

	sgf ".."
)

func main() {

	node, err := sgf.Load(os.Args[1], true)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	mutant := mutate_recursive(node)

	mutant.Save(os.Args[1] + ".mirror.sgf")
}

func mutate_recursive(node *sgf.Node) *sgf.Node {

	mutant := make_mutant(node)

	for _, child := range(node.Children) {
		mutant_child := mutate_recursive(child)
		mutant_child.Parent = mutant
		mutant.Children = append(mutant.Children, mutant_child)
	}

	return mutant
}

func make_mutant(node *sgf.Node) *sgf.Node {

	props := node.AllProperties()

	for _, key := range []string{"B", "W", "AB", "AW", "AE"} {
		for i, s := range props[key] {
			if len(s) == 2 {
				props[key][i] = string(props[key][i][1]) + string(props[key][i][0])		// Diagonal mirror
			}
		}
	}

	mutant := sgf.NewNode(nil, props)		// nil here because we'll handle parent/child relationship manually.

	return mutant
}

