package main

// Read one or more directories of files (non-recursive).

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	sgf ".."
)

func main() {

	start_time := time.Now()

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

	fmt.Printf("Elapsed: %v\n", time.Now().Sub(start_time))
}

func handle_file(dirname, filename string) error {

	path := filepath.Join(dirname, filename)

	_, err := sgf.Load(path)
	if err != nil {
		return err
	}

	return nil
}
