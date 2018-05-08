package types

import (
	"fmt"
	"mjoy.io/common/types/util"
	"mjoy.io/common/types/util/hex"
	"math/big"
	"reflect"
)

//go:generate msgp

const (
	HashLength = 32
)

var (
	hashType int8
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data. It's a hex string/num
type Hash [HashLength]byte

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}
func StringToHash(s string) Hash { return BytesToHash([]byte(s)) }
func BigToHash(b *big.Int) Hash  { return BytesToHash(b.Bytes()) }
func HexToHash(s string) Hash    { return BytesToHash(util.FromHex(s)) }

// Get the string representation of the underlying hash
func (h Hash) Str() string   { return string(h[:]) }
func (h Hash) Bytes() []byte { return h[:] }
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }
func (h Hash) Hex() string   { return hex.Encode(h[:]) }

// Here, we'll pick an arbitrary number between
// 0 and 127 that isn't already in use
func (*Hash) ExtensionType() int8 {
	return hashType
}

// We'll always use 32 bytes to encode the data
func (*Hash) Len() int {
	return HashLength
}

// MarshalBinaryTo simply copies the value
// of the bytes into 'b'
func (h *Hash) MarshalBinaryTo(b []byte) error {
	copy(b, h.Bytes())
	return nil
}

func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

func (h Hash) String() string {
	return h.Hex()
}


// UnmarshalBinary copies the value of 'b'
// into the Hash object. (We might want to add
// a sanity check here later that len(b) <= HashLength.)
func (h *Hash) UnmarshalBinary(b []byte) error {
	// TODO: check b, only hex, len <= HashLength
	if len(b) <= HashLength {
		*h = BytesToHash(b)
		return nil
	}

	return ErrBytesTooLong
}

// for json marshal
func (h Hash) MarshalText() ([]byte, error) {
	// TODO:
	return hex.Bytes(h[:]).MarshalText()
}

// for json unmarshal
func (h *Hash) UnmarshalJSON(b []byte) error {
	return hex.UnmarshalFixedJSON(reflect.TypeOf(Hash{}), b, h[:])
}

// for json unmarshal
func (h *Hash) UnmarshalText(b []byte) error {
	// TODO:
	return hex.UnmarshalFixedText("Hash", b, h[:])
}

// for format print
func (h Hash) Format(s fmt.State, c rune) {
	switch c {
	case 'x' | 'X':
		fmt.Fprintf(s, "%#x", h[:])
	default:
		fmt.Fprintf(s, "%"+string(c), h[:])
	}
}

// Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}
//type Hashs []*Hash

type Hashs struct {
	Hashs []*Hash
}
