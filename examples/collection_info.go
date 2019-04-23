package main

import (
	"fmt"
	"os"
	sgf ".."
)

func main() {
	if len(os.Args) < 2 { return }
	collection, err := sgf.LoadCollection(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	for _, root := range collection {
		fmt.Printf("Found a tree with %d nodes\n", root.TreeSize())
	}
}
