package util

import (
	"fmt"
	"strconv"
	"strings"
)

func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

// Reversing bits in a word, refined basic scheme.
// 24 operations + 6 for loading constants = 30 insns.
// This is Figure 7-1 in HD.
func ReverseBits(x uint32) uint32 {
	x = (x&0x55555555)<<1 | (x>>1)&0x55555555
	x = (x&0x33333333)<<2 | (x>>2)&0x33333333
	x = (x&0x0F0F0F0F)<<4 | (x>>4)&0x0F0F0F0F
	x = (x << 24) | ((x & 0xFF00) << 8) |
		((x >> 8) & 0xFF00) | (x >> 24)
	return x
}
func ReverseUint64Bits(x uint64) uint64 {
	l := uint32(x >> 32)
	r := uint32(x & 0xFFFFFFFF)

	l = ReverseUint32Bits(l)
	r = ReverseUint32Bits(r)
	return uint64(l) | uint64(r)<<32
}

func ReverseUint32Bits(x uint32) uint32 {
	x = (x&0x55555555)<<1 | (x>>1)&0x55555555
	x = (x&0x33333333)<<2 | (x>>2)&0x33333333
	x = (x&0x0F0F0F0F)<<4 | (x>>4)&0x0F0F0F0F
	x = (x << 24) | ((x & 0xFF00) << 8) |
		((x >> 8) & 0xFF00) | (x >> 24)
	return x
}

//GetBitmap returns a bool array of whether bits are set or not set.
func GetBitmap(data []byte) []bool {
	boolMap := make([]bool, len(data)*8)

	for i, b := range data {
		bint := uint(b)
		for n := 0; n < 8; n++ {
			boolMap[(i*8)+n] = bint&(1<<uint(7-n)) > 0
		}
	}
	return boolMap
}

func Spacify(str string) string {
	var newStr string
	for i, s := range str {
		if i > 0 && i%2 == 0 {
			newStr += " "
		}
		newStr += string(s)
	}
	return newStr
}

func Split2(s, sep string) (string, string) {
	if len(s) == 0 {
		return s, s
	}

	array := strings.SplitN(s, sep, 2)

	// Incase no separator were present
	if len(array) == 1 {
		return array[0], ""
	}

	return array[0], array[1]
}

func CheckDigit(number string, base int) (digit int) {
	digits := make([]int64, (len(number)))
	var err error
	for i := 0; i < len(number); i++ {
		digits[i], err = strconv.ParseInt(string(number[i]), base, 64)
		if err != nil {
			fmt.Println(err)
		}
	}
	sum := int64(0)
	for i := len(digits) - 1; i >= 0; i-- {
		sum += 2*digits[i]/int64(base) + 2*digits[i]%int64(base)
		i--
		if i >= 0 {
			sum += digits[i]
		}
	}
	digit = base - int(sum%int64(base))
	if base == digit {
		digit = 0
	}
	return
}

func Checksummed(number string, base int) (check string) {
	return fmt.Sprintf("%s%d", number, CheckDigit(number, base))
}
