package main

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"time"

	// "github.com/ethereum/go-ethereum/crypto"
	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/clearmatics/ion/go_util/util"
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

	// get a block
	blockNumber := big.NewInt(2657422)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// generate Transaction Trie
	txArr := block.Transactions()
	txTrie := util.TxTrie(txArr)

	// generate RLP Proof of Receipt
	txKey := []byte{19}
	txProofArr := util.Proof(txTrie, txKey)
	log.Printf("RLP Proof of Transaction with index %d:\n\t%0x\n", txKey, txProofArr)

	// generate Receipt Trie
	receiptArr := util.GetBlockTxReceipts(client, block)
	receiptTrie := util.ReceiptTrie(receiptArr)

	// generate RLP Proof of Receipt
	receiptKey := []byte{19}
	receiptProofArr := util.Proof(receiptTrie, receiptKey)
	log.Printf("RLP Proof of receipt with index %d:\n\t%0x\n", receiptKey, receiptProofArr)
}
