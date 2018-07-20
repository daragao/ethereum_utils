package trie

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"reflect"

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

func get(db map[string]interface{}, key []byte) interface{} {
	return db[hex.EncodeToString(key)]
}
func put(db map[string]interface{}, key []byte, value interface{}) {
	db[hex.EncodeToString(key)] = value
}
func printDB(db map[string]interface{}) {
	for k, v := range db {
		//fmt.Printf("%s: %x\n",k,v)
		fmt.Printf("%s: [", k)
		for _, el := range v.([][]byte) {
			fmt.Printf("%x, ", el)
		}
		fmt.Println("]")
	}
}
func printDumbTree(db map[string]interface{}, node []byte) {
	curNode := get(db, node)
	if curNode != nil {
		fmt.Printf("%x\n", curNode)
	}
	for _, el := range curNode.([][]byte) {
		printDumbTree(db, el)
	}
}

func dumbUpdate(db map[string]interface{}, node, path, value []byte) []byte {
	var curNode, newNode [][]byte
	newNode = make([][]byte, 17)
	if node == nil {
		curNode = make([][]byte, 17)
	} else {
		// GET DATA FROM DB
		if get(db, node) != nil {
			curNode = get(db, node).([][]byte)
		}
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
	buf := hashBytes(rlp.Encode(newNode))

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

func compactDecode(nibbleArray []byte) []byte {
	bArray := fromNibble(nibbleArray)
	res := bArray
	flag := res[0]
	if len(nibbleArray) == 1 {
		if uint(nibbleArray[0]) == 0x00 {
			return nil
		}
		// nible should be 0x20 // TODO: add exception if not
		return []byte{0x10}
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
	for i := 0; i < len(buf); i++ {
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
func isLeaf(nibble interface{}) bool {
	return (nibble.([]byte)[0] >> 4) > 1
}

func flattenArray(in [][]byte) (out []byte) {
	for _, val := range in {
		if val == nil {
			out = append(out, byte(0x80))
		} else {
			out = append(out, val...)
		}
	}
	return
}

// TUpdate public
func TUpdate(db map[string]interface{}, node, path, value []byte) []byte {
	// only compact encodes path to mkae sure that recursive call account for odd length paths
	encodedPath := compactEncode(fromNibble(path))
	res := trieUpdate(db, node, encodedPath, value)

	resRef := reflect.ValueOf(res)

	if resRef.Len() < 32 {
		res = hashBytes(rlp.Encode(res))
	}
	return reflect.ValueOf(res).Bytes()
}

func hashAndSave(db map[string]interface{}, data interface{}) interface{} {
	dataRLP := rlp.Encode(data)
	dataHash := hashBytes(dataRLP)
	put(db, dataHash, data)
	// only saved as reference if rlp is longer than 32B!!!!!!
	if len(dataRLP) < 32 {
		//return dataRLP
		return data
	}
	return dataHash
}

// FIXME: need to make the array of array an interface type that way I can have multidimension arrays
// I DELETED IT ALL BECAUSE IT WAS CRAP!
// FIXME: Optimized version of Trie
func trieUpdate(db map[string]interface{}, node, encodedPath, value []byte) interface{} {
	// FIXME: please! need to find the right cases to use extension/branch/leaf
	// XXX: this code is "Seven" crime scene! "What's in the box?!" :(
	var curNode, newNode []interface{}

	if get(db, node) != nil {
		curNode = get(db, node).([]interface{})
	}

	pathBytes := compactDecode(encodedPath)

	if curNode == nil {
		newNode = make([]interface{}, 2)
	} else {
		newNode = make([]interface{}, len(curNode))
	}
	copy(newNode, curNode)

	switch len(curNode) {
	case 0:
		newNode[0] = compactEncode(append(pathBytes, byte(0x10)))
		newNode[1] = value
	case 2:
		curNodePathBytes := compactDecode(curNode[0].([]byte))
		if isLeaf(curNode[0]) {
			curNodePathBytes = curNodePathBytes[:len(curNodePathBytes)-1]
		}

		// find common prefix
		var commonPrefix []byte
		minLenPath := int(math.Min(float64(len(curNodePathBytes)), float64(len(pathBytes))))
		for i := 0; i < minLenPath; i++ {
			if curNodePathBytes[i] != pathBytes[i] {
				break
			}
			commonPrefix = append(commonPrefix, pathBytes[i])
		}

		if len(commonPrefix) == 0 {
			// no comon prefix
			// log.Fatalf("NOT IMPLEMENTED: no common prefix!\n\tNew Path: \t% 0x\n\tOld Path: \t% 0x\n", pathBytes, curNodePathBytes)

			// create new lower level branch node
			newBranch := make([]interface{}, 17)

			// if len(trimmedPathBytes) == 0 { newBranch[16] = value }
			if len(curNodePathBytes) == 0 || len(pathBytes) == 0 {
				log.Fatalf("NOT IMPLEMENTED:\n\tNew Path: \t% 0x\n\tOld Path: \t% 0x\n", pathBytes, curNodePathBytes)
			}

			// find index of old path and new path and remove first element of each
			idxPath := uint(pathBytes[0])
			idxCurPath := uint(curNodePathBytes[0])
			newEncodedPath := compactEncode(pathBytes[1:])
			var newEncodedCurPath []byte
			if len(newEncodedCurPath) != 0 {
				newEncodedCurPath = compactEncode(curNodePathBytes[1:])
			}

			// add new node
			lowerLevelNodeHashA := trieUpdate(db, nil, newEncodedPath, value)
			newBranch[idxPath] = lowerLevelNodeHashA // this creates a leaf

			// update current node to an extension that is pointed by this branch
			newNode[0] = newEncodedCurPath
			// save update node
			if isLeaf(curNode[0]) {
				newBranch[16] = curNode[1]
			} else {
				//newNodeHash := hashBytes(rlp.EncodeRLP(newNode))
				//put(db, newNodeHash, newNode)
				newNodeHash := hashAndSave(db, newNode)
				newBranch[idxCurPath] = newNodeHash
			}

			// update newNode to newBranch
			newNode = newBranch
		} else {
			// separate common prefix from paths
			trimmedPathBytes := bytes.TrimPrefix(pathBytes, commonPrefix)
			trimmedCurPathBytes := bytes.TrimPrefix(curNodePathBytes, commonPrefix)
			//log.Fatalf("NOT IMPLEMENTED: no common prefix!\n\tNew Path: \t% 0x\n\tOld Path: \t% 0x\n", pathBytes, curNodePathBytes)
			// there is a common prefix
			newEncodedCommonPrefix := compactEncode(commonPrefix)
			newEncodedPath := compactEncode(trimmedPathBytes[1:])
			var newEncodedCurPath []byte
			if len(newEncodedCurPath) != 0 {
				newEncodedCurPath = compactEncode(trimmedCurPathBytes[1:])
			}
			idxPath := uint(trimmedPathBytes[0])

			// create lower level branch node
			newBranch := make([]interface{}, 17)
			// add new path branches
			lowerLevelNodeHashA := trieUpdate(db, nil, newEncodedPath, value)
			newBranch[idxPath] = lowerLevelNodeHashA // this creates a leaf
			// add old path branches
			if isLeaf(curNode[0]) {
				newBranch[16] = curNode[1]
			} else {
				if len(trimmedCurPathBytes) == 0 {
					//log.Fatalf("NOT IMPLEMENTED: trimmeCurPAth == 0!\n\tNew Path: \t% 0x\n\tOld Path: \t% 0x\n", pathBytes, curNodePathBytes)
					newBranch[16] = curNode[1]
				} else {
					idxCurPath := uint(trimmedCurPathBytes[0])
					newBranch[idxCurPath] = curNode[1]
				}
			}
			// newBranchHash := hashBytes(rlp.EncodeRLP(newBranch))
			// put(db, newBranchHash, newBranch)
			newBranchHash := hashAndSave(db, newBranch)

			// create this level replacement extension node
			newNode[0] = newEncodedCommonPrefix
			newNode[1] = newBranchHash
		}
	case 17:
		if len(pathBytes) == 0 {
			newNode[16] = value
		} else {
			idx := uint(pathBytes[0])
			newPath := compactEncode(pathBytes[:len(pathBytes)-1])
			newNode[idx] = trieUpdate(db, curNode[idx].([]byte), newPath, value)
		}
	default:
		log.Fatal("ERROR: bad size of node!")
	}

	// HASH
	//buf := hashBytes(rlp.EncodeRLP(newNode))
	//put(db, buf, newNode)
	newNodeHash := hashAndSave(db, newNode)

	return newNodeHash
}
