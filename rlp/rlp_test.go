package rlp

import (
	//"fmt"
	"bytes"
	//"encoding/hex"
	"strings"
	"github.com/ethereum/go-ethereum/rlp"
	"testing"
)

func TestEncode(t *testing.T) {
	str := "hello world"
	strB := bytes.TrimLeft([]byte(str), "\x00")

	resA := EncodeRLP(str)
	resB := EncodeRLP(strB)

	if  !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of byte vs string")
	}
}

type DataTestType struct {
	Data string
}

type DataTestMultiType struct {
	DataA string
	DataB uint
}

func TestCompareEncode(t *testing.T) {
	str := "hello world"
	resA, _ := rlp.EncodeToBytes(&DataTestType{Data: str})
	resB := EncodeRLP([]string{str})
	if  !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own")
	}

	str = string(strings.Repeat("a", 1024))
	resA, _ = rlp.EncodeToBytes(&DataTestType{Data: str})
	resB = EncodeRLP([]string{str})
	if  !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own, on long strings")
	}


	str = "hello world"
	val := uint(5)
	resA, _ = rlp.EncodeToBytes(&DataTestMultiType{DataA: str, DataB: val})
	//fmt.Println(hex.EncodeToString(resA))
	resB = EncodeRLP([]interface{}{str,val})
	//fmt.Println(hex.EncodeToString(resB))
	if  !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own, on long strings")
	}
}
