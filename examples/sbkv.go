package main

// This program takes a directory of ELF-GoGoD analysis
// (see https://lifein19x19.com/forum/viewtopic.php?f=18&t=16441)
// and adds Sabaki SBKV tags so that Sabaki can graph the winrates.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sgf ".."
)

func main() {

	dirname := os.Args[1]

	files, err := ioutil.ReadDir(dirname)

	if err != nil {
		panic(err.Error())
	}

	for _, f := range files {
		handle_file(dirname, f.Name())
	}
}

func handle_file(dirname, filename string) error {

	path := filepath.Join(dirname, filename)

	node, err := sgf.Load(path)
	if err != nil {
		return err
	}

	for {

		comment, _ := node.GetValue("C")

		lines := strings.Split(comment, "\n")

		if len(lines) > 0 {

			val, err := strconv.ParseFloat(lines[0], 64)

			if err == nil {
				val *= 100
				node.SetValue("SBKV", fmt.Sprintf("%.2f", val))
			}
		}

		if len(node.Children) == 0 {
			node.Save(path)
			return nil
		}

		node = node.Children[0]
	}
}
