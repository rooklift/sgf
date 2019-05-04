package main

// Displays all GC (game comment) properties in the roots of the SGF files
// in the directory. GoGoD uses these a lot. They have some interesting notes.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	sgf ".."
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s dirname\n", filepath.Base(os.Args[0]))
		return
	}

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	for _, f := range files {
		err := handle_file(os.Args[1], f.Name())
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func handle_file(dirname, filename string) error {

	fullpath := filepath.Join(dirname, filename)

	root, err := sgf.LoadRoot(fullpath)
	if err != nil {
		return err
	}

	gc, _ := root.GetValue("GC")

	if gc != "" {
		gc = strings.Replace(gc, "\n", " ", -1)
		gc = strings.Replace(gc, "\r", " ", -1)
		gc = strings.Replace(gc, "  ", " ", -1)
		fmt.Println(filename, ":", gc)
	}

	return nil
}
