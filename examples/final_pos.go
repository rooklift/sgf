package main

// This shows the final board position in the main line.

import (
	"fmt"
	"os"

	sgf ".."
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Need filename\n")
		return
	}
	root, err := sgf.Load(os.Args[1], true)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	km, _ := root.GetValue("KM")
	re, _ := root.GetValue("RE")
	fmt.Printf("Komi: \"%v\", Result: \"%v\"\n\n", km, re)
	root.GetEnd().Board().Dump()
}
