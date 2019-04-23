package main

// Example of using MutateTree() with a function argument. Flips colours.

import (
	"fmt"
	"os"
	"strings"

	sgf ".."
)

func main() {
	original := sgf.LoadArgOrQuit(1)					// Equivalent to sgf.Load(os.Args[1])
	mutated := original.MutateTree(reverse_colours)
	mutated.Save(os.Args[1] + ".inverted.sgf")
	fmt.Printf("Saved. %d nodes in original, %d nodes in mutated.\n", original.TreeSize(), mutated.TreeSize())
}

// The mutator function is shown the original node and must return the properties
// that it wants the mutated node to have...

func reverse_colours(original *sgf.Node, boardsize int) *sgf.Node {

	reverse_map := map[string]string{
		"B": "W",
		"W": "B",
		"AB": "AW",
		"AW": "AB",
		"PB": "PW",
		"PW": "PB",
	}

	node := original.Copy()

	for old_key, new_key := range reverse_map {
		node.SetValues(new_key, original.AllValues(old_key))
	}

	result, ok := node.GetValue("RE")
	if ok {
		if strings.HasPrefix(result, "B+") {
			node.SetValue("RE", "W+" + result[2:])
		} else if strings.HasPrefix(result, "W+") {
			node.SetValue("RE", "B+" + result[2:])
		}
	}

	komi, ok := node.GetValue("KM")
	if ok {
		if strings.HasPrefix(komi, "-") {
			node.SetValue("KM", komi[1:])
		} else {
			node.SetValue("KM", "-" + komi)
		}
	}

	return node
}
