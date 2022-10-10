package main

import (
	"fmt"

	"github.com/aliensb/blockchain-programming/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()
	chain.AddBlock("tom")
	chain.AddBlock("Jerry")
	chain.AddBlock("Lily")
	chain.AddBlock("Jack")

	for _, block := range chain.Blocks {
		fmt.Printf("chain prevHash: %X \nchain hash: %X \nchain data: %s \n\n", block.PrevHash, block.Hash, block.Data)
	}
}
