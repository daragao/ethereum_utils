package trie

import (
	// "encoding/hex"

	"github.com/clearmatics/ion/go_util/rlp"

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

	// my trie
	db := make(map[string][][]byte)
	rootHash := TrieUpdate(db, nil, key, data)

	// go-ethereum trie
	trieDB := trie.NewDatabase(ethdb.NewMemDatabase())
	trieObj, _ := trie.New(common.Hash{}, trieDB) // empty trie
	trieObj.Update(key, data)
	rootGoLibHash, err := trieObj.Commit(nil)
	if err != nil {
		log.Fatalf("commit error: %v", err)
	}

	// compare
	if !bytes.Equal(rootGoLibHash.Bytes(), rootHash) {
		t.Errorf("~Failed single insert: \n\texpected:\t% 0x \n\tresult:\t\t% 0x\n", rootGoLibHash.Bytes(), rootHash)
	}

	// get a node
	// dbNode, _ := trieDB.Node(trieDB.Nodes()[0])
	// fmt.Printf("Node: \t% 0x\n", dbNode)
	// fmt.Printf("Key: \t% 0x\n", key)
	// fmt.Printf("Value: \t% 0x\n", data)
	// fmt.Printf("My root:\t% 0x\nGeth root:\t% 0x\n", rootHash, rootGoLibHash.Bytes())

	/*
		// iterate through node values
			it := trie.NewIterator(trieObj.NodeIterator(nil))
			for it.Next() {
				fmt.Printf("key: %s\tvalue: %s\n", string(it.Key), string(it.Value))
			}
	*/
}

func printMyNodes(t *testing.T, tDb map[string][][]byte) {
	t.Errorf("My Nodes\n")
	for k, v := range tDb {
		rlpV := rlp.EncodeRLP(v)
		t.Errorf("\tNode[%s]: \t% 0x\n", k, rlpV)
	}
}

func printGoEthereumNodes(t *testing.T, tDb *trie.Database) {
	t.Errorf("Go-Ethereum Nodes\n")
	for idx, node := range tDb.Nodes() {
		dbNode, err := tDb.Node(node)
		if err == nil {
			t.Errorf("\tNode[%x]: \t% 0x\n", node.Bytes(), dbNode)
		} else {
			t.Errorf("\tERROR: Node[%0d]: \t%s\n", idx, err)
		}
	}
}

func TestInsert(t *testing.T) {
	values := [][][]byte{
		[][]byte{[]byte("A"), []byte("1")},
		[][]byte{[]byte("ABCD"), []byte("1")},
		[][]byte{[]byte("ABCDE"), []byte("1")},
		[][]byte{[]byte("B"), []byte("1")},
		[][]byte{[]byte("ABCDE"), []byte("1")},
		[][]byte{[]byte("ABC"), []byte("1")},
	}

	var rootHash []byte
	db := make(map[string][][]byte)

	trieDB := trie.NewDatabase(ethdb.NewMemDatabase())
	trieObj, _ := trie.New(common.Hash{}, trieDB) // empty trie

	for _, v := range values {
		key := v[0]
		data := v[1]
		// my trie
		rootHash = TrieUpdate(db, rootHash, key, data)

		// go ethereum trie
		trieObj.Update(key, data)
		rootGoLibHash, err := trieObj.Commit(nil)
		if err != nil {
			log.Fatalf("commit error: %v", err)
		}

		if !bytes.Equal(rootGoLibHash.Bytes(), rootHash) {
			t.Errorf("%s:%s\n\tMy Trie:\t % 0x\n\tGo Lib Trie:\t % 0x\n", key, data, rootHash, rootGoLibHash.Bytes())
			printMyNodes(t, db)
			printGoEthereumNodes(t, trieDB)
			t.Fatal()
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
