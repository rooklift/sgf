package main

// Scan a directory of SGF files for illegal moves.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	sgf ".."
)

func main() {

	dirs := os.Args[1:]

	for _, d := range dirs {

		files, err := ioutil.ReadDir(d)

		if err != nil {
			panic(err.Error())
		}

		for _, f := range files {
			err := handle_file(d, f.Name())
			if err != nil {
				fmt.Printf("%s: %v\n", f.Name(), err)
			}
		}
	}
}

func handle_file(dirname, filename string) error {

	path := filepath.Join(dirname, filename)

	root, err := sgf.LoadMainLine(path)
	if err != nil {
		return err
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
				return fmt.Errorf("Move %d of %d: %v   %s", i, len(node.GetEnd().GetLine()) - 1, err, re)
			}
		}

		w, _ := child.GetValue("W")
		if w != "" && w != "tt" {
			_, err := board.LegalColour(w, sgf.WHITE)
			if err != nil {
				re, _ := root.GetValue("RE")
				return fmt.Errorf("Move %d of %d: %v   %s", i, len(node.GetEnd().GetLine()) - 1, err, re)
			}
		}

		node = child
	}

	return nil
}
