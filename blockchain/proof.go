package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 12

type ProofOfWork struct {
	Block *Block
	//Target is a number that start with multiple zeros
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.PrevHash,
		pow.Block.Data,
		ToHex(int64(nonce)),
		ToHex(int64(Difficulty)),
	}, []byte{})
	return data
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intiHash big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x -> %d", hash, nonce)
		intiHash.SetBytes(hash[:])
		//if initHash is less than target, then initHash has meat the required zeros at start
		if intiHash.Cmp(pow.Target) == -1 {
			fmt.Printf("\ntarget is %x", pow.Target)
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intiHash big.Int
	var hash [32]byte

	data := pow.InitData(pow.Block.Nonce)
	hash = sha256.Sum256(data)
	fmt.Printf("\r%x", hash)
	intiHash.SetBytes(hash[:])
	return intiHash.Cmp(pow.Target) == -1
}
