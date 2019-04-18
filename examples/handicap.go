package main

import (
	sgf ".."
)

func main() {
	for n := 2; n <= 9; n++ {
		root := sgf.NewSetup(19, sgf.HandicapPoints19(n, false), nil, sgf.WHITE)
		root.Board().Dump()
	}
}
