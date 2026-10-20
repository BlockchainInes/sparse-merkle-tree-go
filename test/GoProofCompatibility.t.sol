// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {SparseMerkleVerifier} from "../contracts/SparseMerkleVerifier.sol";

contract GoProofCompatibilityTest is Test {
    SparseMerkleVerifier private verifier;

    struct Vector {
        string name;
        bytes32 root;
        bytes32 key;
        bytes32 value;
        bool exists;
        bytes32[] siblings;
    }

    function setUp() public {
        verifier = new SparseMerkleVerifier();
    }

    function testGoInclusionProofs() public view {
        Vector[] memory vectors = loadVectors();

        for (uint256 i; i < vectors.length; ++i) {
            if (!vectors[i].exists) {
                continue;
            }

            bytes32[256] memory siblings = toFixedArray(
                vectors[i].siblings
            );

            assertTrue(
                verifier.verify(
                    vectors[i].root,
                    vectors[i].key,
                    vectors[i].value,
                    true,
                    siblings
                ),
                vectors[i].name
            );
        }
    }

    function testGoNonInclusionProof() public view {
        Vector[] memory vectors = loadVectors();

        for (uint256 i; i < vectors.length; ++i) {
            if (vectors[i].exists) {
                continue;
            }

            bytes32[256] memory siblings = toFixedArray(
                vectors[i].siblings
            );

            assertTrue(
                verifier.verify(
                    vectors[i].root,
                    vectors[i].key,
                    vectors[i].value,
                    false,
                    siblings
                ),
                vectors[i].name
            );
        }
    }

    function testGoProofTamperingRejected() public view {
        Vector[] memory vectors = loadVectors();

        for (uint256 i; i < vectors.length; ++i) {
            if (!vectors[i].exists) {
                continue;
            }

            bytes32[256] memory siblings = toFixedArray(
                vectors[i].siblings
            );

            siblings[128] = bytes32(
                uint256(siblings[128]) ^ uint256(1)
            );

            assertFalse(
                verifier.verify(
                    vectors[i].root,
                    vectors[i].key,
                    vectors[i].value,
                    true,
                    siblings
                ),
                vectors[i].name
            );
        }
    }

    function loadVectors() internal view returns (Vector[] memory) {
        string memory path = string.concat(
            vm.projectRoot(),
            "/test/vectors/go_proofs.json"
        );

        string memory json = vm.readFile(path);
        return abi.decode(vm.parseJson(json), (Vector[]));
    }

    function toFixedArray(
        bytes32[] memory source
    ) internal pure returns (bytes32[256] memory result) {
        require(source.length == 256, "Invalid proof length");

        for (uint256 i; i < 256; ++i) {
            result[i] = source[i];
        }
    }
}
