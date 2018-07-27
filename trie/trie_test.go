package trie

import (
	// "encoding/hex"
	"fmt"
	"log"
	"testing"

	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

/*
func TestEncode(t *testing.T) {
	db := make(map[string][][]byte)

	root := dumbUpdate(db,nil,toNibbleArray([]byte("do")),[]byte("verb"))
	root = dumbUpdate(db,root,toNibbleArray([]byte("dog")),[]byte("puppy"))

	//printDB(db)
	//printDumbTree(db, root)
	//fmt.Printf("root: % x\n",root)
}

func TestEncode2(t *testing.T) {
	db := make(map[string][][]byte)
	//str := "dog"
	//resA := rlp.EncodeRLP(str)
	//fmt.Println(hex.EncodeToString(resA))
	root := dumbUpdate(db,nil,toNibbleArray([]byte{1,1,2}),compactEncode(rlp.EncodeRLP("hello")))

	//printDB(db)
	//printDumbTree(db, root)
	//fmt.Printf("root: % x\n",root)
}
*/

func TestCompactEncode(t *testing.T) {
	tests := []struct {
		hex, compact []byte
		isLeaf       bool
	}{
		// empty keys, with and without terminator.
		{hex: []byte{}, compact: []byte{0x00}, isLeaf: false},
		{hex: []byte{16}, compact: []byte{0x20}, isLeaf: true},
		// odd length, no terminator
		{hex: []byte{1, 2, 3, 4, 5}, compact: []byte{0x11, 0x23, 0x45}, isLeaf: false},
		// even length, no terminator
		{hex: []byte{0, 1, 2, 3, 4, 5}, compact: []byte{0x00, 0x01, 0x23, 0x45}, isLeaf: false},
		// odd length, terminator
		{hex: []byte{15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x3f, 0x1c, 0xb8}, isLeaf: true},
		// even length, terminator
		{hex: []byte{0, 15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x20, 0x0f, 0x1c, 0xb8}, isLeaf: true},
	}

	for _, test := range tests {
		if c := compactEncode(test.hex); !bytes.Equal(c, test.compact) {
			t.Errorf("hexToCompact(% 0x) -> % 0x, want % 0x", test.hex, c, test.compact)
		}
		if c := compactDecode(test.compact); !bytes.Equal(c, test.hex) {
			t.Errorf("compactToHex(% 0x) -> % 0x, want % 0x", test.compact, c, test.hex)
		}
		if l := isLeaf(test.compact); l != test.isLeaf {
			t.Errorf("isLeaf(% 0x) -> %v, want %v", test.compact, l, test.isLeaf)
		}
	}
}

func TestSingleInsert(t *testing.T) {
	key := []byte("A")
	data := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	expectedRootHash := []byte{0xd2, 0x37, 0x86, 0xfb, 0x4a, 0x01, 0x0d, 0xa3, 0xce, 0x63, 0x9d, 0x66, 0xd5, 0xe9, 0x04, 0xa1, 0x1d, 0xbc, 0x02, 0x74, 0x6d, 0x1c, 0xe2, 0x50, 0x29, 0xe5, 0x32, 0x90, 0xca, 0xbf, 0x28, 0xab}

	db := make(map[string][][]byte)
	rootHash := trieUpdate(db, nil, key, data)
	if !bytes.Equal(expectedRootHash, rootHash) {
		t.Errorf("~Failed single insert: \n\texpected:\t% 0x \n\tresult:\t\t% 0x\n", expectedRootHash, rootHash)
	}

	trieDB := trie.NewDatabase(ethdb.NewMemDatabase())
	trieObj, _ := trie.New(common.Hash{}, trieDB) // empty trie
	trieObj.Update(key, data)
	rootGoLibHash, err := trieObj.Commit(nil)
	if err != nil {
		log.Fatalf("commit error: %v", err)
	}

	dbNode, _ := trieDB.Node(trieDB.Nodes()[0])
	fmt.Printf("% 0x\n", dbNode)
	fmt.Printf("% 0x\n", key)
	fmt.Printf("% 0x\n", data)

	fmt.Printf("My root:\t% 0x\nGeth root:\t% 0x\n", rootHash, rootGoLibHash.Bytes())

	it := trie.NewIterator(trieObj.NodeIterator(nil))
	for it.Next() {
		fmt.Printf("key: %s\tvalue: %s\n", string(it.Key), string(it.Value))
	}
}

func TestInsert(t *testing.T) {
	values := [][][]byte{
		[][]byte{[]byte("A"), []byte("1")},
		[][]byte{[]byte("ABCD"), []byte("1")},
		[][]byte{[]byte("ABCDE"), []byte("1")},
		[][]byte{[]byte("B"), []byte("1")},
		[][]byte{[]byte("ABCDE"), []byte("1")},
	}

	db := make(map[string][][]byte)

	trieDB := trie.NewDatabase(ethdb.NewMemDatabase())
	trieObj, _ := trie.New(common.Hash{}, trieDB) // empty trie

	for _, v := range values {
		key := v[0]
		data := v[1]
		// my trie
		rootHash := trieUpdate(db, nil, key, data)

		// go ethereum trie
		trieObj.Update(key, data)
		rootGoLibHash, err := trieObj.Commit(nil)
		if err != nil {
			log.Fatalf("commit error: %v", err)
		}

		if !bytes.Equal(rootGoLibHash.Bytes(), rootHash) {
			t.Errorf("%s:%s\n\tMy Trie:\t % 0x\n\tGo Lib Trie:\t % 0x\n", key, data, rootHash, rootGoLibHash.Bytes())
		}
	}
}

/* TODO
func TestInsert(t *testing.T) {
	// values {key,data}
	values := [][][]byte {
		[][]byte{[]byte("doe"),[]byte("reindeer")},
		[][]byte{[]byte("dog"),[]byte("puppy")},
		[][]byte{[]byte("dogglesworth"),[]byte("cat")},
	}
	expectedRootHash := []byte{0x8a,0xad,0x78,0x9d,0xff,0x2f,0x53,0x8b,0xca,0x5d,0x8e,0xa5,0x6e,0x8a,0xbe,0x10,0xf4,0xc7,0xba,0x3a,0x5d,0xea,0x95,0xfe,0xa4,0xcd,0x6e,0x7c,0x3a,0x11,0x68,0xd3,}
}
*/
