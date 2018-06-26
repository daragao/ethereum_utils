package main

import (
		"context"	
    "fmt"
    "log"

		"github.com/ethereum/go-ethereum/crypto"
		_ "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	//_ = client // we'll use this in the upcoming sections

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number.String()) // 5671744

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(count)
	fmt.Printf("\n%+v\n",header)
	rlpHeader, _ := rlp.EncodeToBytes(&header)
	headerHash := crypto.Keccak256Hash(rlpHeader)
	fmt.Printf("\n%x == %x\n",headerHash, block.Hash())


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



data := []byte("hello")
hash := crypto.Keccak256Hash(data)
fmt.Println("\"hello\" keccak hash:\t",hash.Hex())
}
