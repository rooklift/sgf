package main

// Reverses the colours in a game. Also adjusts komi, winner, and player names.

import (
	"fmt"
	"os"
	"strings"

	sgf ".."
)

var	reverse_map = map[string]string{
	"B": "W", "W": "B", "AB": "AW", "AW": "AB", "PB": "PW", "PW": "PB"}

func main() {
	root := sgf.LoadArgOrQuit(1)							// Equivalent to sgf.Load(os.Args[1])
	nodes := root.TreeNodes()

	key_count1, val_count1 := root.TreeKeyValueCount()		// Proving all properties survived.

	invert_km_re(root)
	for _, node := range nodes {
		invert(node)
	}

	key_count2, val_count2 := root.TreeKeyValueCount()		// Proving all properties survived.

	err := root.Save(os.Args[1] + ".inverted.sgf")
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("Saved. Key/value counts: %v/%v and %v/%v\n", key_count1, val_count1, key_count2, val_count2)
	}
}

func invert(node *sgf.Node) {
	dupe := node.Copy()
	for old_key, new_key := range reverse_map {
		node.SetValues(new_key, dupe.AllValues(old_key))
	}
}

func invert_km_re(root *sgf.Node) {
	result, ok := root.GetValue("RE")
	if ok {
		if strings.HasPrefix(result, "B+") {
			root.SetValue("RE", "W+" + result[2:])
		} else if strings.HasPrefix(result, "W+") {
			root.SetValue("RE", "B+" + result[2:])
		}
	}

	komi, ok := root.GetValue("KM")
	if ok {
		if strings.HasPrefix(komi, "-") {
			root.SetValue("KM", komi[1:])
		} else {
			root.SetValue("KM", "-" + komi)
		}
	}
}
