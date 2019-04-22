package main

import (
	"fmt"
	sgf ".."
)

func main() {
	root := sgf.LoadArgOrQuit(1)		// Equivalent to sgf.Load(os.Args[1])
	fmt.Printf("Dyer: %s\n", root.Dyer())
}
