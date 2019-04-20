package main

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
	node, err := sgf.Load(os.Args[1], true)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	for {
		all_children := node.Children()
		if len(all_children) == 0 {
			break
		}
		for _, child := range all_children[1:] {
			node.RemoveChild(child)
		}
		node = node.MainChild()
	}

	node.Save(os.Args[1])
}
