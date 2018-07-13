package rlp

import (
	//"fmt"
	"bytes"
	//"encoding/hex"
	"github.com/ethereum/go-ethereum/rlp"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	str := "hello world"
	strB := bytes.TrimLeft([]byte(str), "\x00")

	resA := Encode(str)
	resB := Encode(strB)

	if !bytes.Equal(resA, resB) {
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
	resB := Encode([]string{str})
	if !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own")
	}

	str = string(strings.Repeat("a", 1024))
	resA, _ = rlp.EncodeToBytes(&DataTestType{Data: str})
	resB = Encode([]string{str})
	if !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own, on long strings")
	}

	str = "hello world"
	val := uint(5)
	resA, _ = rlp.EncodeToBytes(&DataTestMultiType{DataA: str, DataB: val})
	//fmt.Println(hex.EncodeToString(resA))
	resB = Encode([]interface{}{str, val})
	//fmt.Println(hex.EncodeToString(resB))
	if !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own, on long strings")
	}

	strArr := [][]string{{"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}, {"asdf", "qwer", "zxcv"}}
	resA, _ = rlp.EncodeToBytes(strArr)
	resB = Encode(strArr)
	if !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own, on array of arrays of strings", resA, resB)
	}

	strArr1 := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii", "jjj", "kkk", "lll", "mmm", "nnn", "ooo"}
	resA, _ = rlp.EncodeToBytes(strArr1)
	resB = Encode(strArr1)
	if !bytes.Equal(resA, resB) {
		t.Error("Bad encoding of geth vs my own, on array strings\n", "\n", resA, "\n", resB)
	}
}
