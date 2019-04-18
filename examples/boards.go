package main

import (
	sgf ".."
)

func main() {
	sgf.HoshiString = "+"
	for n := 1; n < 26; n++ {
		root := sgf.NewTree(n)
		root.Board().Dump()
	}
}
