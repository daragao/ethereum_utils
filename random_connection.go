package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func printGoEthereumNodes(tDb *trie.Database) {
	log.Printf("Go-Ethereum Nodes\n")
	for idx, node := range tDb.Nodes() {
		dbNode, err := tDb.Node(node)
		if err == nil {
			log.Printf("\tNode[%x]: \t% 0x\n", node.Bytes(), dbNode)
		} else {
			log.Printf("\tERROR: Node[%0d]: \t%s\n", idx, err)
		}
	}
}

func getTxTrie(block *types.Block) *trie.Trie {
	trieDB := trie.NewDatabase(ethdb.NewMemDatabase())
	trieObj, _ := trie.New(common.Hash{}, trieDB) // empty trie
	for idx, tx := range block.Transactions() {

		rlpIdx, _ := rlp.EncodeToBytes(uint(idx))  // rlp encode index of transaction
		rlpTransaction, _ := rlp.EncodeToBytes(tx) // rlp encode transaction

		trieObj.Update(rlpIdx, rlpTransaction) // update trie with the rlp encode index and the rlp encoded transaction
		root, err := trieObj.Commit(nil)       // commit to database (which in this case is stored in memory)
		if err != nil {
			log.Fatalf("commit error: %v", err)
		}

		txRlpHash := crypto.Keccak256Hash(rlpTransaction)

		fmt.Printf("TxHash[%d]: \t% 0x\n\tHash(RLP(Tx)): \t% 0x\n\tTrieRoot: \t% 0x\n",
			idx, tx.Hash().Bytes(), txRlpHash.Bytes(), root.Bytes())
		//fmt.Printf("\n%#v\n% #v\n% 0x\n\n", trieObj, trieObj, root)
	}

	fmt.Printf("\n\nBlock number: %d \n\tBlock.TxHash:\t% 0x \n\tTrie.Root:\t% 0x\n",
		block.Number, block.TxHash().Bytes(), trieObj.Root())

	//printGoEthereumNodes(trieDB)

	return trieObj
}

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	//_ = client // we'll use this in the upcoming sections

	// blockNumber := big.NewInt(7295200)

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	/*
		fmt.Println(header.Number.String()) // 5671744

		count, err := client.TransactionCount(context.Background(), block.Hash())
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*
		// hash header
		fmt.Println(count)
		fmt.Printf("\n%+v\n", header)
		rlpHeader, _ := rlp.EncodeToBytes(&header)
		headerHash := crypto.Keccak256Hash(rlpHeader)
		fmt.Printf("\n%x == %x\n", headerHash, block.Hash())
	*/

	/*
			fmt.Println("==== // ===")
		for _, tx := range block.Transactions() {
			fmt.Println("HASH:\t",tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
			fmt.Println("VALUE:\t",tx.Value().String())    // 10000000000000000
			fmt.Println("GAS:\t",tx.Gas())               // 105000
			fmt.Println("GAS PRICE:\t",tx.GasPrice().Uint64()) // 102000000000
			fmt.Println("NONCE:\t",tx.Nonce())             // 110644
			fmt.Println("DATA:\t",tx.Data())              // []
			fmt.Println("TO:\t",tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
			fmt.Println("==== // ===")
		}
	*/

	trieObj := getTxTrie(block)
	txIdx := uint(random(0, len(block.Transactions())))
	rlpIdx, _ := rlp.EncodeToBytes(txIdx)
	txRlpBytes := trieObj.Get(rlpIdx)
	txRlpHash := crypto.Keccak256Hash(txRlpBytes).Bytes()

	fmt.Printf("Retrieval from Trie == Block TxHash\n\tHash(Trie.Get(%d)): \t% 0x\n\tBlock.TxHash[%d]: \t% 0x\n",
		txIdx, txRlpHash, txIdx, block.Transactions()[txIdx].Hash().Bytes())
	/*
		data := []byte("hello")
		hash := crypto.Keccak256Hash(data)
		fmt.Println("\"hello\" keccak hash:\t", hash.Hex())
	*/
}
