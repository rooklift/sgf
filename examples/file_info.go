package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fohristiwhirl/sgf"
)

func main() {
	if len(os.Args) < 2 { return }

	root, err := sgf.Load(os.Args[1])
	if err != nil {
		fmt.Printf("%v\n", err)
		quit()
	}

	km, _ := root.GetValue("KM")
	re, _ := root.GetValue("RE")
	pb, _ := root.GetValue("PB")
	pw, _ := root.GetValue("PW")
	dt, _ := root.GetValue("DT")
	fmt.Printf("%s (B) vs %s (W)\n", pb, pw)
	fmt.Printf("Komi: %q, Result: %q, Date: %q\n", km, re, dt)
	fmt.Printf("Nodes in tree: %d, Dyer: %s\n", root.TreeSize(), root.Dyer())
	root.GetEnd().Board().Dump()

	quit()
}

// For the sake of being useful when dragging files onto the app in Windows,
// wait until some user input is received before exiting...

func quit() {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	os.Exit(0)
}
