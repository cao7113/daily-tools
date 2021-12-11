package lib

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// SplitToGroupChunks split a number-string into group chunks
func SplitToGroupChunks(num string, g int, padding rune) []string {
	var chunks []string
	l := len(num)
	r := l % g
	if r != 0 {
		c0 := num[:r]
		if padding != 0 {
			c0 = strings.Repeat(string(padding), g-r) + c0
		}
		chunks = append(chunks, c0)
	}
	for i := 0; i < l/g; i++ {
		fi := r + i*g
		chunks = append(chunks, num[fi:fi+g])
	}
	return chunks
}

// GetIntBytes get bytes for int64 number
func GetIntBytes(n int64, be, isBit64 bool) []byte {
	var bo binary.ByteOrder
	if be {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}
	bs := make([]byte, 8)
	bo.PutUint64(bs, uint64(n))

	if isBit64 {
		return bs
	}

	var topSignByte byte = 0x00
	if n < 0 {
		topSignByte = 0xff
	}
	l2r := be
	i := FindFirstByteIndex(bs, l2r, func(b byte) bool {
		return b != topSignByte
	})
	if i == -1 { // return first byte if all bytes are same with top-flag byte
		return bs[:1]
	}

	si := i
	hv := 0x80 & bs[i] // top-bit of critical byte
	if hv != (0x80 & topSignByte) {
		si -= 1
	}

	if be {
		bs = bs[si:]
	} else {
		bs = bs[:si+1]
	}

	return bs
}

// GetBinaryStrings get binary-strings of bytes
func GetBinaryStrings(bts []byte) []string {
	var bs []string
	for _, b := range bts {
		bs = append(bs, fmt.Sprintf("%08b", b))
	}
	return bs
}

func GetIntBinString(n int64, be, isBit64 bool, sep string) string {
	bts := GetIntBytes(n, be, isBit64)
	bs := GetBinaryStrings(bts)
	return strings.Join(bs, sep)
}

// FindFirstByteIndex find first byte index satisfy func condition, return -1 if not found
func FindFirstByteIndex(bs []byte, leftToRight bool, fn func(b byte) bool) int {
	if leftToRight {
		for i, b := range bs {
			if fn(b) {
				return i
			}
		}
		return -1
	}

	for i := len(bs) - 1; i > 0; i-- {
		b := bs[i]
		if fn(b) {
			return i
		}
	}
	return -1
}
