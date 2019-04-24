package main

// Just a speed test on loading a directory.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	sgf ".."
)

func main() {

	st := time.Now()

	dirs := os.Args[1:]

	for _, d := range dirs {

		files, err := ioutil.ReadDir(d)		// Slow the first time you do it on a Windows 8 directory, at least...

		if err != nil {
			panic(err.Error())
		}

		for _, f := range files {
			handle_file(d, f.Name())
		}
	}

	fmt.Printf("%v\n", time.Now().Sub(st))
}

func handle_file(dirname, filename string) error {

	path := filepath.Join(dirname, filename)

	_, err := sgf.Load(path)
	if err != nil {
		return err
	}

	return nil
}
