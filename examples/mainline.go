package main

import (
	"fmt"
	sgf ".."
)

func main() {
	root, err := sgf.LoadSGFMainLine("kifu/2016-03-10a.sgf")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Loaded %d nodes\n", root.TreeSize())
}
