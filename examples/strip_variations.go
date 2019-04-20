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
	node, err := sgf.Load(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	for {
		if node.MainChild() == nil {
			break
		}
		for _, child := range node.Children()[1:] {
			child.Destroy()
		}
		node = node.MainChild()
	}

	node.Save(os.Args[1])
}
