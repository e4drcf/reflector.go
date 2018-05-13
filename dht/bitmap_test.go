package dht

import (
	"testing"

	"github.com/lyoshenka/bencode"
)

func TestBitmap(t *testing.T) {
	a := Bitmap{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
		12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
		36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
	}
	b := Bitmap{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
		12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
		36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 47, 46,
	}
	c := Bitmap{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
	}

	if !a.Equals(a) {
		t.Error("bitmap does not equal itself")
	}
	if a.Equals(b) {
		t.Error("bitmap equals another bitmap with different id")
	}

	if !a.Xor(b).Equals(c) {
		t.Error(a.Xor(b))
	}

	if c.PrefixLen() != 375 {
		t.Error(c.PrefixLen())
	}

	if b.Less(a) {
		t.Error("bitmap fails lessThan test")
	}

	id := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if BitmapFromHexP(id).Hex() != id {
		t.Error(BitmapFromHexP(id).Hex())
	}
}

func TestBitmapMarshal(t *testing.T) {
	b := BitmapFromStringP("123456789012345678901234567890123456789012345678")
	encoded, err := bencode.EncodeBytes(b)
	if err != nil {
		t.Error(err)
	}

	if string(encoded) != "48:123456789012345678901234567890123456789012345678" {
		t.Error("encoding does not match expected")
	}
}

func TestBitmapMarshalEmbedded(t *testing.T) {
	e := struct {
		A string
		B Bitmap
		C int
	}{
		A: "1",
		B: BitmapFromStringP("222222222222222222222222222222222222222222222222"),
		C: 3,
	}

	encoded, err := bencode.EncodeBytes(e)
	if err != nil {
		t.Error(err)
	}

	if string(encoded) != "d1:A1:11:B48:2222222222222222222222222222222222222222222222221:Ci3ee" {
		t.Error("encoding does not match expected")
	}
}

func TestBitmapMarshalEmbedded2(t *testing.T) {
	encoded, err := bencode.EncodeBytes([]interface{}{
		BitmapFromStringP("333333333333333333333333333333333333333333333333"),
	})
	if err != nil {
		t.Error(err)
	}

	if string(encoded) != "l48:333333333333333333333333333333333333333333333333e" {
		t.Error("encoding does not match expected")
	}
}

func TestBitmap_PrefixLen(t *testing.T) {
	tt := []struct {
		hex string
		len int
	}{
		{len: 0, hex: "F00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{len: 0, hex: "800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{len: 1, hex: "700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{len: 1, hex: "400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{len: 384, hex: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{len: 383, hex: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"},
		{len: 382, hex: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002"},
		{len: 382, hex: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003"},
	}

	for _, test := range tt {
		len := BitmapFromHexP(test.hex).PrefixLen()
		if len != test.len {
			t.Errorf("got prefix len %d; expected %d for %s", len, test.len, test.hex)
		}
	}
}

func TestBitmap_ZeroPrefix(t *testing.T) {
	original := BitmapFromHexP("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	tt := []struct {
		zeros    int
		expected string
	}{
		{zeros: -123, expected: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{zeros: 0, expected: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{zeros: 1, expected: "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{zeros: 69, expected: "000000000000000007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{zeros: 383, expected: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"},
		{zeros: 384, expected: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{zeros: 400, expected: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, test := range tt {
		expected := BitmapFromHexP(test.expected)
		actual := original.ZeroPrefix(test.zeros)
		if !actual.Equals(expected) {
			t.Errorf("%d zeros: got %s; expected %s", test.zeros, actual.Hex(), expected.Hex())
		}
	}

	for i := 0; i < nodeIDLength*8; i++ {
		b := original.ZeroPrefix(i)
		if b.PrefixLen() != i {
			t.Errorf("got prefix len %d; expected %d for %s", b.PrefixLen(), i, b.Hex())
		}
	}
}
