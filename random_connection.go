package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"

	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/clearmatics/ion/go_util/util"
	"github.com/ethereum/go-ethereum/ethdb"
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

func main() {
	client := util.Client("https://rinkeby.infura.io")

	fmt.Println("we have a connection")
	//_ = client // we'll use this in the upcoming sections

	blockNumber := big.NewInt(2657422)

	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	// trieObj := util.TxTrie(block.Transactions())

	receiptArr := util.GetBlockTxReceipts(client, block)
	receiptTrie := util.ReceiptTrie(receiptArr)

	// TODO: need to order the proof into and RLP
	proof := ethdb.NewMemDatabase()
	key := []byte{19}
	err = receiptTrie.Prove(key, 0, proof)
	if err != nil {
		log.Fatal("ERROR failed to create proof")
	}

	fmt.Printf("\nProof map for tx receipt with index 0x%0x (#tx 0x%0x)\n", key, block.Transactions()[19].Hash().Bytes())
	for _, key := range proof.Keys() {
		val, _ := proof.Get(key)
		fmt.Printf("\tkey (sha3(value)): 0x%0x\n\t value: 0x%0x\n\t\t\t===================================================\n", key, val)
	}
}
