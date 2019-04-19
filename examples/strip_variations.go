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
	root, err := sgf.Load(os.Args[1], true)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	node := root

	for {
		for len(node.Children) > 1 {
			node.RemoveChild(node.Children[len(node.Children) - 1])
		}
		node = node.MainChild()
		if node == nil {
			break
		}
	}

	root.Save(os.Args[1])
}
