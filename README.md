Golang library for manipulation of SGF trees (i.e. Go / Weiqi / Baduk kifu).

Architecture notes:

* A tree is just a bunch of nodes connected together.
* Nodes do not contain any board representation.
* Boards are generated as needed and cached.
* Therefore, properties (B, W, AB, AW, AE) cannot be altered after node creation.
* Nodes are generally created by playing a move at an existing node.
* Functions that want a point expect it to be an SGF-string e.g. "dd" is the top-left hoshi.
* Such strings can be produced with sgf.Point(3,3) - the numbers are zeroth based.
* Escaping of ] and \ characters is handled invisibly to the user.
