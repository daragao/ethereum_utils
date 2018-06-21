package trie

import (
	"fmt"
	"encoding/hex"
	"github.com/clearmatics/ion/go_util/rlp"
)

func EncodeTrie() {
	str := "dog"

	resA := rlp.EncodeRLP(str)
	fmt.Println(hex.EncodeToString(resA))
}
