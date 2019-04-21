package main

// This shows the final board position in the main line.

import (
	"fmt"
	sgf ".."
)

func main() {
	root := sgf.LoadArgOrQuit(1)		// Equivalent to sgf.Load(os.Args[1])
	km, _ := root.GetValue("KM")
	re, _ := root.GetValue("RE")
	fmt.Printf("Komi: %q, Result: %q\n\n", km, re)
	root.GetEnd().Board().Dump()
}
