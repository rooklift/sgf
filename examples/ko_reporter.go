package main

import (
	"fmt"
	"os"

	sgf ".."
)

func main() {

	root, err := sgf.Load(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	board := root.BoardFromScratch()

	node := root

	for n := 1; true; n++ {
		node = node.MainChild()
		if node == nil {
			break
		}
		board.Update(node)
		if board.HasKo() {
			fmt.Printf("Move %d\n", n)
		}
	}
}
