package main

import (
	"os"
	"github.com/rooklift/sgf"
)

func main() {
	root := sgf.LoadArgOrQuit(1)
	root.Save(os.Args[1] + ".resaved.sgf")
}
