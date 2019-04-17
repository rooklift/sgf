package main

// This shows the final board position in the main line.

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
	node = node.GetEnd()		// Jump to end of the main line
	node.Board().Dump()
}
