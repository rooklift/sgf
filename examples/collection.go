package main

// Load every file in a directory, and save them to a new collection.
// Also serves as a test case - shouldn't generate any boards.

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	sgf ".."
)

func main() {

	dirs := os.Args[1:]

	var roots []*sgf.Node

	fmt.Printf("Finding files...\n")

	for _, d := range dirs {

		files, err := ioutil.ReadDir(d)
		if err != nil { panic(err.Error()) }

		for _, f := range files {
			root, err := sgf.Load(filepath.Join(d, f.Name()))
			if err != nil {
				fmt.Printf("  %v --- %v\n", f.Name(), err)
				continue
			}
			roots = append(roots, root)
			fmt.Printf("  %s\n", f.Name())
		}
	}

	if len(roots) == 0 {
		fmt.Printf("No roots found.\n")
		return
	}

	fmt.Printf("%d files will be included.\n", len(roots))
	fmt.Printf("%d boards were generated.\n", sgf.TotalBoardsGenerated)
	fmt.Printf("Name the file: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	t := scanner.Text()

	err := sgf.SaveCollection(roots, t)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
