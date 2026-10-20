package smt

import (
	"crypto/sha256"
	"testing"
)

func fuzzHash(data []byte) Hash {
	return Hash(sha256.Sum256(data))
}

func FuzzInclusionProof(f *testing.F) {
	f.Add([]byte{0x01}, []byte{0x42})
	f.Add([]byte{0xFF, 0x00}, []byte{0x42, 0x43})
	f.Fuzz(func(t *testing.T, keyBytes, valueBytes []byte) {
		key := fuzzHash(keyBytes)
		value := fuzzHash(valueBytes)
		tree := New()
		tree.Update(key, value)
		proof := tree.Prove(key)
		if !proof.Exists || !Verify(tree.Root(), key, value, proof) {
			t.Fatal("valid inclusion proof rejected")
		}
	})
}

func FuzzTamperedProof(f *testing.F) {
	f.Add([]byte{0x01}, []byte{0x42})
	f.Fuzz(func(t *testing.T, keyBytes, valueBytes []byte) {
		key := fuzzHash(keyBytes)
		value := fuzzHash(valueBytes)
		tree := New()
		tree.Update(key, value)
		proof := tree.Prove(key)
		if !proof.Exists || !Verify(tree.Root(), key, value, proof) {
			t.Fatal("valid inclusion proof rejected")
		}
		for level := 0; level < 256; level++ {
			modified := proof
			modified.Siblings[level][0] ^= 0x01
			if Verify(tree.Root(), key, value, modified) {
				t.Fatalf("tampered proof accepted at level %d", level)
			}
		}
	})
}

func FuzzNonInclusionProof(f *testing.F) {
	f.Add([]byte{0x01}, []byte{0x02}, []byte{0x42})
	f.Fuzz(func(t *testing.T, keyBytes, missingBytes, valueBytes []byte) {
		key := fuzzHash(keyBytes)
		missing := fuzzHash(missingBytes)
		value := fuzzHash(valueBytes)
		if key == missing {
			return
		}
		tree := New()
		tree.Update(key, value)
		proof := tree.Prove(missing)
		if proof.Exists || !Verify(tree.Root(), missing, Hash{}, proof) {
			t.Fatal("valid non-inclusion proof rejected")
		}
		if Verify(tree.Root(), missing, value, proof) && value != (Hash{}) {
			t.Fatal("non-inclusion proof accepted with nonzero value")
		}
	})
}
