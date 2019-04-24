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

	for {

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

		if node.MainChild() == nil {
			node.Save(path)
			fmt.Printf("%s\n", path)
			return nil
		}

		node = node.MainChild()
	}
}
