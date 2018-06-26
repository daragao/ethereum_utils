package rlp

import (
	"fmt"
	"encoding/hex"
	"encoding/binary"
	"bytes"
	// "strings"
	"reflect"
)

func EncodeRLP(data interface{}) []byte {
	return encodeArrayRLP(data,true)
}

func encodeArrayRLP(data interface{}, isFirst bool) []byte {
	var result []byte

	rt := reflect.TypeOf(data)
	switch rt.Kind() {
	case reflect.Slice:
		if rt.Elem().Kind() == reflect.Uint8 {
			b := reflect.ValueOf(data).Bytes()
			return encodeBytesRLP(b)
		} else {
			s := reflect.ValueOf(data)
			dataLen := s.Len()

			for i := 0; i < dataLen; i++ {
				//isFirst = s.Index(i).Kind() == reflect.Slice
				result = append(result[:],encodeArrayRLP(s.Index(i).Interface(),false)[:]...)
			}
			//if isFirst {
				firstByteValue := 0xc0

				firstByte := int2Bytes(uint32(firstByteValue + len(result)))
				if len(result) > 55 {
					lenByte := int2Bytes(uint32(len(result)))
					firstByte = int2Bytes(uint32(firstByteValue + 55 + len(lenByte)))
					firstByte = append(firstByte, lenByte...)
				}
				result = append(firstByte,result[:]...)
			//}
			return result
		}
	case reflect.Bool:
		boolData := reflect.ValueOf(data).Bool() 
		if boolData {
			return []byte{0x01}
		} else {
			return []byte{0x80}
		}
	case reflect.String:
		return encodeBytesRLP([]byte(data.(string)))
	case reflect.Uint:
		b := uint32(reflect.ValueOf(data).Uint())
		return encodeBytesRLP(int2Bytes(b))
	default:
		b := reflect.ValueOf(data).Bytes()
		fmt.Println("UNKNOWN TYPE!",hex.EncodeToString(b),data,rt.Kind())
		return encodeBytesRLP(b)
	}

	return nil
}

func encodeBytesRLP(data []byte) []byte {
	dataLen := len(data)
	if dataLen == 1 {
		if data[0] <= 0x7f {
			return []byte{data[0]}
		} else {
			return []byte{0x81,data[0]}
		}
	}
	if dataLen > 1 && dataLen < 56 {
		return append([]byte{byte(0x80+dataLen)}[:],data[:]...)
	}
	if dataLen > 55 {
		bs := int2Bytes(uint32(dataLen))
		return append(append([]byte{byte(0xb7+len(bs))}[:],bs[:]...),data[:]...)
	}
	//if dataLen == 0 {
	return []byte{0x80}
	//}
}

func int2Bytes(data uint32) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, data)
	return bytes.TrimLeft(bs, "\x00")
}

/*
func main() {
	fmt.Println("testing trie")



	// TODO: encodeRLP accept any type
	str := "hello world"
	strB := bytes.TrimLeft([]byte(str), "\x00")

	fmt.Println(str," -> 0x",hex.EncodeToString(strB))

	fmt.Println(hex.EncodeToString(encodeRLP(str)))
	fmt.Println(hex.EncodeToString(encodeRLP(strB)))
	fmt.Println(hex.EncodeToString(encodeRLP(true)))
	fmt.Println(hex.EncodeToString(encodeRLP(false)))

	data := string(strings.Repeat("a", 1024))
	fmt.Println(hex.EncodeToString(encodeRLP(data)))
	fmt.Println(hex.EncodeToString(encodeRLP([]string{"hello","world"})))
	fmt.Println(hex.EncodeToString(encodeRLP([]string{"hello","world",data})))
}
*/


// TODO: write tests for the RLP encoder
// TODO: write RLP decoder
// TODO: write Trie functions
