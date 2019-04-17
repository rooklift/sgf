package main

// This shows the final board position in the main line.

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
	root.GetEnd().Board().Dump()
}
