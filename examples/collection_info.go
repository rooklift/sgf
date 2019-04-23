package main

import (
	"fmt"
	"os"
	sgf ".."
)

func main() {
	if len(os.Args) < 2 { return }
	collection, _ := sgf.LoadCollection(os.Args[1])
	for _, root := range collection {
		fmt.Printf("Found a tree with %d nodes\n", root.TreeSize())
	}
}
