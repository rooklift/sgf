package main

// Scan a directory of SGF files for illegal moves. Recursive.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rooklift/sgf"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <dir>\n", filepath.Base(os.Args[0]))
		return
	}
	filepath.Walk(os.Args[1], handle_file)
}

func handle_file(path string, _ os.FileInfo, err error) error {

	// Returning an error halts the whole walk. So don't.

	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}

	root, err := sgf.LoadMainLine(path)
	if err != nil {
		return nil
	}

	i := 0
	node := root

	for {
		child := node.MainChild()
		if child == nil {
			break
		}

		i++

		err := node.Validate()

		if err != nil {
			re, _ := root.GetValue("RE")
			fmt.Printf("%s: Move %d of %d: %v -- %s\n", filepath.Base(path), i, len(node.GetEnd().GetLine()) - 1, err, re)
			return nil
		}

		node = child
	}

	return nil
}
