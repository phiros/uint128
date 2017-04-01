package uint128

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"math"
)

var (
	target = map[string]struct{ H, L uint64 }{
		"456":               {0x0, 0x456},
		"10000000000000456": {0x1, 0x456},
		"e0000000000000009": {0xe, 0x9},
		"10000000000000000": {0x1, 0x0},
	}
)

func TestEncode(t *testing.T) {

	for k, v := range target {
		u := new(Uint128)
		u.H = v.H
		u.L = v.L
		if u.HexString() != k {
			z := []rune(u.HexString())
			for i, r := range k {
				if r != z[i] {
					t.Error("missmatch expect\n", k, "\n", u.HexString(), "\n",
						strings.Repeat(" ", i), "^")
					break
				}
			}
		}
	}
}

func TestLoadFromByte(t *testing.T) {
	for k, v := range target {
		u := &Uint128{}
		b, _ := hex.DecodeString(fmt.Sprintf("%032s", k))
		err := binary.Read(bytes.NewReader(b),
			binary.BigEndian, u)

		if err != nil || u.H != v.H || u.L != v.L {
			i := new(big.Int)
			i.SetString(fmt.Sprintf("0x%032s", k), 0)
			t.Error("missmatch ", fmt.Sprintf("%032s", k), i.Text(16), u, err)
		}
	}
}

func TestXor(t *testing.T) {
	xor := []struct{ s, x, cmp string }{
		{"1", "1", "0"},
		{"2", "1", "3"},
		{"e0000000000000009", "f0000000000000000", "10000000000000009"},
	}
	for _, entry := range xor {
		s, err := NewFromString(entry.s)
		if err != nil {
			t.Error(err)
		}
		x, _ := NewFromString(entry.x)
		if err != nil {
			t.Error(err)
		}
		cmp, _ := NewFromString(entry.cmp)
		if err != nil {
			t.Error(err)
		}

		t.Log(s, x, cmp)

		s.Xor(x)
		if s.Compare(cmp) != 0 {
			t.Error("failed xor at", entry, s, x, cmp)
		}
	}
}

func TestUint128_BigEndianBytes(t *testing.T) {
	var expectedBytesMaxUint128 []byte = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						    0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	var expectedAscendingBytesPattern []byte = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
							  0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
	var expectedDescendingBytesPattern []byte = []byte{0xff, 0xfe, 0xfd, 0xfc, 0xfb, 0xfa, 0xf0, 0xf9,
							   0xf8, 0xf7, 0xf6, 0xf5, 0xf4, 0xf3, 0xf2, 0xf1}
	var uint128 *Uint128 = new(Uint128)
	uint128.H = math.MaxUint64
	uint128.L = math.MaxUint64

	var actualBytes = uint128.BigEndianBytes()
	if bytes.Compare(actualBytes, expectedBytesMaxUint128) != 0 {
		t.Error("expected maxUint128 but got something else")
	}

	uint128.H = 0x0001020304050607
	uint128.L = 0x0809101112131415
	actualBytes = uint128.BigEndianBytes()
	if bytes.Compare(actualBytes, expectedAscendingBytesPattern) != 0 {
		t.Error("expected ascending pattern but got something else", actualBytes,
			expectedAscendingBytesPattern)
	}

	uint128.H = 0xfffefdfcfbfaf0f9
	uint128.L = 0xf8f7f6f5f4f3f2f1
	actualBytes = uint128.BigEndianBytes()
	if bytes.Compare(actualBytes, expectedDescendingBytesPattern) != 0 {
		t.Error("expected ascending pattern but got something else", actualBytes, expectedDescendingBytesPattern)
	}
}

func TestNewFromBigEndianBytes(t *testing.T) {
	var maxUint128 []byte = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				       0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	var ascendingBytesPattern []byte = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
						  0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
	var descendingBytesPattern []byte = []byte{0xff, 0xfe, 0xfd, 0xfc, 0xfb, 0xfa, 0xf0, 0xf9,
						   0xf8, 0xf7, 0xf6, 0xf5, 0xf4, 0xf3, 0xf2, 0xf1}

	var uint128Bytes, err = NewFromBigEndianBytes(maxUint128)
	if err != nil {
		t.Error(err)
	}
	if uint128Bytes.H != math.MaxUint64 || uint128Bytes.L != math.MaxUint64 {
		t.Error("expected maxUint128 but got something else", uint128Bytes)
	}

	uint128Bytes, err = NewFromBigEndianBytes(ascendingBytesPattern)
	if err != nil {
		t.Error(err)
	}
	if uint128Bytes.H != 0x0001020304050607 || uint128Bytes.L != 0x0809101112131415 {
		t.Error("expected ascending bit patterm but got something else", uint128Bytes)
	}

	uint128Bytes, err = NewFromBigEndianBytes(descendingBytesPattern)
	if err != nil {
		t.Error(err)
	}
	if uint128Bytes.H != 0xfffefdfcfbfaf0f9 || uint128Bytes.L != 0xf8f7f6f5f4f3f2f1 {
		t.Error("expected descending bit patterm but got something else", uint128Bytes)
	}
}

func TestNewFromUint64(t *testing.T) {
	const ascendingPattern = 0x0809101112131415
	const descendingPattern = 0xf8f7f6f5f4f3f2f1

	var uint128 = NewFromUint64(math.MaxUint64)
	if uint128.H != 0 || uint128.L != math.MaxUint64 {
		t.Error("expected MaxUint64 but got something else")
	}

	uint128 = NewFromUint64(ascendingPattern)
	if uint128.H != 0 || uint128.L !=  ascendingPattern {
		t.Error("expected ascending bit patterm but got something else")
	}

	uint128 = NewFromUint64(descendingPattern)
	if uint128.H != 0 || uint128.L !=  descendingPattern {
		t.Error("expected descending bit patterm but got something else")
	}
}

func TestUint128_Add(t *testing.T) {
	var maxUint128Bytes []byte = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					    0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}


	var u128number = NewFromUint64(math.MaxUint64)
	u128number.Add(NewFromUint64(1))
	if u128number.H != 1 || u128number.L != 0 {
		t.Error("expected MaxUint64 + 1 but got", u128number.H, u128number.L)
	}

	u128numberMax, _ := NewFromBigEndianBytes(maxUint128Bytes)
	u128numberMax.Add(NewFromUint64(1))
	if u128numberMax.H != 0 || u128numberMax.L != 0 {
		t.Error("expected MaxUint128 + 1 = 0 but got", u128numberMax.H, u128numberMax.L)
	}

	u128numberMax1, _ := NewFromBigEndianBytes(maxUint128Bytes)
	u128numberMax2, _ := NewFromBigEndianBytes(maxUint128Bytes)
	u128numberMax1.Add(u128numberMax2)
	if u128numberMax1.H != 0xffffffffffffffff || u128numberMax1.L != 0xfffffffffffffffe {
		t.Error("expected MaxUint128 - 1 but got", u128numberMax1.H, u128numberMax1.L)
	}
}

func TestUint128_Sub(t *testing.T) {
	var maxUint128Bytes []byte = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					    0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}


	u128numberMax, _ := NewFromBigEndianBytes(maxUint128Bytes)
	u128numberMax.Sub(NewFromUint64(1))
	if u128numberMax.H != 0xffffffffffffffff || u128numberMax.L != 0xfffffffffffffffe {
		t.Error("expected MaxUint128 -1 but got", u128numberMax.H, u128numberMax.L)
	}

	u128numberZero := NewFromUint64(0)
	u128numberZero.Sub(NewFromUint64(1))
	if u128numberZero.H != 0xffffffffffffffff || u128numberZero.L != 0xffffffffffffffff {
		t.Error("expected MaxUint128 but got", u128numberZero.H, u128numberZero.L)
	}
}
