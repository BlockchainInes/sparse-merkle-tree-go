package smt

type Proof struct {
	Exists   bool
	Siblings [256]Hash
}

func (t *Tree) Prove(key Hash) Proof {
	var proof Proof
	_, proof.Exists = t.nodes[nodeKey{depth: 256, prefix: key}]
	for depth := uint16(256); depth > 0; depth-- {
		parentDepth := depth - 1
		siblingPrefix := prefixAt(key, depth)
		siblingPrefix[parentDepth/8] ^= 1 << (7 - parentDepth%8)
		sibling, exists := t.nodes[nodeKey{depth: depth, prefix: siblingPrefix}]
		if !exists {
			sibling = t.empty[256-int(depth)]
		}
		proof.Siblings[256-int(depth)] = sibling
	}
	return proof
}

func Verify(root, key, value Hash, proof Proof) bool {
	var current Hash
	if proof.Exists {
		current = hashLeaf(key, value)
	} else {
		if value != (Hash{}) {
			return false
		}
		current = keccak256([]byte{0x02})
	}
	for level := 0; level < 256; level++ {
		depth := uint16(255 - level)
		sibling := proof.Siblings[level]
		if bitAt(key, depth) == 0 {
			current = hashNode(current, sibling)
		} else {
			current = hashNode(sibling, current)
		}
	}
	return current == root
}
