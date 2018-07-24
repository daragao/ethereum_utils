package util

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client gets client or fails if no connection
func Client(url string) *ethclient.Client {
	client, err := ethclient.Dial(url)
	log.Println(client)
	if err != nil {
		log.Fatal("Client failed to connect: ", err)
	} else {
		fmt.Println("Connected to: ", url)
	}
	return client
}

// GetBlockTxReceipts get the receipts for all the transactions in a block
func GetBlockTxReceipts(ec *ethclient.Client, block *types.Block) []*types.Receipt {
	var receiptsArr []*types.Receipt
	for _, tx := range block.Transactions() {
		receipt, err := ec.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal("TransactionReceipt ERROR:", err)
		}
		receiptsArr = append(receiptsArr, receipt)
	}
	return receiptsArr
}
