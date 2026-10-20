package smt

import "testing"

func TestHashDomainSeparation(t *testing.T) {
	var key Hash
	var value Hash
	key[31] = 0x01
	value[31] = 0x02
	leaf := hashLeaf(key, value)
	node := hashNode(key, value)
	if leaf == node {
		t.Fatal("leaf and node hashes must differ")
	}
}

func TestEmptyHashes(t *testing.T) {
	hashes := emptyHashes()
	if hashes[0] != keccak256([]byte{0x02}) {
		t.Fatal("invalid empty leaf hash")
	}
	for i := 1; i <= 256; i++ {
		expected := hashNode(hashes[i-1], hashes[i-1])
		if hashes[i] != expected {
			t.Fatalf("invalid empty hash at level %d", i)
		}
	}
}
