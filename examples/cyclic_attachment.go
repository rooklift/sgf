package main

import (
	"math/rand"
	"time"
	sgf ".."
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {

	root := sgf.Load("kifu/2016-03-10a.sgf")

	// Choose a random node...

	all_nodes := root.TreeNodes()
	node := all_nodes[rand.Intn(len(all_nodes))]

	// Choose a random node in that node's subtree...

	descendents := node.SubtreeNodes()
	d := descendents[rand.Intn(len(descendents))]

	// Trying to attach the node to a descendent or itself should panic...

	node.SetParent(d)
}
