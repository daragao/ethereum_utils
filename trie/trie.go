package trie

import (
	"encoding/hex"
	"fmt"
	"github.com/clearmatics/ion/go_util/rlp"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func toNibbleArray(bArr []byte) []byte {
	var nibble []byte
	for _, b := range bArr {
		nibble = append(nibble, b>>4, b&0x0f)
	}
	return nibble
}

func get(db map[string][][]byte, key []byte) [][]byte {
	return db[hex.EncodeToString(key)]
}
func put(db map[string][][]byte, key []byte, value [][]byte) {
	db[hex.EncodeToString(key)] = value
}
func printDB(db map[string][][]byte) {
	for k, v := range db {
		//fmt.Printf("%s: %x\n",k,v)
		fmt.Printf("%s: [", k)
		for _, el := range v {
			fmt.Printf("%x, ", el)
		}
		fmt.Println("]")
	}
}
func printDumbTree(db map[string][][]byte, node []byte) {
	curNode := get(db, node)
	if curNode != nil {
		fmt.Printf("%x\n", curNode)
	}
	for _, el := range curNode {
		printDumbTree(db, el)
	}
}

func dumbUpdate(db map[string][][]byte, node, path, value []byte) []byte {
	var curNode, newNode [][]byte
	newNode = make([][]byte, 17)
	if node == nil {
		curNode = make([][]byte, 17)
	} else {
		// GET DATA FROM DB
		curNode = get(db, node)
	}
	copy(newNode, curNode)
	if path == nil || len(path) == 0 {
		// last element of the array is the value
		newNode[16] = value
	} else {
		newIndex := dumbUpdate(db, curNode[uint(path[0])], path[1:], value)
		newNode[uint(path[0])] = newIndex
	}

	// HASH
	buf := hashBytes(rlp.EncodeRLP(newNode))

	// INSERT DATA TO DB
	put(db, buf, newNode)

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
	flags := byte(2*term + oddLen)
	var result []byte
	if oddLen != 0 {
		result = append([]byte{flags}, partial...)
	} else {
		result = append([]byte{flags, 0}, partial...)
	}
	return toNibble(result)
}

// decode compact
func compactDecode(nibbleArray []byte) []byte {
	// TODO: the inverse of compactEncode
	bArray := fromNibble(nibbleArray)
	res := bArray
	flag := res[0]
	if len(nibbleArray) == 1 {
		if uint(nibbleArray[0]) == 0x00 {
			return nil
		} else { // nible should be 0x20 // TODO: add exception if not
			return []byte{0x10}
		}
	}
	if uint(flag) > 1 { // it should be 2 or 3
		res = append(res, byte(0x10)) // remove terminator
	}
	if uint(flag)%2 == 0 { // it should be either 0 or 2
		res = res[2:]
	} else { // it should be 1 or 3
		res = res[1:]
	}
	// fmt.Printf("bArray= % 0x fromNibble=% 0x flag=% 0x res=% 0x\n",nibbleArray,bArray,flag,res)
	return res
}

// convert byte array into nibble only (to save space)
func toNibble(partial []byte) []byte {
	buf := make([]byte, len(partial)/2)
	for i := 0; i < len(buf); i += 1 {
		buf[i] |= partial[2*i]<<4 | partial[(2*i)+1]
	}
	return buf
}
func fromNibble(partial []byte) []byte {
	buf := make([]byte, len(partial)*2)
	for i := 0; i < len(buf); i += 2 {
		buf[i] = partial[i/2] >> 4
		buf[i+1] = partial[(i/2)] & byte(0x0f)
		//fmt.Printf("% 0x -> % 0x",partial[i/2],buf[i:i+2])
	}
	// fmt.Printf("\n")
	return buf
}

// FIXME: Optimized version of Trie
func trieUpdate(db map[string][][]byte, node, path, value []byte) []byte {
	var newNode, curNode [][]byte
	if node == nil {
		curNode = make([][]byte, 2)
	} else {
		// GET DATA FROM DB
		curNode = get(db, node)
	}
	copy(newNode, curNode)

	// leaf nodes can be converted into extension nodes, and extension nodes can be converted into branch nodes (but not reversable)
	// leaf node -> extension node -> branch node
	if node == nil || len(node) == 0 {
		// LEAF NODE
		newNode[0] = compactEncode(append(path, byte(0x10))) // TODO: encode like it should be!
		newNode[1] = value
	} else if len(node) == 2 {
		// TODO: EXTENSION NODE
		// TODO: convert leaf node to extension node if needed
	} else {
		// TODO: convert extension node to branch node if needed
		// BRANCH NODE
		if path == nil || len(path) == 0 {
			// if the path has ended than this is the value
			// last element of the array is the value
			newNode[16] = value
		} else {
			newIndex := dumbUpdate(db, newNode[uint(path[0])], path[1:], value)
			newNode[uint(path[0])] = newIndex
		}
	}

	// HASH
	buf := hashBytes(rlp.EncodeRLP(newNode))

	// INSERT DATA TO DB
	put(db, buf, newNode)

	return buf
}
