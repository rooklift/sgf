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
	node = node.GetEnd()
	node.Board().Dump()
}
