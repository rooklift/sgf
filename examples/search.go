package main

// Search a directory for the named player.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	sgf ".."
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s dirname playername\n", filepath.Base(os.Args[0]))
		return
	}

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	for _, f := range files {
		err := handle_file(os.Args[1], f.Name(), os.Args[2])
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func handle_file(dirname, filename, search string) error {

	fullpath := filepath.Join(dirname, filename)

	root, err := sgf.LoadRoot(fullpath)
	if err != nil {
		return err
	}

	pb, _ := root.GetValue("PB")
	pw, _ := root.GetValue("PW")

	if pb == search || pw == search {
		fmt.Printf("%v\n", fullpath)
	}

	return nil
}
