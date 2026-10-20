package smt

import (
	"encoding/hex"
	"testing"
)

func TestKeccak256ReferenceVectors(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
		{
			name:     "abc",
			input:    []byte("abc"),
			expected: "4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, err := hex.DecodeString(tt.expected)
			if err != nil {
				t.Fatal(err)
			}
			actual := keccak256(tt.input)
			if hex.EncodeToString(actual[:]) != hex.EncodeToString(expected) {
				t.Fatalf("unexpected Keccak-256 hash for %s", tt.name)
			}
		})
	}
}

func TestSingleLeafRootReference(t *testing.T) {
	var key Hash
	var value Hash

	key[31] = 0x01
	value[31] = 0x2A

	tree := New()
	tree.Update(key, value)

	leafData := make([]byte, 65)
	leafData[0] = 0x00
	copy(leafData[1:33], key[:])
	copy(leafData[33:], value[:])

	current := keccak256(leafData)

	empty := keccak256([]byte{0x02})

	for level := 0; level < 256; level++ {
		var nodeData [65]byte
		nodeData[0] = 0x01

		if (key[31-level/8]>>(level%8))&1 == 0 {
			copy(nodeData[1:33], current[:])
			copy(nodeData[33:], empty[:])
		} else {
			copy(nodeData[1:33], empty[:])
			copy(nodeData[33:], current[:])
		}

		current = keccak256(nodeData[:])

		var emptyData [65]byte
		emptyData[0] = 0x01
		copy(emptyData[1:33], empty[:])
		copy(emptyData[33:], empty[:])
		empty = keccak256(emptyData[:])
	}

	if tree.Root() != current {
		t.Fatalf(
			"single-leaf root mismatch: got %x, expected %x",
			tree.Root(),
			current,
		)
	}
}

func TestMultipleLeafRootReference(t *testing.T) {
	var leftKey, rightKey Hash
	var leftValue, rightValue Hash

	leftKey[31] = 0x01
	rightKey[0] = 0x80
	rightKey[31] = 0x02

	leftValue[31] = 0x2A
	rightValue[31] = 0x7B

	tree := New()
	tree.Update(leftKey, leftValue)
	tree.Update(rightKey, rightValue)

	referenceLeaf := func(key, value Hash) Hash {
		var data [65]byte
		data[0] = 0x00
		copy(data[1:33], key[:])
		copy(data[33:], value[:])
		return keccak256(data[:])
	}

	referenceNode := func(left, right Hash) Hash {
		var data [65]byte
		data[0] = 0x01
		copy(data[1:33], left[:])
		copy(data[33:], right[:])
		return keccak256(data[:])
	}

	empty := keccak256([]byte{0x02})

	left := referenceLeaf(leftKey, leftValue)
	right := referenceLeaf(rightKey, rightValue)

	for level := 0; level < 255; level++ {
		if (leftKey[31-level/8]>>(level%8))&1 == 0 {
			left = referenceNode(left, empty)
		} else {
			left = referenceNode(empty, left)
		}

		if (rightKey[31-level/8]>>(level%8))&1 == 0 {
			right = referenceNode(right, empty)
		} else {
			right = referenceNode(empty, right)
		}

		empty = referenceNode(empty, empty)
	}

	expected := referenceNode(left, right)

	if tree.Root() != expected {
		t.Fatalf(
			"multiple-leaf root mismatch: got %x, expected %x",
			tree.Root(),
			expected,
		)
	}
}

func TestSharedPrefixRootReference(t *testing.T) {
	var keyA, keyB Hash
	var valueA, valueB Hash

	keyA[31] = 0x01
	keyB[31] = 0x03

	valueA[31] = 0x2A
	valueB[31] = 0x7B

	tree := New()
	tree.Update(keyA, valueA)
	tree.Update(keyB, valueB)

	referenceLeaf := func(key, value Hash) Hash {
		var data [65]byte
		data[0] = 0x00
		copy(data[1:33], key[:])
		copy(data[33:], value[:])
		return keccak256(data[:])
	}

	referenceNode := func(left, right Hash) Hash {
		var data [65]byte
		data[0] = 0x01
		copy(data[1:33], left[:])
		copy(data[33:], right[:])
		return keccak256(data[:])
	}

	empty := keccak256([]byte{0x02})

	leafA := referenceLeaf(keyA, valueA)
	leafB := referenceLeaf(keyB, valueB)

	branchA := referenceNode(empty, leafA)
	branchB := referenceNode(empty, leafB)

	current := referenceNode(branchA, branchB)

	empty = referenceNode(empty, empty)
	empty = referenceNode(empty, empty)

	for level := 2; level < 256; level++ {
		current = referenceNode(current, empty)
		empty = referenceNode(empty, empty)
	}

	if tree.Root() != current {
		t.Fatalf(
			"shared-prefix root mismatch: got %x, expected %x",
			tree.Root(),
			current,
		)
	}
}
