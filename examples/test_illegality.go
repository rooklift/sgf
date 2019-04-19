package main

import (
	"fmt"

	sgf ".."
)

func main() {
	root, err := sgf.Load("test_illegality.sgf", true)
	node := root.GetEnd()
	original_end := node

	// Since illegal moves don't generate new nodes, it should always be that node == original_end

	node, err = node.PlayMove(sgf.Point(10,8))
	fmt.Printf("%v\n", err)
	node, err = node.PlayMove(sgf.Point(11,9))
	fmt.Printf("%v\n", err)
	node, err = node.PlayMove(sgf.Point(11,10))
	fmt.Printf("%v\n", err)
	node, err = node.PlayMove(sgf.Point(19,19))
	fmt.Printf("%v\n", err)
	node, err = node.PlayMoveColour(sgf.Point(10,8), sgf.WHITE)		// This will succeed.
	fmt.Printf("%v\n", err)

	fmt.Printf("Node is original node? %v.\n", node == original_end)
	fmt.Printf("Node has %v children.\n", len(node.Children))
	node.Board().Dump()
}
