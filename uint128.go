package uint128

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
)

const (
	lessThan = iota - 1
	equal
	greaterThan

	Len = 32
	LenBytes = 16
)

var (
	MaxUint128 = Uint128{ math.MaxUint64, math.MaxUint64}
)

// Big endian uint128
type Uint128 struct {
	H, L uint64
}

func (u Uint128) Compare(o Uint128) int {
	if u.H < o.H {
		return lessThan
	} else if u.H > o.H {
		return greaterThan
	}

	if u.L < o.L {
		return lessThan
	} else if u.L > o.L {
		return greaterThan
	}

	return equal
}

func (u Uint128) And(o Uint128) (ans Uint128) {
	ans.H = u.H & o.H
	ans.L = u.L & o.L
	return
}

func (u Uint128) Or(o Uint128) (ans Uint128) {
	ans.H = u.H | o.H
	ans.L = u.L | o.L
	return
}

func (u Uint128) Xor(o Uint128) (ans Uint128) {
	ans.H = u.H ^ o.H
	ans.L = u.L ^ o.L
	return
}

// See: https://www.codeproject.com/Tips/617214/UInt-Addition-Subtraction
// for an explanation
func (u Uint128) Add(o Uint128) (ans Uint128) {
	var C uint64 = (((u.L & o.L) & 1) + (u.L >> 1) + (o.L >> 1)) >> 63
	ans.H = u.H + o.H + C
	ans.L = u.L + o.L
	return ans
}

// See: https://www.codeproject.com/Tips/617214/UInt-Addition-Subtraction
// for an explanation
func (u Uint128) Sub(o Uint128) (ans Uint128) {
	ans.L = u.L - o.L
	var C uint64 = (((ans.L & o.L) & 1) + (o.L >> 1) + (ans.L >> 1)) >> 63
	ans.H = u.H - (o.H + C)
	return
}

func NewFromUint64(uint64 uint64) (ans Uint128) {
	ans.L = uint64
	ans.H = 0
	return
}

func NewFromString(s string) (u Uint128, err error) {

	if len(s) > Len {
		return NewFromUint64(0), fmt.Errorf("s:%s length greater than 32", s)
	}

	b, err := hex.DecodeString(fmt.Sprintf("%032s", s))
	if err != nil {
		return NewFromUint64(0), err
	}
	rdr := bytes.NewReader(b)
	err = binary.Read(rdr, binary.BigEndian, &u)
	return
}

func NewFromBigEndianBytes(b []byte) (ans Uint128, err error) {
	if len(b) > LenBytes {
		return NewFromUint64(0), fmt.Errorf("length greater than 16 bytes")
	}
	rdr := bytes.NewReader(b)
	err = binary.Read(rdr, binary.BigEndian, &ans)
	return ans, err
}

func (u Uint128) BigEndianBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, u)
	return buf.Bytes()
}

func (u Uint128) HexString() string {
	if u.H == 0 {
		return fmt.Sprintf("%x", u.L)
	}
	return fmt.Sprintf("%x%016x", u.H, u.L)
}

func (u Uint128) String() string {
	return fmt.Sprintf("0x%032x", u.HexString())
}

func (u Uint128) BinaryString() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatUint(u.H, 2))
	buffer.WriteString(strconv.FormatUint(u.L, 2))
	return buffer.String()
}
