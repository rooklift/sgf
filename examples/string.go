package main

import (
	"fmt"

	"github.com/rooklift/sgf"
)

func main() {
	root := sgf.LoadArgOrQuit(1)
	fmt.Println(root.SGF())
}
