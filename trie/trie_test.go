package trie

import (
	"fmt"
	// "encoding/hex"
	"github.com/clearmatics/ion/go_util/rlp"
	"testing"

	"bytes"
)

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

	printDB(db)
	printDumbTree(db, root)
	fmt.Printf("root: % x\n",root)
}

func TestCompactEncode(t *testing.T) {
	partial := []byte{1,2,3,4,5}
	res := compactEncode(partial)
	fmt.Printf("% x\n",res)

	partial = []byte{0, 1,2,3,4,5}
	res = compactEncode(partial)
	fmt.Printf("% x\n",res)

	partial = []byte{0, 0xf, 1, 0xc, 0xb, 8, 0x10}
	res = compactEncode(partial)
	fmt.Printf("% x\n",res)
	
	partial = []byte{0xf, 1, 0xc, 0xb, 8, 0x10}
	res = compactEncode(partial)
	fmt.Printf("% x\n",res)

	tests := []struct{ hex, compact []byte }{
		// empty keys, with and without terminator.
		{hex: []byte{}, compact: []byte{0x00}},
		{hex: []byte{16}, compact: []byte{0x20}},
		// odd length, no terminator
		{hex: []byte{1, 2, 3, 4, 5}, compact: []byte{0x11, 0x23, 0x45}},
		// even length, no terminator
		{hex: []byte{0, 1, 2, 3, 4, 5}, compact: []byte{0x00, 0x01, 0x23, 0x45}},
		// odd length, terminator
		{hex: []byte{15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x3f, 0x1c, 0xb8}},
		// even length, terminator
		{hex: []byte{0, 15, 1, 12, 11, 8, 16 /*term*/}, compact: []byte{0x20, 0x0f, 0x1c, 0xb8}},
	}

	for _, test := range tests {
		if c := compactEncode(test.hex); !bytes.Equal(c, test.compact) {
			t.Errorf("hexToCompact(%x) -> %x, want %x", test.hex, c, test.compact)
		}
	}
}
