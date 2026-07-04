package main

// Rotates a game tree.

import (
	"fmt"
	"os"

	"github.com/rooklift/sgf"
)

func main() {

	root := sgf.LoadArgOrQuit(1)							// Equivalent to sgf.Load(os.Args[1])
	nodes := root.TreeNodes()

	width, height := root.RootBoardSize()

	for _, node := range nodes {
		rotate(node, width, height)
	}

	if width != height {
		root.SetValue("SZ", fmt.Sprintf("%v:%v", height, width))
	}

	err := root.Save(os.Args[1] + ".rotated.sgf")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func rotate(node *sgf.Node, width, height int) {
	for _, key := range []string{"AB", "AW", "AE", "B", "CR", "MA", "SL", "SQ", "TR", "W"} {
		all_values := node.AllValues(key)
		for i, val := range all_values {
			x, y, onboard := sgf.ParsePoint(val, width, height)
			if onboard {
				all_values[i] = sgf.Point(height - 1 - y, x)
			}
		}
		node.SetValues(key, all_values)
	}
}
