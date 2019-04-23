package main

import (
	"fmt"
	sgf ".."
)

func main() {
	root := sgf.LoadArgOrQuit(1)		// Equivalent to sgf.Load(os.Args[1])
	km, _ := root.GetValue("KM")
	re, _ := root.GetValue("RE")
	pb, _ := root.GetValue("PB")
	pw, _ := root.GetValue("PW")
	fmt.Printf("%s (B) vs %s (W)\n", pb, pw)
	fmt.Printf("Komi: %q, Result: %q\n", km, re)
	fmt.Printf("Nodes in tree: %d\n", root.TreeSize())
	fmt.Printf("Final board:\n")
	root.GetEnd().Board().Dump()
}
