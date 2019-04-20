package main

// This program takes a directory of ELF-GoGoD analysis
// (see https://lifein19x19.com/forum/viewtopic.php?f=18&t=16441)
// and adds Sabaki SBKV tags so that Sabaki can graph the winrates.

import (
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
			handle_file(d, f.Name())
		}
	}
}

func handle_file(dirname, filename string) error {

	path := filepath.Join(dirname, filename)

	node, err := sgf.Load(path)
	if err != nil {
		return err
	}

	node.GetEnd().Board()

	return nil
}
