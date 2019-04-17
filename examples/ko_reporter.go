package main

import (
	"fmt"
	"os"

	sgf ".."
)

func main() {
	node, err := sgf.Load(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	for n := 1; true; n++ {
		node = node.MainChild()
		if node == nil {
			break
		}
		board := node.Board()			// This is super-slow and stupid, remaking from scratch every time. FIXME.
		if board.Ko.X != -1 {
			fmt.Printf("Move %d\n", n)
		}
	}
}
