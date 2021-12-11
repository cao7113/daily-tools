package lib

import (
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/suite"
	"math/bits"
	"strconv"
	"testing"
)

func (s *FunSuite) TestSplitToGroupChunks() {
	s.Equal([]string{"12"}, SplitToGroupChunks("12", 3, 0))
	s.Equal([]string{"1", "234"}, SplitToGroupChunks("1234", 3, 0))

	s.Equal([]string{"012"}, SplitToGroupChunks("12", 3, '0'))
}

func (s *FunSuite) TestRune() {
	c := '\u65e5'
	s.Equal('\U000065e5', c)
	s.Equal(rune(26085), c)
	s.Equal("日", string(c))
	s.Equal("严", "\u4e25")
	s.Equal(rune(10), '\n')
	s.Equal(rune(13), '\r')
}

func (s *FunSuite) TestParse() {
	n, err := strconv.ParseInt("-11", 10, 64)
	s.Nil(err)
	s.Equal(int64(-11), n)
}

func (s *FunSuite) TestXOR() {
	s.Equal(0, 123^123)
	s.Equal(123, 123^0)

	s.Run("swap var in-place", func() {
		x := 123
		y := 234
		x0 := x
		y0 := y
		x = x ^ y
		y = x ^ y
		x = x ^ y
		s.Equal(x0, y)
		s.Equal(y0, x)
	})

	s.Run("find two diff num", func() {
		nums := []int{2, 3, 4, 5, 7, 8, 9, 10, 12}
		x, y := findDiff2(nums, 2, 11)
		s.Equal(6+11, x+y)
	})

	s.Run("n-th bit flip", func() {
		n := 12356
		fn := 10
		c := 1 << fn
		nf := n ^ c
		b0 := n & c
		b1 := nf & c

		fmt.Printf("bin: %s flip-num %d \n", GetIntBinString(int64(c), true, false, " "), c)
		fmt.Printf("bin: %s raw num: %d\n", GetIntBinString(int64(n), true, false, " "), n)
		fmt.Printf("bin: %s flipped num: %d\n", GetIntBinString(int64(nf), true, false, " "), nf)

		if b0 == 0 {
			s.Equal(c, b1)
		} else {
			s.Equal(c, b0)
			s.Equal(0, b1)
		}
	})

	s.Run("flip count", func() {
		a := 0b0000_1011
		b := 0b0010_0010
		c := a ^ b
		s.Equal(3, bits.OnesCount(uint(c)))
	})
}

func (s *FunSuite) TestBitwiseOps() {
	s.Run("right shift as divide 2", func() {
		i := -1      // 0b1111_1110
		i1 := i >> 1 // 0b111_1111
		s.EqualValues(-1, i1)
		s.EqualValues(0xff, uint8(i1))

		j := int8(2) // 0b0000_0010
		j1 := j >> 1 // 0b0000_0001
		s.EqualValues(1, j1)
		s.EqualValues(0x01, uint8(j1))
	})

	s.Run("left shift as multiply 2", func() {
		i := int16(-66) // 0b1111_1111_1011_1110
		i1 := i << 1    // 0b1111_1111_0111_1100
		s.Equal(int16(-132), i1)
		s.Equal(uint8(0b0111_1100), uint8(i1))

		j := int8(2) // 0b0000_0010
		j1 := j << 1 // 0b0000_0100
		s.EqualValues(4, j1)
		s.EqualValues(0b0000_0100, uint8(j1))
	})

	s.Run("two's component", func() {
		ui := binary.BigEndian.Uint16([]byte{0xff, 0xfe})
		fmt.Printf("%016b\n", ui) // 11111111 11111110
		i := int16(ui)
		s.EqualValues(-2, i)

		j := -1
		s.EqualValues(uint8(0xff), uint8(j)) // 0xff --> 1111 1111
	})

	s.Run("ones count", func() {
		x := 0b0110_1100_0001_0101
		wCnt := bits.OnesCount(uint(x))
		cnt := onesCount(x)
		s.Equal(wCnt, cnt)
	})

	s.Run("total hamming distances", func() {
		type hArg struct {
			nums []int32
			hSum int
		}
		args := []*hArg{
			{
				nums: []int32{4, 12, 28},
				hSum: 4,
			},
			{
				nums: []int32{4, 2, 14},
				hSum: 6,
			},
		}

		for _, a := range args {
			hs := sumHammingDistances(a.nums)
			s.Equal(a.hSum, hs)
			for _, n := range args[0].nums {
				fmt.Printf("%08b num:	%d\n", n, n)
			}
			fmt.Printf("sum hamming-distances of %v is: %d\n", a.hSum, hs)
		}
	})
}

func sumHammingDistances(nums []int32) int {
	hs := 0
	for i := 0; i < 32; i++ {
		c0 := 0
		c1 := 0
		for _, n := range nums {
			if (n>>i)&1 == 1 {
				c1++
			} else {
				c0++
			}
		}
		hs += c1 * c0
	}
	return hs
}

func onesCount(x int) int {
	cnt := 0
	for x != 0 {
		x &= x - 1
		cnt++
	}
	return cnt
}

func findDiff2(nums []int, start, count int) (int, int) {
	xorResult := 0
	ns := make([]int, count)
	for i := 0; i < count; i++ {
		ns = append(ns, start+i)
	}

	for _, n := range nums {
		xorResult ^= n
	}
	for _, n := range ns {
		xorResult ^= n
	}
	b1 := xorResult & (-xorResult)

	nx := 0
	for _, n := range nums {
		if n&b1 == 0 {
			nx ^= n
		}
	}
	for _, n := range ns {
		if n&b1 == 0 {
			nx ^= n
		}
	}

	ny := xorResult ^ nx
	return nx, ny
}

func (s *FunSuite) TestFirst1() {
	x := 100
	f := x & (-x)
	s.Equal(0x04, f)
}

func (s *FunSuite) TestGetIntBytes() {
	type tArgs struct {
		num    int64
		be     bool
		is64   bool
		result []byte
	}

	demos := map[string]tArgs{
		"zero": {
			num:    0,
			be:     true,
			is64:   false,
			result: []byte{0x00},
		},
		// min-mode
		"positive below 128 number with big-endian and min-mode": {
			num:    2,
			be:     true,
			is64:   false,
			result: []byte{0x02},
		},
		"positive single-byte msb changed with big-endian and min-mode": {
			num:    255,
			be:     true,
			is64:   false,
			result: []byte{0x00, 0xff},
		},
		"negative single-byte with big-endian and min-mode": {
			num:    -2,
			be:     true,
			is64:   false,
			result: []byte{0xfe},
		},
		"negative single-byte msb changed with big-endian and min-mode": {
			num:    -129,
			be:     true,
			is64:   false,
			result: []byte{0xff, 0x7f},
		},
		"negative number with big-endian and min-mode": {
			num:    -257,
			be:     true,
			is64:   false,
			result: []byte{0xfe, 0xff},
		},
		"positive number with little-endian and min-mode": {
			num:    258,
			be:     false,
			is64:   false,
			result: []byte{0x02, 0x01},
		},
		"negative number with little-endian and min-mode": {
			num:    -257,
			be:     false,
			is64:   false,
			result: []byte{0xff, 0xfe},
		},

		// 64bit-mode
		"positive number with big-endian and 64-bit": {
			num:    1,
			be:     true,
			is64:   true,
			result: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x0, 0x00, 0x01},
		},
		"negative number with big-endian and 64-bit": {
			num:    -2,
			be:     true,
			is64:   true,
			result: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe},
		},
		"positive number with little-endian and 64-bit": {
			num:    1,
			be:     false,
			is64:   true,
			result: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0, 0x00},
		},
		"negative number with little-endian and 64-bit": {
			num:    -2,
			be:     false,
			is64:   true,
			result: []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
	}

	for n, a := range demos {
		s.Run(n, func() {
			bs := GetIntBytes(a.num, a.be, a.is64)
			s.EqualValues(a.result, bs)
		})
	}
}

func (s *FunSuite) TestFindFirstByteIndex() {
	s.Equal(2, FindFirstByteIndex([]byte{0x00, 0x00, 0x01}, true, func(b byte) bool {
		return b != 0x00
	}))
	s.Equal(0, FindFirstByteIndex([]byte{0x01, 0x00, 0x00}, true, func(b byte) bool {
		return b != 0x00
	}))
	s.Equal(1, FindFirstByteIndex([]byte{0xff, 0xfe, 0xff}, false, func(b byte) bool {
		return b != 0xff
	}))
}

func TestFunSuite(t *testing.T) {
	suite.Run(t, &FunSuite{})
}

type FunSuite struct {
	suite.Suite
}
