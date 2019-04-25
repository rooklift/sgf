package main

import (
	"fmt"
	sgf ".."
)

func main() {

	for sz := 4; sz <= 19; sz++ {
		for n := 2; n <= 9; n++ {
			root := sgf.NewTree(sz)
			for _, stone := range sgf.HandicapPoints(sz, n, false) {
				root.AddValue("AB", stone)
			}
			root.Board().DumpBoard()
			fmt.Printf("\n")
		}
	}
}
