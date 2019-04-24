package main

import (
	"os"
	sgf ".."
)

func main() {

	node := sgf.LoadArgOrQuit(1)		// Equivalent to sgf.Load(os.Args[1])

	for {
		if node.MainChild() == nil {
			break
		}
		for _, child := range node.Children()[1:] {
			child.Detach()
		}
		node = node.MainChild()
	}

	node.Save(os.Args[1])
}
