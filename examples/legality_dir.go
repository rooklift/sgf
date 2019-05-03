package main

// Scan a directory of SGF files for illegal moves. Recursive.

import (
	"fmt"
	"os"
	"path/filepath"

	sgf ".."
)

func main() {
	if len(os.Args) < 2 { return }
	filepath.Walk(os.Args[1], handle_file)
}

func handle_file(path string, _ os.FileInfo, err error) error {

	// Returning an error halts the whole walk. So don't.

	if err != nil {
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

		board := node.Board()

		b, _ := child.GetValue("B")
		if b != "" && b != "tt" {
			_, err := board.LegalColour(b, sgf.BLACK)
			if err != nil {
				re, _ := root.GetValue("RE")
				fmt.Printf("%s: Move %d of %d: %v -- %s\n", filepath.Base(path), i, len(node.GetEnd().GetLine()) - 1, err, re)
				return nil
			}
		}

		w, _ := child.GetValue("W")
		if w != "" && w != "tt" {
			_, err := board.LegalColour(w, sgf.WHITE)
			if err != nil {
				re, _ := root.GetValue("RE")
				fmt.Printf("%s: Move %d of %d: %v -- %s\n", filepath.Base(path), i, len(node.GetEnd().GetLine()) - 1, err, re)
				return nil
			}
		}

		node = child
	}

	return nil
}
