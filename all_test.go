package sgf

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("forced write error")
}

// Play() must reject ko recaptures, suicide, occupied points, and off-board
// points, leaving the tree unchanged.
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

// A file containing multiple game trees must load as a collection, with each
// tree complete.
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

// Attaching a node to itself or one of its own descendents would make a cycle
// in the tree, and must panic.
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

// The Dyer signature (the coordinates of moves 20, 40, 60, 31, 51, 71) of a
// known game must match the known value.
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

// Escaped ] and \ characters in values must be unescaped at load time.
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

// LoadMainLine must discard all variations, keeping only each node's main
// child.
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

// GIB format (Tygem) files must load, with handicap converted to HA and AB
// properties.
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

// NGF format (WBaduk) files must load, with handicap converted to HA and AB
// properties.
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

// An NGF file with an unreadable board size line must produce an error, not a
// panic.
func TestNgfBadBoardSize(t *testing.T) {
	fmt.Printf("TestNgfBadBoardSize\n")

	badNgf := "header\nnot-a-size\nwhite\nblack\nx\n0\nx\n6.5\n20230525\nx\nblack win\nPM00BJJ\n"

	if _, err := LoadNGF(badNgf); err == nil {
		t.Errorf("Expected bad board size to return an error")
	}
}

// A real 9 handicap game must have its stones present as AB values.
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

// TreeKeyValueCount must count every key and value in a large tree.
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

// Multi-byte UTF-8 values must survive loading intact.
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

// Boards must be generated lazily (one update per node actually needed), with
// the right stones and capture counts at the end of a real game.
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
	for x := 0; x < board.Width; x++ {
		for y := 0; y < board.Height; y++ {
			if board.State[x][y] != EMPTY {
				stones++
			}
		}
	}
	if stones != 203 {
		t.Errorf("Stones not as expected")
	}
}

// Group info methods: Stones, Liberties, HasLiberties, DestroyGroup. They
// must also tolerate illegal positions and invalid points without crashing.
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

// The board cache must fill as boards are requested, and be purged for all
// affected nodes when a board-altering property or structure change occurs.
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

// Copy must copy a node's properties but not its family relationships.
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

// Property editing basics: AddValue ignores duplicates, SetValue replaces all
// values, and deleting the last value deletes the key.
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

// LoadRoot must return the root node only, with no children attached.
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

// GetLine must return the whole path from the root to the node, inclusive.
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

// Direct board edits: Play fails on occupied points, while ForceStone always
// succeeds; both must leave the correct player to move.
func TestBoardEdits(t *testing.T) {
	fmt.Printf("TestBoardEdits\n")

	board := NewBoard(19, 19)

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

// Playing thousands of random (sometimes offboard or illegal) moves directly
// on a board, and via node.Play(), must give identical boards and errors.
func TestLegalMovesEquivalence(t *testing.T) {
	fmt.Printf("TestLegalMovesEquivalence\n")

	const alpha = "abcdefghijklmnopqrst"		// 20 chars, so sometimes generates offboard

	for i := 0; i < 10; i++ {

		board := NewBoard(19, 19)
		node := NewTree(19, 19)

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

// Same idea with forced stones and setup properties: direct board edits must
// match boards generated from the equivalent SGF nodes.
func TestForcedMovesEquivalence(t *testing.T) {
	fmt.Printf("TestForcedMovesEquivalence\n")

	const alpha = "abcdefghijklmnopqrst"		// 20 chars, so sometimes generates offboard

	for i := 0; i < 10; i++ {

		board := NewBoard(19, 19)
		node := NewTree(19, 19)

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

// A parsed tree must serialise back to the exact input string (fragile in
// principle, since key order is arbitrary in SGF).
func TestLoadSGF(t *testing.T) {
	fmt.Printf("TestLoadSGF\n")
	sgf := "(;GM[1]FF[4]CA[UTF-8]AP[Sabaki:0.52.2]KM[6.5]SZ[13]DT[2023-03-30];B[aa];W[ba];B[ca])"
	s, err := LoadSGF(sgf)
	if err != nil {
		t.Errorf("Failed to parse the SGF contents")
	}
	if s.SGF() != sgf {
		t.Errorf("Parsed and generated SGF should be the same")		// How safe is this test? Key order is arbitrary in SGF...
	}
}

// Errors from the underlying writer must be propagated when saving.
func TestWriteTreeError(t *testing.T) {
	fmt.Printf("TestWriteTreeError\n")

	if err := NewTree(19, 19).write_tree(failingWriter{}); err == nil {
		t.Errorf("Expected write_tree to return writer error")
	}
}

// A literal % in a value must not be mangled by any printf-style path.
func TestPercentageSignInComment(t *testing.T) {
	fmt.Printf("TestPercentageSignInComment\n")
	sgfData := "(;C[test%test])"
	s, _ := LoadSGF(sgfData)
	if sgfData != s.SGF() {
		t.Errorf("percentage sign not serialized correctly: %s", s.SGF())
	}
}

// Whitespace between property values is valid SGF and must be tolerated,
// e.g. line-wrapped point lists like AB[dd]\n[pp].
func TestWhitespaceBetweenValues(t *testing.T) {
	fmt.Printf("TestWhitespaceBetweenValues\n")

	for _, sgf := range []string{
		"(;GM[1]SZ[19]AB[dd]\n[pp])",
		"(;GM[1]SZ[19]AB[dd] [pp])",
		"(;GM[1]SZ[19]AB[dd]\r\n\t[pp])",
	} {
		root, err := LoadSGF(sgf)
		if err != nil {
			t.Errorf("%q did not load: %v", sgf, err)
			continue
		}
		ab := root.AllValues("AB")
		if len(ab) != 2 || ab[0] != "dd" || ab[1] != "pp" {
			t.Errorf("%q gave AB values %v", sgf, ab)
		}
	}
}

// A lowercase ident after a completed key must not silently attach its value
// to the previous key (weight[0.5] once became SZ[19][0.5]). It is an error.
// Meanwhile, FF[3] style idents like CoPyright must still work.
func TestLowercaseKeys(t *testing.T) {
	fmt.Printf("TestLowercaseKeys\n")

	if _, err := LoadSGF("(;GM[1]SZ[19]weight[0.5])"); err == nil {
		t.Errorf("Lowercase ident did not cause an error")
	}

	root, err := LoadSGF("(;GM[1]SZ[19]CoPyright[meh])")
	if err != nil {
		t.Errorf("FF[3] ident failed to load: %v", err)
	} else if cp, _ := root.GetValue("CP"); cp != "meh" {
		t.Errorf("FF[3] ident CoPyright did not become CP")
	}
}

// A UTF-8 byte order mark must be tolerated by every loading entry point.
func TestBOM(t *testing.T) {
	fmt.Printf("TestBOM\n")

	bom_sgf := "\xef\xbb\xbf(;GM[1]FF[4]SZ[19];B[dd];W[pp])"

	filename := filepath.Join(t.TempDir(), "bom.sgf")
	if err := os.WriteFile(filename, []byte(bom_sgf), 0644); err != nil {
		t.Fatalf(err.Error())
	}

	if _, err := LoadSGF(bom_sgf); err != nil {
		t.Errorf("LoadSGF: %v", err)
	}
	if _, err := Load(filename); err != nil {
		t.Errorf("Load: %v", err)
	}
	if _, err := LoadRoot(filename); err != nil {
		t.Errorf("LoadRoot: %v", err)
	}
	if _, err := LoadMainLine(filename); err != nil {
		t.Errorf("LoadMainLine: %v", err)
	}

	// A collection: the parser's returned character counts must stay aligned
	// with the input string despite the skipped BOM bytes...

	roots, err := LoadCollectionSGF("\xef\xbb\xbf(;GM[1]SZ[19];B[dd])(;GM[1]SZ[9];B[cc])")
	if err != nil {
		t.Errorf("LoadCollectionSGF: %v", err)
	} else if len(roots) != 2 {
		t.Errorf("LoadCollectionSGF: got %d trees, wanted 2", len(roots))
	} else {
		sz, _ := roots[1].GetValue("SZ")
		if sz != "9" {
			t.Errorf("LoadCollectionSGF: second tree had SZ %q", sz)
		}
	}

	// A BOM alone is not a tree...

	if _, err := LoadSGF("\xef\xbb\xbf"); err == nil {
		t.Errorf("BOM-only string did not cause an error")
	}
}

// GTP coordinates skip the letter I.
func TestParseGTP(t *testing.T) {
	fmt.Printf("TestParseGTP\n")

	tests := map[string]string{
		"I5":  "",			// No I in GTP
		"i5":  "",
		"H5":  "ho",
		"J5":  "io",		// J is adjacent to H
		"A1":  "as",
		"T19": "sa",
		"A0":  "",			// Off board
		"U1":  "",			// Off board
	}

	for s, expected := range tests {
		if result := ParseGTP(s, 19, 19); result != expected {
			t.Errorf("ParseGTP(%q, 19, 19) returned %q, expected %q", s, result, expected)
		}
	}
}

// MakeMainLine must preserve the relative order of the displaced siblings.
func TestMakeMainLineOrder(t *testing.T) {
	fmt.Printf("TestMakeMainLineOrder\n")

	root := NewTree(19, 19)
	for _, p := range []string{"aa", "bb", "cc", "dd"} {
		root.Play(p)
	}

	get_order := func() string {
		s := ""
		for _, child := range root.Children() {
			mv, _ := child.GetValue("B")
			s += mv
		}
		return s
	}

	root.Children()[2].MakeMainLine()
	if get_order() != "ccaabbdd" {
		t.Errorf("Expected order ccaabbdd, got %s", get_order())
	}

	root.Children()[0].MakeMainLine()		// Already main line: no change.
	if get_order() != "ccaabbdd" {
		t.Errorf("Expected order ccaabbdd, got %s", get_order())
	}

	root.Children()[3].MakeMainLine()
	if get_order() != "ddccaabb" {
		t.Errorf("Expected order ddccaabb, got %s", get_order())
	}
}

// Writing a key that the parser could not read back must panic.
func TestBadKeyPanics(t *testing.T) {
	fmt.Printf("TestBadKeyPanics\n")

	expect_panic := func(desc string, fn func()) {
		defer func() {
			if recover() == nil {
				t.Errorf("%s did not panic", desc)
			}
		}()
		fn()
	}

	node := NewNode(nil)

	expect_panic("Lowercase key", func() { node.SetValue("weight", "0.5") })
	expect_panic("Empty key", func() { node.SetValue("", "x") })
	expect_panic("Mixed case key", func() { node.AddValue("aB", "dd") })
	expect_panic("Key with bracket", func() { node.AddValue("B[", "dd") })
}

// Setup properties (AB / AW / AE) must not change the player to move.
// Exception by convention: a root that sets up Black stones only (i.e. a
// handicap game) with no PL property means White is next to move.
func TestSetupPlayer(t *testing.T) {
	fmt.Printf("TestSetupPlayer\n")

	tests := []struct {
		sgf      string
		expected Colour
	}{
		{"(;GM[1]SZ[19]HA[2]AB[pd][dp])", WHITE},			// Handicap convention
		{"(;GM[1]SZ[19]AB[pd][dp]PL[B])", BLACK},			// PL beats the convention
		{"(;GM[1]SZ[19]AB[pd]AW[dp])", BLACK},				// Mixed setup: default stands
		{"(;GM[1]SZ[19];AB[dd])", BLACK},					// Mid-tree AB: player unchanged
		{"(;GM[1]SZ[19];B[pd];AW[dd])", WHITE},				// Mid-tree AW: player unchanged
		{"(;GM[1]SZ[19];B[pd];AE[pd])", WHITE},				// Mid-tree AE: player unchanged
	}

	for _, test := range tests {
		root, err := LoadSGF(test.sgf)
		if err != nil {
			t.Errorf("%q did not load: %v", test.sgf, err)
			continue
		}
		if player := root.GetEnd().Board().Player; player != test.expected {
			t.Errorf("%q: next player was %v, expected %v", test.sgf, player.Word(), test.expected.Word())
		}
	}

	// Board-level: AddStone and AddList must leave the player alone...

	board := NewBoard(19, 19)
	board.AddStone("dd", WHITE)
	board.AddList("aa:bb", WHITE)
	if board.Player != BLACK {
		t.Errorf("AddStone / AddList changed the player")
	}

	// End to end: at the root of a real handicap game, Play() must choose White...

	root, err := Load("test_kifu/3handicap.gib")
	if err != nil {
		t.Fatalf(err.Error())
	}
	node, err := root.Play("jj")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if _, ok := node.GetValue("W"); !ok {
		t.Errorf("Play() at handicap root did not choose White")
	}
}

// -------------------------------------------------------------------------------------------------
// Rectangular boards, i.e. SZ[width:height]...

// RootBoardSize must handle both SZ formats, defaulting to 19x19 whenever the
// value is missing, malformed, or out of range.
func TestRectangularSZ(t *testing.T) {
	fmt.Printf("TestRectangularSZ\n")

	tests := []struct {
		sgf				string
		width, height	int
	}{
		{"(;GM[1]SZ[19])", 19, 19},
		{"(;GM[1]SZ[19:9])", 19, 9},
		{"(;GM[1]SZ[9:19])", 9, 19},
		{"(;GM[1]SZ[52:52])", 52, 52},			// Spec says square must not use this format, but tolerate it
		{"(;GM[1])", 19, 19},					// No SZ at all
		{"(;GM[1]SZ[foo])", 19, 19},
		{"(;GM[1]SZ[19:])", 19, 19},
		{"(;GM[1]SZ[:9])", 19, 19},
		{"(;GM[1]SZ[0:5])", 19, 19},
		{"(;GM[1]SZ[53:19])", 19, 19},
	}

	for _, test := range tests {
		root, err := LoadSGF(test.sgf)
		if err != nil {
			t.Errorf("%q did not load: %v", test.sgf, err)
			continue
		}
		width, height := root.RootBoardSize()
		if width != test.width || height != test.height {
			t.Errorf("%q gave size %dx%d, expected %dx%d", test.sgf, width, height, test.width, test.height)
		}
	}
}

// A board from a rectangular game must have the right dimensions, and the x
// and y bounds must not be interchangeable.
func TestRectangularBoard(t *testing.T) {
	fmt.Printf("TestRectangularBoard\n")

	// 19 wide, 9 tall, with stones in two corners. The final B[dj] is
	// off-board (though it would be fine on 19x19) and is thus a pass...

	root, err := LoadSGF("(;GM[1]FF[4]SZ[19:9];B[sa];W[si];B[dj])")
	if err != nil {
		t.Fatalf(err.Error())
	}

	board := root.GetEnd().Board()

	if board.Width != 19 || board.Height != 9 {
		t.Errorf("Board was %dx%d, expected 19x9", board.Width, board.Height)
	}

	if board.Get("sa") != BLACK || board.Get("si") != WHITE {
		t.Errorf("Corner stones not as expected")
	}

	if board.Player != WHITE {
		t.Errorf("Off-board move did not act as a pass")
	}

	if ValidPoint("dj", 19, 9) || ValidPoint("dj", 19, 19) == false {
		t.Errorf("ValidPoint bounds not as expected")
	}

	// String() must produce Height rows of Width points. This is a regression
	// test: the loop bounds were once transposed, panicking on any board with
	// Width != Height...

	lines := strings.Split(strings.TrimRight(board.String(), "\n"), "\n")
	if len(lines) != 9 {
		t.Errorf("Board printout had %d rows, expected 9", len(lines))
	}
	for _, line := range lines {
		if len(line) != 19 * 2 {				// Each point prints as 2 characters
			t.Errorf("Board printout row had length %d, expected 38", len(line))
		}
	}

	// Adjacency at the edges of a 19x9 board...

	if len(AdjacentPoints("si", 19, 9)) != 2 {	// Bottom right corner
		t.Errorf("Corner had wrong number of adjacent points")
	}
	if len(AdjacentPoints("se", 19, 9)) != 3 {	// Middle of right edge
		t.Errorf("Edge point had wrong number of adjacent points")
	}

	// Point lists (e.g. AB[hd:ie] style rectangles) are bounded by both
	// dimensions...

	if len(ParsePointList("hd:ie", 9, 5)) != 4 {
		t.Errorf("ParsePointList not as expected")
	}
	if ParsePointList("hd:if", 9, 5) != nil {	// "if" is off a 9x5 board
		t.Errorf("Point list going off-board should be nil")
	}

	// Boards of transposed dimensions are not equal...

	if NewBoard(9, 5).Equals(NewBoard(5, 9)) {
		t.Errorf("Boards of transposed dimensions compared equal")
	}
}

// Captures and suicide detection must work at the edges of a rectangular
// board; in particular there must be no phantom liberties off the short sides.
func TestRectangularCaptures(t *testing.T) {
	fmt.Printf("TestRectangularCaptures\n")

	// On a 9x5 board, capture a stone in the bottom right corner...

	root := NewTree(9, 5)
	root.AddValue("AB", Point(8, 4))
	root.AddValue("AW", Point(7, 4))

	node, err := root.PlayColour(Point(8, 3), WHITE)
	if err != nil {
		t.Fatalf(err.Error())
	}

	board := node.Board()

	if board.Get(Point(8, 4)) != EMPTY {
		t.Errorf("Corner stone was not captured")
	}
	if board.CapturesBy[WHITE] != 1 {
		t.Errorf("Wrong capture count")
	}

	// Playing back into that corner would now be suicide...

	if _, err := node.PlayColour(Point(8, 4), BLACK); err == nil {
		t.Errorf("Corner suicide was allowed")
	}
}

// GTP coordinates on a rectangular board: the row number is counted from the
// bottom, so the conversion depends on the height.
func TestRectangularGTP(t *testing.T) {
	fmt.Printf("TestRectangularGTP\n")

	tests := map[string]string{		// For a board 9 wide and 5 tall
		"A5": "aa",					// Top left
		"A1": "ae",					// Bottom left
		"J5": "ia",					// Top right (no I in GTP)
		"J1": "ie",					// Bottom right
		"A6": "",					// Too high
		"K1": "",					// Too far right
	}

	for s, expected := range tests {
		if result := ParseGTP(s, 9, 5); result != expected {
			t.Errorf("ParseGTP(%q, 9, 5) returned %q, expected %q", s, result, expected)
		}
	}
}

// NewTree must emit SZ[w:h] only when the board is not square, and a
// rectangular game must survive a save / load cycle intact.
func TestRectangularRoundTrip(t *testing.T) {
	fmt.Printf("TestRectangularRoundTrip\n")

	if sz, _ := NewTree(13, 13).GetValue("SZ"); sz != "13" {
		t.Errorf("Square SZ was %q, expected \"13\"", sz)
	}

	root := NewTree(19, 9)
	if sz, _ := root.GetValue("SZ"); sz != "19:9" {
		t.Errorf("Rectangular SZ was %q, expected \"19:9\"", sz)
	}

	node, err := root.Play(Point(18, 8))
	if err != nil {
		t.Fatalf(err.Error())
	}
	node, err = node.Play(Point(0, 0))
	if err != nil {
		t.Fatalf(err.Error())
	}

	reload, err := LoadSGF(root.SGF())
	if err != nil {
		t.Fatalf(err.Error())
	}

	width, height := reload.RootBoardSize()
	if width != 19 || height != 9 {
		t.Errorf("Reloaded size was %dx%d, expected 19x9", width, height)
	}

	if reload.GetEnd().Board().Equals(node.Board()) == false {
		t.Errorf("Reloaded board did not equal the original")
	}
}
