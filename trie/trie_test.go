package trie

import (
	//"fmt"
	// "encoding/hex"
	// "github.com/clearmatics/ion/go_util/rlp"
	"testing"

	"bytes"
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
	tests := []struct{ hex, compact []byte; isLeaf bool }{
		// empty keys, with and without terminator.
		{hex: []byte{}, compact: []byte{0x00}, isLeaf: false },
		{hex: []byte{16}, compact: []byte{0x20}, isLeaf: true },
		// odd length, no terminator
		{hex: []byte{1, 2, 3, 4, 5}, compact: []byte{0x11, 0x23, 0x45}, isLeaf: false },
		// even length, no terminator
		{hex: []byte{0, 1, 2, 3, 4, 5}, compact: []byte{0x00, 0x01, 0x23, 0x45}, isLeaf: false },
		// odd length, terminator
		{hex: []byte{15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x3f, 0x1c, 0xb8}, isLeaf: true },
		// even length, terminator
		{hex: []byte{0, 15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x20, 0x0f, 0x1c, 0xb8}, isLeaf: true },
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
