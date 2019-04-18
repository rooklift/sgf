Golang library for manipulation of SGF trees (i.e. Go / Weiqi / Baduk kifu).

Architecture notes:

* A tree is just a bunch of nodes connected together.
* Nodes do not contain any board representation.
* Boards are generated as needed and cached.
* Therefore, properties (B, W, AB, AW, AE) cannot be altered after node creation.
* Generally nodes are created by playing a move.
* Boards cannot be directly manipulated.
* Coordinates are zeroth-based, from top left.
* Escaping of ] and \ characters is handled invisibly to the user.
