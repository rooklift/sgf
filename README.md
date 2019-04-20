Golang library for manipulation of SGF trees (i.e. Go / Weiqi / Baduk kifu).

# Technical notes

* SGF nodes are based on `map[string][]string`.
* Nodes also have a parent node, and zero or more child nodes.
* A tree is just a bunch of nodes connected together.
* Nodes do not contain any board representation.
* Boards are generated only as needed, and cached.
* Thus, board-altering properties (B, W, AB, AW, AE, PL) can't be set after node creation.
* Nodes are generally created by playing a move at an existing node.
* Functions that want a point expect it to be an SGF-string e.g. "dd" is the top-left hoshi.
* Such strings can be produced with sgf.Point(3,3) - the numbers are zeroth based.
* Escaping of ] and \ characters is handled invisibly to the user.
* Behind the scenes, properties are stored in an escaped state.

# Limitations

* Not unicode aware. Some potential problems if a unicode character contains a ] or \ byte.
* Assumes an SGF file has one game (the normal case) and doesn't handle collections.

# Example

TODO
