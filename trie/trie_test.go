package trie

import (
	"testing"
)

func TestEncode(t *testing.T) {
	db := make(map[string]byte)
	db["test"] = byte(0)

	EncodeTrie()
}
