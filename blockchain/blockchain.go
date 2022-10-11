package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
	lh     = "lh"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) AddBlock(data string) {
	//find the latest block
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(lh))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)
			return nil
		})
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte(lh), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)

}

// 创世纪块
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if item, err := txn.Get([]byte(lh)); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set([]byte(genesis.Hash), genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte(lh), genesis.Hash)

			lastHash = genesis.Hash
			return err
		} else {
			err := item.Value(func(val []byte) error {
				lastHash = append(lastHash, val...)
				return nil
			})
			return err

		}
	})
	Handle(err)
	return &BlockChain{lastHash, db}
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.Database}
}

func (it *BlockChainIterator) Next() *Block {
	var block *Block

	err := it.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(it.CurrentHash)
		Handle(err)
		var data []byte
		err = item.Value(func(val []byte) error {
			data = append(data, val...)
			block = Deserialize(data)
			return nil
		})
		return err
	})
	Handle(err)
	it.CurrentHash = block.PrevHash
	return block
}
