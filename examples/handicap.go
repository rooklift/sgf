package main

import (
	"fmt"
	sgf ".."
)

func main() {
	for n := 2; n <= 9; n++ {
		root := sgf.NewTree(19)
		for _, stone := range sgf.HandicapPoints19(n, false) {
			root.AddValue("AB", stone)
		}
		root.Board().DumpBoard()
		fmt.Printf("\n")
	}
}
