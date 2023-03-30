package sgf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// SaveCollection creates a new file, and saves each tree given into that file.
// It is useful for saving the rarely-used SGF collection format. Note that the
// location of the nodes in their trees is irrelevant: in each case, the whole
// tree is always saved.
func SaveCollection(nodes []*Node, filename string) error {

	var roots []*Node

	for _, node := range nodes {
		if node != nil {
			roots = append(roots, node.GetRoot())		// Note that node.Save() relies on GetRoot() here.
		}
	}

	if len(roots) == 0 {
		return fmt.Errorf("SaveCollection(): No non-nil roots supplied")
	}

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(outfile)		// bufio for speedier output if file is huge.
	for _, root := range roots {
		root.write_tree(w)
	}
	w.Flush()							// "After all data has been written, the client should call the Flush method"

	// We didn't defer outfile.Close() like normal people so we can check its error here, just in case...

	err = outfile.Close()
	if err != nil {
		return err
	}

	return nil
}

// Save saves the entire game tree to the specified file. It does not need to be
// called from the root node, but can be called from any node in an SGF tree -
// the whole tree is always saved.
func (self *Node) Save(filename string) error {
	return SaveCollection([]*Node{self}, filename)		// Not using self.GetRoot() since SaveCollection does.
}

// SGF returns the entire tree as a string in SGF format.
func (self *Node) SGF() string {
	if self == nil {
		return "<nil>"
	}
	var buf bytes.Buffer
	self.GetRoot().write_tree(&buf)
	return buf.String()
}

func (self *Node) write_tree(w io.Writer) {

	node := self

	fmt.Fprintf(w, "(")

	for {
		node.WriteTo(w)
		if len(node.children) > 1 {
			for _, child := range node.children {
				child.write_tree(w)
			}
			break
		} else if len(node.children) == 1 {
			node = node.children[0]
			continue
		} else {
			break
		}
	}

	fmt.Fprintf(w, ")")

	// We could print a newline...
	// fmt.Fprintf(w, "\n")

	return
}

func escape_string(s string) string {

	// Note the danger of building up strings with += string(c): https://play.golang.org/p/435YV7klTuI

	var buf bytes.Buffer

	for n := 0; n < len(s); n++ {
		if s[n] == '\\' || s[n] == ']' {
			buf.WriteByte('\\')
		}
		buf.WriteByte(s[n])
	}

	return buf.String()
}

// Load reads an SGF file (or GIB or NGF file, if the filename has that
// extension) creating a tree of SGF nodes, and returning the root. The input
// file is closed automatically. If the file has more than one SGF tree (a
// rarity) only the first is loaded - use LoadCollection() for such files
// instead.
func Load(filename string) (*Node, error) {

	file_bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	data := string(file_bytes)
	return LoadData(data, filename)
}

func LoadData(data, filename string) (*Node, error) {
	// If RAM wastage was ever an issue, one can do the super-spooky:
	// data := *(*string)(unsafe.Pointer(&file_bytes))
	// See https://github.com/golang/go/issues/25484 for details.
	root, _, err := load_sgf_tree(data, nil)

	if err != nil {
		if strings.HasSuffix(strings.ToLower(filename), ".gib") {
			root, err = load_gib(data)
		} else if strings.HasSuffix(strings.ToLower(filename), ".ngf") {
			root, err = load_ngf(data)
		}
	}

	if err != nil {
		return nil, err
	}

	return root, nil
}

func LoadSGF(sgf string) (*Node, error) {
	root, _, err := load_sgf_tree(sgf, nil)
	return root, err
}

func LoadGIB(sgf string) (*Node, error) {
	return load_gib(sgf)
}

func LoadNGF(sgf string) (*Node, error) {
	return load_ngf(sgf)
}

func load_sgf_tree(sgf string, parent_of_local_root *Node) (*Node, int, error) {

	// A tree is whatever is between ( and ).
	//
	// FIXME: this is not unicode aware. Potential problems exist if
	// a unicode code point contains a meaningful character, especially
	// the bytes ] and \ although this is impossible for utf-8.

	var root *Node
	var node *Node
	var tree_started bool
	var inside_value bool
	var value bytes.Buffer						// I used to use string and += string(c), but
	var key bytes.Buffer						// ran into https://play.golang.org/p/435YV7klTuI
	var keycomplete bool

	for i := 0; i < len(sgf); i++ {

		c := sgf[i]

		if tree_started == false {
			if c <= ' ' {						// Reasonable definition of whitespace, where ' ' is byte 32.
				continue
			} else if c == '(' {
				tree_started = true
				continue
			} else {
				return nil, 0, fmt.Errorf("load_sgf_tree(): unexpected byte 0x%x before (", c)
			}
		}

		if inside_value {

			if c == '\\' {
				if len(sgf) <= i + 1 {
					return nil, 0, fmt.Errorf("load_sgf_tree(): escape character at end of input")
				}
				value.WriteByte(sgf[i + 1])
				i++								// Skip 1 character.
			} else if c == ']' {
				inside_value = false
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): value ended by ] but node was nil")
				}
				node.AddValue(key.String(), value.String())
			} else {
				value.WriteByte(c)
			}

		} else {

			if c <= ' ' || (c >= 'a' && c <= 'z') {
				continue												// Silently discard whitespace and lowercase ASCII
			} else if c == '[' {
				if node == nil {
					// The tree has ( but no ; before its first property. We could return an error.
					// Alternatively, we can tolerate this...
					node = NewNode(parent_of_local_root)
					root = node											// First node we saw in the tree.
				}
				value.Reset()
				inside_value = true
				keycomplete = true
				if key.String() == "" {
					return nil, 0, fmt.Errorf("load_sgf_tree(): value started with [ but key was \"\"")
				}
			} else if c == '(' {
				if node == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): new subtree started but node was nil")
				}
				_, chars_to_skip, err := load_sgf_tree(sgf[i:], node)	// Substrings are memory efficient in Golang.
				if err != nil {
					return nil, 0, err
				}
				i += chars_to_skip - 1		// Subtract 1: the ( character we have read is also counted by the recurse.
			} else if c == ')' {
				if root == nil {
					return nil, 0, fmt.Errorf("load_sgf_tree(): subtree ended but local root was nil")
				}
				return root, i + 1, nil									// Return characters read.
			} else if c == ';' {
				if node == nil {
					node = NewNode(parent_of_local_root)
					root = node											// First node we saw in the tree.
				} else {
					node = NewNode(node)
				}
				key.Reset()
				keycomplete = false
			} else if c >= 'A' && c <= 'Z' {
				if keycomplete {
					key.Reset()
					keycomplete = false
				}
				key.WriteByte(c)
			} else {
				return nil, 0, fmt.Errorf("load_sgf_tree(): unacceptable byte 0x%x while expecting key", c)
			}
		}
	}

	if root == nil {
		return nil, 0, fmt.Errorf("load_sgf_tree(): local root was nil at function end")
	}

	// Just being here must mean we reached the actual end of the file without
	// reading a final ')' character. Still, we can return what we have.
	// Note that load_special() relies on this.

	return root, len(sgf), nil		// Return characters read.
}

// LoadCollection loads an SGF file and returns a slice of all root nodes found
// in it. It is useful for reading the rare SGF files that are in the
// "collection" format. The input file is closed automatically. Note that it is
// OK to use this function on normal, single-tree SGF files, in which case a
// slice of length 1 will be returned.
func LoadCollection(filename string) ([]*Node, error) {

	var ret []*Node

	file_bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return ret, err
	}

	data := string(file_bytes)
	data = strings.TrimSpace(data)		// Otherwise any trailing characters will trigger an extra attempt to read a tree.

	for {
		if len(data) == 0 {
			return ret, nil
		}

		root, chars_read, err := load_sgf_tree(data, nil)

		if err != nil {
			return ret, err
		}

		ret = append(ret, root)

		data = data[chars_read:]
	}
}

// LoadMainLine loads the main line of an SGF file. Unlike Load, the whole file
// is never read into memory, making this efficient for batch statistics
// collection.
func LoadMainLine(filename string) (*Node, error) {
	return load_special(filename, false)
}

// LoadRoot loads the root node of an SGF file. Unlike Load, the whole file is
// never read into memory, making this efficient for batch statistics
// collection.
func LoadRoot(filename string) (*Node, error) {
	return load_special(filename, true)
}

func load_special(filename string, root_only bool) (*Node, error) {

	// Pull out the bare minimum bytes necessary to parse the root / mainline.
	// This relies on the parser being OK with sudden end of input.

	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	data := bytes.NewBuffer(make([]byte, 0, 256))		// Start buffer with len 0 cap 256
	reader := bufio.NewReader(infile)

	inside_value := false
	escape_flag := false
	semicolons := 0
	brackets := 0

	for {
		c, err := reader.ReadByte()
		if err != nil {
			root, _, err := load_sgf_tree(data.String(), nil)
			return root, err
		}
		data.WriteByte(c)

		if inside_value {
			if escape_flag {
				escape_flag = false
				continue
			}
			if c == '\\' {
				escape_flag = true
				continue
			}
			if c == ']' {
				inside_value = false
				continue
			}
		} else {
			if c == '[' {
				inside_value = true
				continue
			}
			if c == ')' {
				root, _, err := load_sgf_tree(data.String(), nil)
				return root, err
			}
			if c == ';' {
				semicolons++
				if root_only && semicolons >= 2 {
					data.Truncate(data.Len() - 1)		// Delete the second ; from the data.
					root, _, err := load_sgf_tree(data.String(), nil)
					return root, err
				}
				continue
			}
			if c == '(' {
				brackets++
				if root_only && brackets >= 2 {
					data.Truncate(data.Len() - 1)		// Delete the second ( from the data.
					root, _, err := load_sgf_tree(data.String(), nil)
					return root, err
				}
				continue
			}
		}
	}
}
