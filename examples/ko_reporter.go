package main

import (
	"fmt"
	sgf ".."
)

func main() {

	node := sgf.LoadArgOrQuit(1)		// Equivalent to sgf.Load(os.Args[1])

	for n := 1; true; n++ {
		node = node.MainChild()
		if node == nil {
			break
		}
		board := node.Board()
		if board.HasKo() {
			fmt.Printf("Move %d\n", n)
			// As a test, try making a new node by playing in the ko square. Should fail...
			_, err := node.PlayMove(board.Ko)
			fmt.Printf("%v\n", err)
		}
	}
}
