package sgf

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Printf("\n")
}

func TestIllegality(t *testing.T) {
	fmt.Printf("TestIllegality\n")

	root, err := Load("test_kifu/illegality.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	node := root.GetEnd()
	original_end := node

	node, err = node.PlayMove(Point(10,8))
	if err == nil {
		t.Errorf("Recaptured a ko")
	}

	node, err = node.PlayMove(Point(11,9))
	if err == nil {
		t.Errorf("Played a suicide move")
	}

	node, err = node.PlayMove(Point(11,10))
	if err == nil {
		t.Errorf("Played on top of a stone")
	}

	node, err = node.PlayMove(Point(19,19))
	if err == nil {
		t.Errorf("Played an off-board move")
	}

	if node != original_end {
		t.Errorf("node was not original_end")
	}

	if len(node.children) != 0 {
		t.Errorf("node gained a child somehow")
	}
}

func TestCollection(t *testing.T) {
	fmt.Printf("TestCollection\n")

	collection, err := LoadCollection("test_kifu/collection.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if len(collection) != 3 {
		t.Errorf("Collection was not of expected size")
	}

	expectations := []int{44, 244, 3793}

	for i, root := range collection {
		if root.TreeSize() != expectations[i] {
			t.Errorf("A tree was not of expected size")
		}
	}
}

func TestCyclicAttachment(t *testing.T) {
	fmt.Printf("TestCyclicAttachment\n")

	root, err := Load("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Choose a random node...

	all_nodes := root.TreeNodes()
	node := all_nodes[rand.Intn(len(all_nodes))]

	// Choose a random node in that node's subtree...

	descendents := node.SubtreeNodes()
	d := descendents[rand.Intn(len(descendents))]

	// Trying to attach the node to a descendent or itself should panic...

	defer func() {
        r := recover()
        if r == nil {
            t.Errorf("The cyclic attachment did not cause a panic")
        }
    }()

	node.SetParent(d)
}

func TestDyer(t *testing.T) {
	fmt.Printf("TestDyer\n")

	root, err := Load("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if root.Dyer() != "comhcledemrd" {
		t.Errorf("Dyer signature was not what was expected")
	}
}

func TestUnescaping(t *testing.T) {
	fmt.Printf("TestUnescaping\n")

	root, err := Load("test_kifu/escaped.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	node := root.GetEnd()

	label, _ := node.GetValue("LB")
	if label != "pd:\\" {
		t.Errorf("Label not as expected")
	}

	comment, _ := node.GetValue("C")
	if comment != "This comment has a \\ character." {
		t.Errorf("Comment not as expected")
	}
}

func TestMainLineLoader(t *testing.T) {
	fmt.Printf("TestMainLineLoader\n")

	root, err := LoadSGFMainLine("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if root.TreeSize() != 212 {
		t.Errorf("Wrong number of nodes in tree")
	}
}

func TestGibLoader(t *testing.T) {
	fmt.Printf("TestGibLoader\n")

	root, err := Load("test_kifu/3handicap.gib")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if root.TreeSize() != 253 {
		t.Errorf("Wrong number of nodes in tree")
	}

	ha, _ := root.GetValue("HA")
	if ha != "3" {
		t.Errorf("Wrong handicap")
	}

	stones := root.AllValues("AB")
	if len(stones) != 3 {
		t.Errorf("Wrong AB property")
	}
}

func TestNgfLoader(t *testing.T) {
	fmt.Printf("TestNgfLoader\n")

	root, err := Load("test_kifu/3handicap.ngf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if root.TreeSize() != 284 {
		t.Errorf("Wrong number of nodes in tree")
	}

	ha, _ := root.GetValue("HA")
	if ha != "3" {
		t.Errorf("Wrong handicap")
	}

	stones := root.AllValues("AB")
	if len(stones) != 3 {
		t.Errorf("Wrong AB property")
	}
}

func TestHandicap(t *testing.T) {
	fmt.Printf("TestHandicap\n")

	root, err := Load("test_kifu/9handicap.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	ha, _ := root.GetValue("HA")
	if ha != "9" {
		t.Errorf("Wrong handicap")
	}

	stones := root.AllValues("AB")
	if len(stones) != 9 {
		t.Errorf("Wrong AB property")
	}
}

func TestKeyValues(t *testing.T) {
	fmt.Printf("TestKeyValues\n")

	root, err := Load("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	key_count, value_count := root.TreeKeyValueCount()

	if key_count != 9562 || value_count != 9562 {
		t.Errorf("Wrong number of keys or values in tree")
	}
}

func TestUnicode(t *testing.T) {
	fmt.Printf("TestUnicode\n")

	root, err := Load("test_kifu/unicode.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	pb, _ := root.GetValue("PB")
	pw, _ := root.GetValue("PW")

	if pb != "播放機" || pw != "戰鬥機" {
		t.Errorf("Got unexpected string when reading unicode")
	}
}

func TestCache(t *testing.T) {
	fmt.Printf("TestCache\n")

	root, err := Load("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	nodes := root.SubtreeNodes()

	for _, node := range nodes {
		node.Board()
	}

	for _, node := range nodes {
		if node.__board_cache == nil {
			t.Errorf("Board cache was not made (1)")
		}
	}

	root.AddValue("AB", "aa")

	for _, node := range nodes {
		if node.__board_cache != nil {
			t.Errorf("Board cache was not purged (1)")
		}
	}

	for _, node := range nodes {
		node.Board()
	}

	for _, node := range nodes {
		if node.__board_cache == nil {
			t.Errorf("Board cache was not made (2)")
		}
	}

	root.MainChild().Detach()

	for _, node := range nodes {
		if node != root {
			if node.__board_cache != nil {
				t.Errorf("Board cache was not purged (2)")
			}
		} else {
			if node.__board_cache == nil {
				t.Errorf("Board cache of root was purged for no reason")
			}
		}
	}
}
