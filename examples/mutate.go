package main

// Example of mutating an entire game tree.
// The node.MutateTree() function is called with a function argument
// which takes the current property map and returns the new one.

import (
	"fmt"
	"os"

	sgf ".."
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Need filename\n")
		return
	}
	original, err := sgf.Load(os.Args[1], true)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	mutated := original.MutateTree(mirror_diagonal)
	mutated.Save(os.Args[1] + ".mirror.sgf")
	original.Save(os.Args[1] + ".original.sgf")		// Unharmed
}

func mirror_diagonal(props map[string][]string) map[string][]string {
	for _, key := range []string{"B", "W", "AB", "AW", "AE"} {
		for i, s := range props[key] {
			if len(s) == 2 {
				props[key][i] = string(props[key][i][1]) + string(props[key][i][0])		// Diagonal mirror
			}
		}
	}
	return props
}
