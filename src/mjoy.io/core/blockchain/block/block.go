////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 The mjoy-go Authors.
//
// The mjoy-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// @File: block.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package block

import (
	"mjoy.io/core/transaction"
	"mjoy.io/common/types"
	"mjoy.io/trie"
	"encoding/binary"
	"bytes"
	"sync/atomic"
	"math/big"
	"time"
	"mjoy.io/utils/crypto/sha3"
	"sort"
	"mjoy.io/common/types/util"
	"mjoy.io/common"
	"github.com/tinylib/msgp/msgp"
	"fmt"
	"errors"
	"mjoy.io/utils/bloom"
)

type DerivableList interface {
	Len() int
	GetMsgp(i int) []byte
}

func DeriveSha(list DerivableList) types.Hash {
	keyBytesBuf := bytes.NewBuffer([]byte{})
	trie := new(trie.Trie)
	for i := 0; i < list.Len(); i++ {
		keyBytesBuf.Reset()
		binary.Write(keyBytesBuf, binary.BigEndian, i)
		trie.Update(keyBytesBuf.Bytes(), list.GetMsgp(i))
	}
	return trie.Hash()
}

//go:generate msgp
// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.

//msgp:shim BlockNonce as:interface{} using:toBytes/fromBytes mode:convert
type BlockNonce [8]byte
func toBytes(v BlockNonce) (interface{}, error) {
	b := make([]byte, 8)
	copy(b, v[:])
	return b,nil
}
func fromBytes(s interface{}) (BlockNonce, error) {
	var out [8]byte
	if s == nil {
		return out,errors.New("input nil")
	}
	buf, ok := s.([]byte)
	if !ok {
		return out,errors.New("input type is not []byte")
	}
	if len(buf) != 8{
		return out,errors.New("input []byte len error")
	}
	copy(out[:],buf)
	return out,nil
}



// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// block header



type Header struct {
	ParentHash  types.Hash     		  `json:"parentHash" `
	Coinbase   	types.Address  		  `json:"blockProducer" `
	StateHash   types.Hash            `json:"stateRoot" `
	TxHash      types.Hash            `json:"transactionsRoot" `
	ReceiptHash types.Hash            `json:"receiptsRoot" `
	Bloom       types.Bloom           `json:"logsBloom" `
	Number      *types.BigInt         `json:"number" `
	Time        *types.BigInt         `json:"timestamp" `
	Extra       []byte                `json:"extraData" `
	MixHash     types.Hash            `json:"mixHash" `
	Nonce       BlockNonce            `json:"nonce" `
}

type Body struct {
	Transactions []*transaction.Transaction
}
type Block struct {
	B_header      *Header               // block header
	B_body        Body                  // all transactions in this block

	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package mjoy to track
	// inter-peer block relay.
	ReceivedAt   time.Time      `msg:"-"`
	ReceivedFrom interface{}    `msg:"-"`
}

func (h *Header) Hash() (out types.Hash) {
	var buf bytes.Buffer
	err := msgp.Encode(&buf, h)
	if err == nil{
		return sh3Hash(buf.Bytes())
	}else{
		return types.Hash{}
	}
}

func sh3Hash(x interface{}) (h types.Hash) {
	h3 := sha3.NewKeccak256()
	h3.Write(x.([]byte))
	h3.Sum(h[:0])
	return h
}


var (
	EmptyRootHash  = DeriveSha(transaction.Transactions{})
)
// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs and receipts.
func NewBlock(header *Header, txs []*transaction.Transaction, receipts []*transaction.Receipt) *Block {
	b := &Block{B_header: CopyHeader(header)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.B_header.TxHash = EmptyRootHash
	} else {
		b.B_header.TxHash = DeriveSha(transaction.Transactions(txs))
		b.B_body.Transactions = make(transaction.Transactions, len(txs))
		copy(b.B_body.Transactions, txs)
	}

	if len(receipts) == 0 {
		b.B_header.ReceiptHash = EmptyRootHash
	} else {
		b.B_header.ReceiptHash = DeriveSha(transaction.Receipts(receipts))
		bloomIn := []bloom.BloomByte{}
		for _, receipt := range receipts {
			for _, log := range receipt.Logs {
				bloomIn = append(bloomIn, log.Address)
				for _, b := range log.Topics {
					bloomIn = append(bloomIn, b)
				}
			}
		}
		b.B_header.Bloom = bloom.CreateBloom(bloomIn)
	}
	return b
}
// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Time = new(types.BigInt); h.Time != nil {
		cpy.Time.Put(h.Time.IntVal)
	}

	if cpy.Number = new(types.BigInt); h.Number != nil {
		cpy.Number.Put(h.Number.IntVal)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	return &cpy
}

func NewBlockWithHeader(header *Header) *Block {
	return &Block{B_header: CopyHeader(header)}
}

func (header *Header) HashNoNonce() types.Hash {
	v := &HeaderNoNonce{
		header.ParentHash,
		header.Coinbase,
		header.StateHash,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Number,
		header.Time,
		header.Extra,
	}
	return v.Hash()
}

func (b *Block) Transactions() transaction.Transactions { return b.B_body.Transactions }

func (b *Block) Transaction(hash types.Hash) *transaction.Transaction {
	for _, transaction := range b.B_body.Transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}


func (b *Block) Number() *big.Int     { return new(big.Int).Set(&b.B_header.Number.IntVal) }
func (b *Block) Time() *big.Int       { return new(big.Int).Set(&b.B_header.Time.IntVal) }

func (b *Block) NumberU64() uint64       { return b.B_header.Number.IntVal.Uint64() }
func (b *Block) MixDigest() types.Hash   { return b.B_header.MixHash }
func (b *Block) Nonce() uint64           { return binary.BigEndian.Uint64(b.B_header.Nonce[:]) }
func (b *Block) Bloom() types.Bloom      { return b.B_header.Bloom }
func (b *Block) Coinbase() types.Address { return b.B_header.Coinbase }
func (b *Block) Root() types.Hash        { return b.B_header.StateHash }
func (b *Block) ParentHash() types.Hash  { return b.B_header.ParentHash }
func (b *Block) TxHash() types.Hash      { return b.B_header.TxHash }
func (b *Block) ReceiptHash() types.Hash { return b.B_header.ReceiptHash }
func (b *Block) Extra() []byte           { return util.CopyBytes(b.B_header.Extra) }

func (b *Block) Header() *Header { return CopyHeader(b.B_header) }

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.B_body.Transactions} }


func (b *Block) HashNoNonce() types.Hash {
	return b.B_header.HashNoNonce()
}

func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := 0
	var buf bytes.Buffer
	msgp.Encode(&buf, b)
	c = buf.Len()
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}


// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header

	return &Block{
		B_header:       &cpy,
		B_body:         b.B_body,
	}
}

// WithBody returns a new block with the given transaction  contents.
func (b *Block) WithBody(body *Body) *Block {
	block := &Block{
		B_header:       CopyHeader(b.B_header),
	}
	block.B_body.Transactions = make([]*transaction.Transaction, len(body.Transactions))
	copy(block.B_body.Transactions, body.Transactions)
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() types.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(types.Hash)
	}
	v := b.B_header.Hash()
	b.hash.Store(v)
	return v
}


func (b *Block) String() string {
	str := fmt.Sprintf(`Block(#%v): Size: %v {
BlockproducerHash: %x
%v
Transactions:
%v
}
`, b.Number(), b.Size(), b.B_header.HashNoNonce(), b.B_header, b.B_body.Transactions)
	return str
}

func (h *Header) String() string {
	return fmt.Sprintf(`Header(%x):
[
	ParentHash:	    %x
	Coinbase:	    %x
	Root:		    %x
	TxSha		    %x
	ReceiptSha:	    %x
	Bloom:		    %x
	Number:		    %v
	Time:		    %v
	Extra:		    %s
	MixDigest:      %x
	Nonce:		    %x
]`, h.Hash(), h.ParentHash, h.Coinbase, h.StateHash, h.TxHash, h.ReceiptHash, h.Bloom, h.Number, h.Time, h.Extra, h.MixHash, h.Nonce)
}


//blocks part

type Blocks []*Block

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.B_header.Number.IntVal.Cmp(&b2.B_header.Number.IntVal) < 0 }



// header wihtout mixhash and nonce

type HeaderNoNonce struct {
	ParentHash  types.Hash
	Coinbase   types.Address
	StateHash   types.Hash
	TxHash      types.Hash
	ReceiptHash types.Hash
	Bloom       types.Bloom
	Number      *types.BigInt
	Time        *types.BigInt
	Extra       []byte
}
func (h *HeaderNoNonce) Hash() types.Hash {
	var buf bytes.Buffer
	err := msgp.Encode(&buf, h)
	if err == nil{
		return sh3Hash(buf.Bytes())
	}else{
		return types.Hash{}
	}
}

//type Headers []*Header
type Headers struct {
	Headers []*Header
}
