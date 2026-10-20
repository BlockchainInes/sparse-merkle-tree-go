 // SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract SparseMerkleVerifier {
    function verify(
        bytes32 root,
        bytes32 key,
        bytes32 value,
        bool exists,
        bytes32[256] calldata siblings
    ) external pure returns (bool) {
        bytes32 current;

        if (exists) {
            current = keccak256(
                abi.encodePacked(bytes1(0x00), key, value)
            );
        } else {
            if (value != bytes32(0)) {
                return false;
            }

            current = keccak256(hex"02");
        }

        uint256 keyBits = uint256(key);

        for (uint256 level = 0; level < 256; ++level) {
            bytes32 sibling = siblings[level];

            if (((keyBits >> level) & 1) == 0) {
                current = keccak256(
                    abi.encodePacked(bytes1(0x01), current, sibling)
                );
            } else {
                current = keccak256(
                    abi.encodePacked(bytes1(0x01), sibling, current)
                );
            }
        }

        return current == root;
    }
}
