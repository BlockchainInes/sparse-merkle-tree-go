package smt

import (
	"math/rand"
	"testing"
)

func TestMultipleKeyInclusion(t *testing.T) {
	tree := New()
	for i := 0; i < 32; i++ {
		var key, value Hash
		key[31] = byte(i)
		value[31] = byte(i + 1)
		tree.Update(key, value)
	}
	for i := 0; i < 32; i++ {
		var key, value Hash
		key[31] = byte(i)
		value[31] = byte(i + 1)
		proof := tree.Prove(key)
		if !proof.Exists || !Verify(tree.Root(), key, value, proof) {
			t.Fatalf("invalid inclusion proof for key %d", i)
		}
	}
}

func TestMultipleKeyNonInclusion(t *testing.T) {
	tree := New()
	for i := 0; i < 32; i++ {
		var key, value Hash
		key[31] = byte(i)
		value[31] = byte(i + 1)
		tree.Update(key, value)
	}
	for i := 32; i < 64; i++ {
		var key Hash
		key[31] = byte(i)
		proof := tree.Prove(key)
		if proof.Exists || !Verify(tree.Root(), key, Hash{}, proof) {
			t.Fatalf("invalid non-inclusion proof for key %d", i)
		}
	}
}

func TestRandomizedUpdateConsistency(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	first := New()
	second := New()
	entries := make(map[Hash]Hash)
	for i := 0; i < 100; i++ {
		var key, value Hash
		_, _ = rng.Read(key[:])
		_, _ = rng.Read(value[:])
		first.Update(key, value)
		entries[key] = value
	}
	for key, value := range entries {
		second.Update(key, value)
	}
	if first.Root() != second.Root() {
		t.Fatal("inconsistent root after randomized updates")
	}
	for key, value := range entries {
		proof := first.Prove(key)
		if !proof.Exists || !Verify(first.Root(), key, value, proof) {
			t.Fatal("invalid proof after randomized updates")
		}
	}
}

func TestTamperedProofAtEveryLevel(t *testing.T) {
	tree := New()
	var key, value Hash
	key[31] = 0x01
	value[31] = 0x42
	tree.Update(key, value)
	proof := tree.Prove(key)
	for level := 0; level < 256; level++ {
		modified := proof
		modified.Siblings[level][0] ^= 0x01
		if Verify(tree.Root(), key, value, modified) {
			t.Fatalf("tampered proof accepted at level %d", level)
		}
	}
}

func TestIncorrectRootRejected(t *testing.T) {
	tree := New()
	var key, value Hash
	key[31] = 0x01
	value[31] = 0x42
	tree.Update(key, value)
	proof := tree.Prove(key)
	incorrectRoot := tree.Root()
	incorrectRoot[0] ^= 0x01
	if Verify(incorrectRoot, key, value, proof) {
		t.Fatal("incorrect root accepted")
	}
}
