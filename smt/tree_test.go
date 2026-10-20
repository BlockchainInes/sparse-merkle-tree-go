package smt

import "testing"

func TestEmptyTreeRoot(t *testing.T) {
	tree := New()
	expected := emptyHashes()[256]
	if tree.Root() != expected {
		t.Fatal("unexpected empty tree root")
	}
}

func TestUpdateChangesRoot(t *testing.T) {
	tree := New()
	var key, value Hash
	key[31] = 0x01
	value[31] = 0x42
	initialRoot := tree.Root()
	tree.Update(key, value)
	if tree.Root() == initialRoot {
		t.Fatal("root unchanged after update")
	}
}

func TestUpdateDeterministic(t *testing.T) {
	var key, value Hash
	key[31] = 0x01
	value[31] = 0x42
	first := New()
	second := New()
	first.Update(key, value)
	second.Update(key, value)
	if first.Root() != second.Root() {
		t.Fatal("identical updates produced different roots")
	}
}

func TestUpdateOrderIndependence(t *testing.T) {
	var keyA, keyB, valueA, valueB Hash
	keyA[31] = 0x01
	keyB[31] = 0x02
	valueA[31] = 0xAA
	valueB[31] = 0xBB
	first := New()
	second := New()
	first.Update(keyA, valueA)
	first.Update(keyB, valueB)
	second.Update(keyB, valueB)
	second.Update(keyA, valueA)
	if first.Root() != second.Root() {
		t.Fatal("update order changed the final root")
	}
}

func TestOverwriteDeterministic(t *testing.T) {
	var key, oldValue, newValue Hash
	key[31] = 0x01
	oldValue[31] = 0xAA
	newValue[31] = 0xBB
	first := New()
	second := New()
	first.Update(key, oldValue)
	first.Update(key, newValue)
	second.Update(key, newValue)
	if first.Root() != second.Root() {
		t.Fatal("overwriting a value produced an inconsistent root")
	}
}
