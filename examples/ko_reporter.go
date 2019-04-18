package main

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

	for n := 1; true; n++ {
		node = node.MainChild()
		if node == nil {
			break
		}
		board := node.Board()
		if board.HasKo() {
			fmt.Printf("Move %d\n", n)
			// As a test, try making a new node by playing in the ko square. Should fail...
			_, err := node.PlayMove(board.GetKo())
			fmt.Printf("%v\n", err)
		}
	}
}
