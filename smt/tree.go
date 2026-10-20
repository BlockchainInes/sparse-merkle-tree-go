package smt

type Tree struct {
	nodes  map[nodeKey]Hash
	empty  [257]Hash
	root   Hash
}

type nodeKey struct {
	depth  uint16
	prefix Hash
}

func New() *Tree {
	empty := emptyHashes()
	return &Tree{
		nodes: make(map[nodeKey]Hash),
		empty: empty,
		root:  empty[256],
	}
}

func (t *Tree) Root() Hash {
	return t.root
}

func bitAt(key Hash, depth uint16) byte {
	return (key[depth/8] >> (7 - depth%8)) & 1
}

func prefixAt(key Hash, depth uint16) Hash {
	var prefix Hash
	copy(prefix[:], key[:])
	if depth == 256 {
		return prefix
	}
	byteIndex := depth / 8
	bitOffset := depth % 8
	if bitOffset != 0 {
		prefix[byteIndex] &= 0xFF << (8 - bitOffset)
		byteIndex++
	}
	for i := byteIndex; i < 32; i++ {
		prefix[i] = 0
	}
	return prefix
}

func (t *Tree) Update(key, value Hash) {
	leaf := hashLeaf(key, value)
	t.nodes[nodeKey{depth: 256, prefix: key}] = leaf
	current := leaf
	for depth := uint16(256); depth > 0; depth-- {
		parentDepth := depth - 1
		siblingPrefix := prefixAt(key, depth)
		byteIndex := parentDepth / 8
		bitOffset := parentDepth % 8
		siblingPrefix[byteIndex] ^= 1 << (7 - bitOffset)
		siblingKey := nodeKey{depth: depth, prefix: siblingPrefix}
		sibling, exists := t.nodes[siblingKey]
		if !exists {
			sibling = t.empty[256-int(depth)]
		}
		if bitAt(key, parentDepth) == 0 {
			current = hashNode(current, sibling)
		} else {
			current = hashNode(sibling, current)
		}
		parentKey := nodeKey{depth: parentDepth, prefix: prefixAt(key, parentDepth)}
		t.nodes[parentKey] = current
	}
	t.root = current
}
