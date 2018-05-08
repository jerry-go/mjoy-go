package types

import (
	"fmt"
	"math/big"
	"mjoy.io/common/types/util/hex"
)

//go:generate msgp

const (
	// BloomByteLength represents the number of bytes used in a header log bloom.
	BloomByteLength = 256

	// BloomBitLength represents the number of bits used in a header log bloom.
	BloomBitLength = 8 * BloomByteLength
)

var (
	bloomType int8
)

// Bloom represents a 2048 bit bloom filter.
type Bloom [BloomByteLength]byte

// BytesToBloom converts a byte slice to a bloom filter.
// It panics if b is not of suitable size.
func BytesToBloom(b []byte) Bloom {
	var bloom Bloom
	bloom.SetBytes(b)
	return bloom
}

// SetBytes sets the content of b to the given bytes.
// It panics if d is not of suitable size.
func (b *Bloom) SetBytes(d []byte) {
	if len(b) < len(d) {
		panic(fmt.Sprintf("bloom bytes too big %d %d", len(b), len(d)))
	}
	copy(b[BloomByteLength-len(d):], d)
}

// Big converts b to a big integer.
func (b Bloom) Big() *big.Int {
	return new(big.Int).SetBytes(b[:])
}

func (b Bloom) Bytes() []byte {
	return b[:]
}

// MarshalText encodes b as a hex string with 0x prefix.
func (b Bloom) MarshalText() ([]byte, error) {
	return hex.Bytes(b[:]).MarshalText()
}

// UnmarshalText b as a hex string with 0x prefix.
func (b *Bloom) UnmarshalText(input []byte) error {
	return hex.UnmarshalFixedText("Bloom", input, b[:])
}

// Here, we'll pick an arbitrary number between
// 0 and 127 that isn't already in use
func (bloom *Bloom) ExtensionType() int8 {
	return bloomType
}

// We'll always use 256 bytes to encode the data
func (bloom *Bloom) Len() int {
	return BloomByteLength
}

func (bloom *Bloom) MarshalBinaryTo(b []byte) error {
	copy(b[:], bloom[:])
	return nil
}

func (bloom *Bloom) UnmarshalBinary(b []byte) error {
	if BloomByteLength < len(b) {
		return ErrBytesTooLong
	}
	copy(bloom[BloomByteLength-len(bloom):], b)
	return nil
}

