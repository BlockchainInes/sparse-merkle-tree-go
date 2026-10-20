 // SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {SparseMerkleVerifier} from "../contracts/SparseMerkleVerifier.sol";

contract SparseMerkleVerifierTest {
    SparseMerkleVerifier private verifier;

    bytes32 private constant KEY = bytes32(uint256(1));
    bytes32 private constant VALUE = bytes32(uint256(42));

    function setUp() public {
        verifier = new SparseMerkleVerifier();
    }

    function testValidInclusionProof() public {
        (bytes32 root, bytes32[256] memory siblings) =
            buildSingleLeafProof(KEY, VALUE);

        assert(
            verifier.verify(root, KEY, VALUE, true, siblings)
        );
    }

    function testValidNonInclusionProof() public {
        bytes32[256] memory siblings;
        bytes32 current = keccak256(hex"02");

        for (uint256 level; level < 256; ++level) {
            siblings[level] = current;
            current = hashNode(current, current);
        }

        assert(
            verifier.verify(
                current,
                KEY,
                bytes32(0),
                false,
                siblings
            )
        );
    }

    function testIncorrectValueRejected() public {
        (bytes32 root, bytes32[256] memory siblings) =
            buildSingleLeafProof(KEY, VALUE);

        assert(
            !verifier.verify(
                root,
                KEY,
                bytes32(uint256(43)),
                true,
                siblings
            )
        );
    }

    function testIncorrectRootRejected() public {
        (, bytes32[256] memory siblings) =
            buildSingleLeafProof(KEY, VALUE);

        assert(
            !verifier.verify(
                bytes32(uint256(999)),
                KEY,
                VALUE,
                true,
                siblings
            )
        );
    }

    function testTamperedSiblingRejected() public {
        (bytes32 root, bytes32[256] memory siblings) =
            buildSingleLeafProof(KEY, VALUE);

        siblings[128] = bytes32(uint256(123));

        assert(
            !verifier.verify(root, KEY, VALUE, true, siblings)
        );
    }

    function testNonInclusionWithNonzeroValueRejected() public {
        bytes32[256] memory siblings;

        assert(
            !verifier.verify(
                bytes32(0),
                KEY,
                VALUE,
                false,
                siblings
            )
        );
    }

    function buildSingleLeafProof(
        bytes32 key,
        bytes32 value
    )
        internal
        pure
        returns (
            bytes32 root,
            bytes32[256] memory siblings
        )
    {
        bytes32 current = keccak256(
            abi.encodePacked(bytes1(0x00), key, value)
        );

        bytes32 empty = keccak256(hex"02");
        uint256 keyBits = uint256(key);

        for (uint256 level; level < 256; ++level) {
            siblings[level] = empty;

            if (((keyBits >> level) & 1) == 0) {
                current = hashNode(current, empty);
            } else {
                current = hashNode(empty, current);
            }

            empty = hashNode(empty, empty);
        }

        root = current;
    }

    function hashNode(
        bytes32 left,
        bytes32 right
    ) internal pure returns (bytes32) {
        return keccak256(
            abi.encodePacked(bytes1(0x01), left, right)
        );
    }
}
