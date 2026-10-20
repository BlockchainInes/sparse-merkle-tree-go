package smt

import "testing"

func TestInclusionProof(t *testing.T) {
	tree := New()
	var key, value Hash
	key[31] = 0x01
	value[31] = 0x42
	tree.Update(key, value)
	proof := tree.Prove(key)
	if !proof.Exists {
		t.Fatal("expected inclusion proof")
	}
	if !Verify(tree.Root(), key, value, proof) {
		t.Fatal("valid inclusion proof rejected")
	}
}

func TestNonInclusionProof(t *testing.T) {
	tree := New()
	var key, value, missing Hash
	key[31] = 0x01
	value[31] = 0x42
	missing[31] = 0x02
	tree.Update(key, value)
	proof := tree.Prove(missing)
	if proof.Exists {
		t.Fatal("expected non-inclusion proof")
	}
	if !Verify(tree.Root(), missing, Hash{}, proof) {
		t.Fatal("valid non-inclusion proof rejected")
	}
}

func TestTamperedProof(t *testing.T) {
	tree := New()
	var key, value Hash
	key[31] = 0x01
	value[31] = 0x42
	tree.Update(key, value)
	proof := tree.Prove(key)
	proof.Siblings[0][0] ^= 0x01
	if Verify(tree.Root(), key, value, proof) {
		t.Fatal("tampered proof accepted")
	}
}

func TestIncorrectValue(t *testing.T) {
	tree := New()
	var key, value, incorrect Hash
	key[31] = 0x01
	value[31] = 0x42
	incorrect[31] = 0x43
	tree.Update(key, value)
	proof := tree.Prove(key)
	if Verify(tree.Root(), key, incorrect, proof) {
		t.Fatal("incorrect value accepted")
	}
}

func TestEmptyTreeNonInclusion(t *testing.T) {
	tree := New()
	var key Hash
	key[31] = 0x01
	proof := tree.Prove(key)
	if proof.Exists {
		t.Fatal("empty tree returned inclusion proof")
	}
	if !Verify(tree.Root(), key, Hash{}, proof) {
		t.Fatal("empty tree non-inclusion proof rejected")
	}
}
