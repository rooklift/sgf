package main

import (
	"os"
	sgf ".."
)

func main() {
	root := sgf.LoadArgOrQuit(1)
	root.Save(os.Args[1] + ".resaved.sgf")
}
