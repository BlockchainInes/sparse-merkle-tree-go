package smt

import "golang.org/x/crypto/sha3"

type Hash [32]byte

func keccak256(data []byte) Hash {
	h := sha3.NewLegacyKeccak256()
	_, _ = h.Write(data)
	var result Hash
	h.Sum(result[:0])
	return result
}

func hashLeaf(key, value Hash) Hash {
	var data [65]byte
	data[0] = 0x00
	copy(data[1:33], key[:])
	copy(data[33:], value[:])
	return keccak256(data[:])
}

func hashNode(left, right Hash) Hash {
	var data [65]byte
	data[0] = 0x01
	copy(data[1:33], left[:])
	copy(data[33:], right[:])
	return keccak256(data[:])
}

func emptyHashes() [257]Hash {
	var hashes [257]Hash
	hashes[0] = keccak256([]byte{0x02})
	for i := 1; i <= 256; i++ {
		hashes[i] = hashNode(hashes[i-1], hashes[i-1])
	}
	return hashes
}
