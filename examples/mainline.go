package main

import (
	"fmt"
	sgf ".."
)

func main() {
	root, _ := sgf.LoadSGFMainLine("kifu/2016-03-10a.sgf")
	fmt.Printf("%v\n", root)
}
