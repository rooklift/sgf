package sgf

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
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

	node, err = node.Play(Point(10,8))
	if err == nil {
		t.Errorf("Recaptured a ko")
	}

	node, err = node.Play(Point(11,9))
	if err == nil {
		t.Errorf("Played a suicide move")
	}

	node, err = node.Play(Point(11,10))
	if err == nil {
		t.Errorf("Played on top of a stone")
	}

	node, err = node.Play(Point(19,19))
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

	root, err := LoadMainLine("test_kifu/2016-03-10a.sgf")
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

func TestBoard(t *testing.T) {
	fmt.Printf("TestBoard\n")

	root, err := Load("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	total_board_updates = 0			// Reset global

	root.Board()
	if total_board_updates != 1 {
		t.Errorf("total_board_updates not as expected")
	}

	// Real tests...

	board := root.GetEnd().Board()
	if total_board_updates != 212 {	//
		t.Errorf("total_board_updates not as expected")
	}

	if board.CapturesBy[BLACK] != 3 || board.CapturesBy[WHITE] != 5 {
		t.Errorf("Captures not as expected")
	}

	stones := 0
	for x := 0; x < board.Size; x++ {
		for y := 0; y < board.Size; y++ {
			if board.State[x][y] != EMPTY {
				stones++
			}
		}
	}
	if stones != 203 {
		t.Errorf("Stones not as expected")
	}
}

func TestGroups(t *testing.T) {
	fmt.Printf("TestGroups\n")

	root, err := Load("test_kifu/group_info.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	board := root.Board()

	if len(board.Stones("aa")) != 57 {
		t.Errorf("len(board.Stones()) not as expected")
	}

	if len(board.Liberties("aa")) != 37 {
		t.Errorf("len(board.Liberties()) not as expected")
	}

	if board.HasLiberties("pd") {
		t.Errorf("Empty point was considered as having liberties")
	}

	if board.HasLiberties("pp") {
		t.Errorf("Empty point was considered as having liberties")
	}

	if board.DestroyGroup("aa") != 57 {
		t.Errorf("DestroyGroup did not return the expected value")
	}

	// Try adding some stones to make an illegal position...

	root, err = Load("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	board = root.GetEnd().Board()
	board.AddStone("jk", WHITE)
	board.AddStone("kk", WHITE)
	if board.HasLiberties("kk") == true || len(board.Liberties("kk")) != 0 {
		t.Errorf("Group with no liberties reported as having liberties")
	}

	// None of the group info methods should crash if given an invalid point...

	board.Stones("ZZ")
	board.HasLiberties("ZZ")
	board.Liberties("ZZ")
	board.Singleton("ZZ")
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

func TestNodeCopy(t *testing.T) {
	fmt.Printf("TestNodeCopy\n")

	root := NewNode(nil)
	node := NewNode(root)
	NewNode(node)			// Add a child.

	node.AddValue("AB", "dd")
	node.AddValue("AB", "pp")

	c := node.Copy()

	if len(c.AllKeys()) != 1 || c.KeyCount() != 1 {
		t.Errorf("Copy had wrong number of keys")
	}

	if len(c.AllValues("AB")) != 2 || c.ValueCount("AB") != 2 {
		t.Errorf("Copy had wrong number of values")
	}

	if c.Parent() != nil {
		t.Errorf("Copy had a parent")
	}

	if c.MainChild() != nil {
		t.Errorf("Copy had a child")
	}
}

func TestNodeUpdates(t *testing.T) {
	fmt.Printf("TestNodeUpdates\n")

	expect_keys := func(node *Node, n int) {
		if len(node.AllKeys()) != n || node.KeyCount() != n {
			t.Errorf("Wrong number of keys")
		}
	}

	expect_vals := func(node *Node, key string, n int) {
		if len(node.AllValues(key)) != n || node.ValueCount(key) != n {
			t.Errorf("Wrong number of values")
		}
	}

	node := NewNode(nil)
	expect_keys(node, 0)
	expect_vals(node, "AB", 0)

	node.AddValue("AB", "dd")
	expect_keys(node, 1)
	expect_vals(node, "AB", 1)

	node.AddValue("AW", "dd")
	expect_keys(node, 2)
	expect_vals(node, "AB", 1)
	expect_vals(node, "AW", 1)

	node.DeleteKey("AW")
	expect_keys(node, 1)
	expect_vals(node, "AB", 1)
	expect_vals(node, "AW", 0)

	node.AddValue("AB", "dd")			// Duplicate value, shouldn't add.
	expect_keys(node, 1)
	expect_vals(node, "AB", 1)

	node.AddValue("AB", "pp")
	expect_keys(node, 1)
	expect_vals(node, "AB", 2)

	node.AddValue("AB", "dp")
	expect_keys(node, 1)
	expect_vals(node, "AB", 3)

	node.SetValue("AB", "jj")			// SetValue should delete all others.
	expect_keys(node, 1)
	expect_vals(node, "AB", 1)

	node.DeleteValue("AB", "dd")		// Deleting a non-existant value does nothing.
	expect_keys(node, 1)
	expect_vals(node, "AB", 1)

	node.DeleteValue("AB", "AB")		// Check this doesn't delete the key.
	expect_keys(node, 1)
	expect_vals(node, "AB", 1)

	node.DeleteValue("AB", "jj")
	expect_keys(node, 0)
	expect_vals(node, "AB", 0)
}

func TestRootLoader(t *testing.T) {
	fmt.Printf("TestRootLoader\n")

	root, err := LoadRoot("test_kifu/instabranch.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if root.MainChild() != nil {
		t.Errorf("root had a child")
	}
}

func TestLine(t *testing.T) {
	fmt.Printf("TestLine\n")

	root, err := Load("test_kifu/2016-03-10a.sgf")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	end := root.GetEnd()
	line := end.GetLine()

	if len(line) != 212 {
		t.Errorf("line was not the expected length")
	}
}

func TestBoardEdits(t *testing.T) {
	fmt.Printf("TestBoardEdits\n")

	board := NewBoard(19)

	expect_next_player := func(board *Board, colour Colour) {
		if board.Player != colour {
			t.Errorf("Wrong colour to play")
		}
	}

	board.Play("pp")
	expect_next_player(board, WHITE)

	board.Play("pp")					// Fails
	expect_next_player(board, WHITE)

	board.ForceStone("pp", WHITE)		// Succeeds
	expect_next_player(board, BLACK)

	board.ForceStone("pp", WHITE)		// Succeeds
	expect_next_player(board, BLACK)

	board.ForceStone("pp", BLACK)		// Succeeds
	expect_next_player(board, WHITE)

	board.Play("dd")
	expect_next_player(board, BLACK)

	board.Pass()
	expect_next_player(board, WHITE)
}

func TestLegalMovesEquivalence(t *testing.T) {
	fmt.Printf("TestLegalMovesEquivalence\n")

	const alpha = "abcdefghijklmnopqrst"		// 20 chars, so sometimes generates offboard

	for i := 0; i < 10; i++ {

		board := NewBoard(19)
		node := NewTree(19)

		var node_err, board_err error

		for n := 0; n < 1000; n++ {
			x := rand.Intn(20)					// See above
			y := rand.Intn(20)
			p := fmt.Sprintf("%c%c", alpha[x], alpha[y])

			// Sometimes switch the colours up...

			if rand.Intn(8) == 0 {
				board_err = board.PlayColour(p, board.Player.Opposite())
				node, node_err = node.PlayColour(p, node.Board().Player.Opposite())
			} else {
				board_err = board.Play(p)
				node, node_err = node.Play(p)
			}

			if (board_err == nil && node_err != nil) || (board_err != nil && node_err == nil) {
				t.Errorf("Got differing errors")
				break
			}

			if board.Equals(node.Board()) != true {
				t.Errorf("Got differing boards")
				break
			}
		}
	}
}

func TestForcedMovesEquivalence(t *testing.T) {
	fmt.Printf("TestForcedMovesEquivalence\n")

	const alpha = "abcdefghijklmnopqrst"		// 20 chars, so sometimes generates offboard

	for i := 0; i < 10; i++ {

		board := NewBoard(19)
		node := NewTree(19)

		for n := 0; n < 1000; n++ {
			x := rand.Intn(20)					// See above
			y := rand.Intn(20)
			p := fmt.Sprintf("%c%c", alpha[x], alpha[y])

			colour := BLACK
			key := "B"
			if rand.Intn(2) == 0 {
				colour = WHITE
				key = "W"
			}

			if rand.Intn(8) == 0 {

				// Sometimes do direct board
				// manipulation with no captures.

				board.Set(p, colour)
				board.Player = colour.Opposite()
				board.ClearKo()

				key = "A" + key
				node = NewNode(node)
				node.SetValue(key, p)			// Key is AB or AW

			} else {

				// Sometimes do stone placement
				// with captures.

				board.ForceStone(p, colour)

				node = NewNode(node)
				node.SetValue(key, p)			// Key is B or W

			}

			if board.Equals(node.Board()) != true {
				t.Errorf("Got differing boards at move %d", n)
				board.Dump()
				node.Board().Dump()
				node.GetRoot().write_tree(os.Stdout)
				fmt.Printf("\n")
				break
			}
		}

		// node.GetRoot().Save("meh.sgf")
	}
}
