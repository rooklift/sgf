package main

import (
	"fmt"
	sgf ".."
)

func main() {
	node, _ := sgf.Load("escaped.sgf", true)
	node = node.GetEnd()
	label, _ := node.GetValue("LB")
	fmt.Printf("%v\n", label)
	comment, _ := node.GetValue("C")
	fmt.Printf("%v\n", comment)
}
