package main

import (
	"fmt"
	sgf ".."
)

func main() {
	root := sgf.LoadArgOrQuit(1)		// Equivalent to sgf.Load(os.Args[1])
	km, _ := root.GetValue("KM")
	re, _ := root.GetValue("RE")
	fmt.Printf("Komi: %q, Result: %q\n", km, re)
	fmt.Printf("Nodes in tree: %d\n", root.TreeSize())
	fmt.Printf("Final board:\n")
	root.GetEnd().Board().Dump()
}
