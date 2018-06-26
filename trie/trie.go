package trie

import (
	"fmt"
	"encoding/hex"
	"github.com/clearmatics/ion/go_util/rlp"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func toNibbleArray(bArr []byte) []byte {
	var nibble []byte
	for _, b := range bArr {
		nibble = append(nibble, b >> 4, b & 0x0f)
	}
	return nibble
}

func get(db map[string][][]byte, key []byte) [][]byte {
	return db[hex.EncodeToString(key)]
}
func put(db map[string][][]byte, key []byte ,value [][]byte) {
	db[hex.EncodeToString(key)] = value
}
func printDB(db map[string][][]byte) {
	for k, v := range db {
		//fmt.Printf("%s: %x\n",k,v)
		fmt.Printf("%s: [",k)
		for _, el := range v {
			fmt.Printf("%x, ",el)
		}
		fmt.Println("]")
	}
}
func printDumbTree(db map[string][][]byte, node []byte) {
	curNode := get(db,node)
	if curNode != nil {
		fmt.Printf("%x\n",curNode)
	}
	for _, el := range curNode {
		printDumbTree(db, el)
	}
}

func dumbUpdate(db map[string][][]byte, node,path,value []byte) []byte {
	var curNode, newNode [][]byte
	newNode = make([][]byte,17)
	if node == nil {
		curNode = make([][]byte,17)
	} else {
		// GET DATA FROM DB
		curNode = get(db,node)
	}
	copy(newNode, curNode)
	if path == nil || len(path) == 0 {
		// last element of the array is the value
		newNode[16] = value
	} else {
		newIndex := dumbUpdate(db,curNode[uint(path[0])],path[1:],value)
		newNode[uint(path[0])] = newIndex
	}

	// HASH
	buf := hashBytes(rlp.EncodeRLP(newNode))

	// INSERT DATA TO DB
	put(db,buf,newNode)

	return buf
}

func hashBytes(b []byte) []byte {
	hash := sha3.NewKeccak256()
	var buf []byte
	hash.Write(b)
	buf = hash.Sum(buf)
	return buf
}

func compactEncode(partial []byte) []byte {
	if len(partial) == 0 {
		return []byte{0}
	}
	term := 0
	// has terminator (0x10)
	if partial[len(partial)-1] == 0x10 {
		partial = partial[:len(partial)-1]
		term = 1
	}
	oddLen := len(partial) % 2
	flags := byte(2 * term + oddLen)
	if oddLen != 0 {
		return decodeNibble(append([]byte{flags}, partial...))
	} else {
		return decodeNibble(append([]byte{flags, 0}, partial...))
	}
}

// convert byte array into nibble only (to save space)
func decodeNibble(partial []byte) []byte {
	buf := make([]byte, len(partial)/2)
	for i := 0; i < len(buf); i += 1 {
		buf[i] |= partial[2*i] << 4 | partial[(2*i)+1]
	}
	return buf
}
// TODO: Optimized version of Trie
