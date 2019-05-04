package main

// This program takes a directory of ELF-GoGoD analysis
// (see https://lifein19x19.com/forum/viewtopic.php?f=18&t=16441)
// and adds Sabaki SBKV tags so that Sabaki can graph the winrates.

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	root, err := sgf.Load(path)
	if err != nil {
		return nil
	}

	for node := root; node != nil; node = node.MainChild() {

		comment, _ := node.GetValue("C")
		lines := strings.Split(comment, "\n")			// Always returns at least one string
		val, err := strconv.ParseFloat(strings.TrimSpace(lines[0]), 64)
		if err == nil {
			val *= 100
			node.SetValue("SBKV", fmt.Sprintf("%.2f", val))
		}
		if node.Parent() != nil {
			for _, sibling := range node.Parent().Children() {
				_, ok := sibling.GetValue("TE")
				if ok {
					b, _ := sibling.GetValue("B")
					w, _ := sibling.GetValue("W")
					if b != "" {
						node.AddValue("TR", b)
					}
					if w != "" {
						node.AddValue("TR", w)
					}
				}
			}
		}
	}

	root.Save(path)
	fmt.Printf("%s\n", path)
	return nil
}
