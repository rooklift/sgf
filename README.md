Golang library for manipulation of SGF trees (i.e. Go / Weiqi / Baduk kifu).

Architecture notes:

* A tree is just a bunch of nodes connected together. Nodes do not contain any board representation.
* Boards can be generated as needed and cached.
* Coordinates are zeroth-based.
