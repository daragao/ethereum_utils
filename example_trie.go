package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func generateProof(trie *trie.Trie, path []byte) []interface{} {
	proof := ethdb.NewMemDatabase()
	err := trie.Prove(path, 0, proof)
	if err != nil {
		log.Fatal("ERROR failed to create proof")
	}

	var proofArr []interface{}
	for nodeIt := trie.NodeIterator(nil); nodeIt.Next(true); {
		if val, err := proof.Get(nodeIt.Hash().Bytes()); val != nil && err == nil {
			var decodedVal interface{}
			err = rlp.DecodeBytes(val, &decodedVal)
			if err != nil {
				log.Fatalf("ERROR(%s) failed decoding RLP: 0x%0x\n", err, val)
			}
			proofArr = append(proofArr, decodedVal)
		}
	}
	return proofArr
}

func main() {
	trieDB := trie.NewDatabase(ethdb.NewMemDatabase())
	trieObj, _ := trie.New(common.Hash{}, trieDB) // empty trie

	data := map[string]string{
		"a":    "test1",
		"ab":   "t",
		"abc":  "test3",
		"abcd": "test4",
		"abed": "test5",
	}

	for key, value := range data {
		p, _ := rlp.EncodeToBytes(key)
		v, _ := rlp.EncodeToBytes(value)

		trieObj.Update(p, v) // update trie with the rlp encode index and the rlp encoded transaction
	}

	_, err := trieObj.Commit(nil) // commit to database (which in this case is stored in memory)
	if err != nil {
		log.Fatalf("commit error: %v", err)
	}

	fmt.Printf("Trie Root: %0x\n", trieObj.Hash().Bytes())

	fmt.Printf("Trie Leafs\n")
	it := trie.NewIterator(trieObj.NodeIterator(nil))
	for it.Next() {
		var decodedKey, decodedValue string
		rlp.DecodeBytes(it.Key, &decodedKey)
		rlp.DecodeBytes(it.Value, &decodedValue)

		//fmt.Printf("Key: % 0x\tValue: % 0x\n", it.Key, it.Value)
		fmt.Printf("\tKey: %s (%0x)\tValue: %s (%0x)\n", decodedKey, it.Key, decodedValue, it.Value)
	}

	fmt.Printf("\nTrie Nodes DB (RLP decoded and not ordered)\n")
	// print branches
	for _, node := range trieDB.Nodes() {
		dbNode, _ := trieDB.Node(node)
		var slice interface{}
		rlp.DecodeBytes(dbNode, &slice)
		fmt.Printf("\t%0x = SHA3(%0x) = SHA3(RLP(%0x))\n", node, dbNode, slice)
	}

	//generate a proof
	proofPath := "abcd"
	proofKey, _ := rlp.EncodeToBytes("abcd")
	proofDecoded := generateProof(trieObj, proofKey)
	proofEncoded, _ := rlp.EncodeToBytes(proofDecoded)
	fmt.Printf("\nProof of path %s (%0x)\n\tRLP Decoded: %0x\n\tRLP Encoded: %0x\n", proofPath, proofKey, proofDecoded, proofEncoded)

}

/*
Expected output:

Trie Root: da2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c
Trie Leafs
	Key: a (61)	Value: test1 (857465737431)
	Key: ab (826162)	Value: t (74)
	Key: abc (83616263)	Value: test3 (857465737433)
	Key: abcd (8461626364)	Value: test4 (857465737434)
	Key: abed (8461626564)	Value: test5 (857465737435)

Trie Nodes DB (RLP decoded and not ordered)
	6b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccc = SHA3(e5808080ca8220648685746573743480ca822064868574657374358080808080808080808080) = SHA3(RLP([   [2064 857465737434]  [2064 857465737435]           ]))
	5d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2 = SHA3(e583161626a06b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccc) = SHA3(RLP([161626 6b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccc]))
	207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a97 = SHA3(f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080) = SHA3(RLP([  [206162 74] [20616263 857465737433] 5d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2            ]))
	da2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c = SHA3(f839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080) = SHA3(RLP([      [31 857465737431]  207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a97        ]))

Proof of path abcd (8461626364)
	RLP Decoded: [[      [31 857465737431]  207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a97        ] [  [206162 74] [20616263 857465737433] 5d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2            ] [161626 6b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccc] [   [2064 857465737434]  [2064 857465737435]           ]]
	RLP Encoded: f8cbf839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080e583161626a06b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccce5808080ca8220648685746573743480ca822064868574657374358080808080808080808080

*/
