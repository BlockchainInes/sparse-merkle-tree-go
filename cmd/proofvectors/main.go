package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BlockchainInes/sparse-merkle-tree-go/smt"
)

type Vector struct {
	Name     string   `json:"name"`
	Root     string   `json:"root"`
	Key      string   `json:"key"`
	Value    string   `json:"value"`
	Exists   bool     `json:"exists"`
	Siblings []string `json:"siblings"`
}

func encodeHash(hash smt.Hash) string {
	return "0x" + hex.EncodeToString(hash[:])
}

func makeVector(
	name string,
	tree *smt.Tree,
	key smt.Hash,
	value smt.Hash,
) Vector {
	proof := tree.Prove(key)

	siblings := make([]string, len(proof.Siblings))
	for i, sibling := range proof.Siblings {
		siblings[i] = encodeHash(sibling)
	}

	return Vector{
		Name:     name,
		Root:     encodeHash(tree.Root()),
		Key:      encodeHash(key),
		Value:    encodeHash(value),
		Exists:   proof.Exists,
		Siblings: siblings,
	}
}

func main() {
	var keyA, keyB, missingKey smt.Hash
	var valueA, valueB smt.Hash

	keyA[31] = 0x01
	keyB[31] = 0x03
	missingKey[31] = 0x05

	valueA[31] = 0x2A
	valueB[31] = 0x7B

	tree := smt.New()
	tree.Update(keyA, valueA)
	tree.Update(keyB, valueB)

	vectors := []Vector{
		makeVector("inclusion-key-a", tree, keyA, valueA),
		makeVector("inclusion-key-b", tree, keyB, valueB),
		makeVector("non-inclusion", tree, missingKey, smt.Hash{}),
	}

	outputPath := filepath.Join("test", "vectors", "go_proofs.json")

	data, err := json.MarshalIndent(vectors, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %d Go proof vectors: %s\n", len(vectors), outputPath)
}
