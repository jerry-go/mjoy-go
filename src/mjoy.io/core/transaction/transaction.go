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
// @File: transaction.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package transaction

import (
	"errors"
	"math/big"
	"sync/atomic"

	"mjoy.io/utils/crypto"
	"mjoy.io/common/types/util"
	"mjoy.io/common"
	"fmt"
	"mjoy.io/utils/crypto/sha3"
	"mjoy.io/common/types"
	"bytes"
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp
//msgp:ignore Message TransactionsByPriceAndNonce

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
	errNoSigner   = errors.New("missing signing methods")
)

// deriveSigner makes a *best* guess about which signer to use.
func deriveSigner(V *big.Int) Signer {
	return NewMSigner(deriveChainId(V))
}

type Transaction struct {
	Data Txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

func (this * Transaction)PrintDataInfo(){
	logger.Debug("Txdata.Amount Value:",this.Data.Amount)
	//fmt.Println("txData.Amount Value:",this.Data.Amount)
}

//for test
func (this * Transaction)PrintVSR(){
	fmt.Printf("V:=%d,S:=%d,R:=%d\n",this.Data.V.IntVal.Int64(),
											this.Data.S.IntVal.Int64(),
											this.Data.R.IntVal.Int64())
}

func (this *Transaction)msgpHash()(h types.Hash){
	var buf bytes.Buffer
	err := msgp.Encode(&buf, this)
	if err != nil{
		return types.Hash{}
	}else{
		hw:=sha3.NewKeccak256()
		hw.Write(buf.Bytes())
		hw.Sum(h[:0])

		return h
	}
}


type Txdata struct {
	AccountNonce uint64         `json:"nonce"    gencodec:"required"`
	Recipient    *types.Address `json:"to"       msgp:"nil"` // nil means contract creation
	Amount       *types.BigInt  `json:"value"    gencodec:"required"`
	Payload      []byte         `json:"input"    gencodec:"required"`

	// Signature values
	V *types.BigInt `json:"v" gencodec:"required"`
	R *types.BigInt `json:"r" gencodec:"required"`
	S *types.BigInt `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *types.Hash `json:"hash" msgp:"-"`
}


func NewTransaction(nonce uint64, to types.Address, amount *big.Int, no1 uint64, no2 *big.Int, data []byte) *Transaction {
	return newTransaction(nonce, &to, amount, no1, no2, data)
}

func NewContractCreation(nonce uint64, amount *big.Int, no1 uint64, no2 *big.Int, data []byte) *Transaction {
	return newTransaction(nonce, nil, amount, no1, no2, data)
}

func newTransaction(nonce uint64, to *types.Address, amount *big.Int, no1 uint64, no2 *big.Int, data []byte) *Transaction {
	if len(data) > 0 {
		data = util.CopyBytes(data)
	}
	d := Txdata{
		AccountNonce: nonce,
		Recipient:    to,
		Payload:      data,
		Amount:       new(types.BigInt),
		V:            new(types.BigInt),
		R:            new(types.BigInt),
		S:            new(types.BigInt),
	}
	if amount != nil {
		d.Amount.IntVal = *amount
	}

	return &Transaction{Data: d}
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainId() *big.Int {
	return deriveChainId(&tx.Data.V.IntVal)
}

// Protected returns whether the transaction is protected from replay protection.
func (tx *Transaction) Protected() bool {
	return isProtectedV(&tx.Data.V.IntVal)
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		//if v is 27 or 28,return false
		return v != 27 && v != 28
	}
	// anything not 27 or 28 are considered unprotected
	return true
}


// MarshalJSON encodes the web3 RPC transaction format.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	hash := tx.Hash()
	data := tx.Data
	data.Hash = &hash
	return data.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec Txdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}
	var V byte
	if isProtectedV(&dec.V.IntVal) {
		chainID := deriveChainId(&dec.V.IntVal).Uint64()
		V = byte(dec.V.IntVal.Uint64() - 35 - 2*chainID)
	} else {
		V = byte(dec.V.IntVal.Uint64() - 27)
	}
	if !crypto.ValidateSignatureValues(V, &dec.R.IntVal, &dec.S.IntVal, false) {
		return ErrInvalidSig
	}
	*tx = Transaction{Data: dec}
	return nil
}

func (tx *Transaction) DataPayload() []byte       { return util.CopyBytes(tx.Data.Payload) }
func (tx *Transaction) Value() *big.Int    {
	if tx.Data.Amount == nil{
		return nil
	}
	return new(big.Int).Set(&tx.Data.Amount.IntVal)
}
func (tx *Transaction) Nonce() uint64      { return tx.Data.AccountNonce }
func (tx *Transaction) CheckNonce() bool   { return true }

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (tx *Transaction) To() *types.Address {
	if tx.Data.Recipient == nil {
		return nil
	}
	to := *tx.Data.Recipient
	return &to
}

// Hash hashes the Msgp encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() types.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(types.Hash)
	}

	v := tx.msgpHash()
	tx.hash.Store(v)
	return v
}



type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}


// Size returns the true MSGP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	var buf bytes.Buffer
	err := msgp.Encode(&buf, tx)
	if err != nil{
		c = writeCounter(0)
	}else {
		c = writeCounter(len(buf.Bytes()))
	}

	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (tx *Transaction) AsMessage(s Signer) (Message, error) {
	msg := Message{
		nonce:      tx.Data.AccountNonce,
		to:         tx.Data.Recipient,
		amount:     &tx.Data.Amount.IntVal,
		data:       tx.Data.Payload,
		checkNonce: true,
	}

	var err error
	msg.from, err = Sender(s, tx)
	return msg, err
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{Data: tx.Data}
	cpy.Data.R, cpy.Data.S, cpy.Data.V = &types.BigInt{*r}, &types.BigInt{*s}, &types.BigInt{*v}
	return cpy, nil
}


// Cost returns amount.
func (tx *Transaction) Cost() *big.Int {
	total := big.NewInt(0)
	total.Add(total, &tx.Data.Amount.IntVal)
	return total
}

func (tx *Transaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return &tx.Data.V.IntVal, &tx.Data.R.IntVal, &tx.Data.S.IntVal
}

func (tx *Transaction) String() string {
	var from, to string
	if tx.Data.V != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(&tx.Data.V.IntVal)
		if f, err := Sender(signer, tx); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if tx.Data.Recipient == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", tx.Data.Recipient[:])
	}
	var buf bytes.Buffer
	msgp.Encode(&buf, tx)
	return fmt.Sprintf(`
	TX(%x)
	Contract: %v
	From:     %s
	To:       %s
	Nonce:    %v
	Value:    %#x
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x
	Hex:      %x
`,
		tx.Hash(),
		tx.Data.Recipient == nil,
		from,
		to,
		tx.Data.AccountNonce,
		tx.Data.Amount,
		tx.Data.Payload,
		tx.Data.V,
		tx.Data.R,
		tx.Data.S,
		buf.Bytes(),
	)
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetMsgp implements Msgpable and returns the i'th element of s in msgp.
func (s Transactions)GetMsgp(i int)[]byte{
	var buf bytes.Buffer
	err := msgp.Encode(&buf, s[i])
	if err != nil{
		return nil
	}

	return buf.Bytes()
}

// TxDifference returns a new set t which is the difference between a to b.
func TxDifference(a, b Transactions) (keep Transactions) {
	keep = make(Transactions, 0, len(a))

	remove := make(map[types.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}
//for block producing
type TransactionForProducing struct {
	txs map[types.Address]Transactions	//all the transactions with address
	heads Transactions
	signer Signer
}

func NewTransactionsForProducing(signer Signer , txs map[types.Address]Transactions ) * TransactionForProducing{
	heads := new(Transactions)
	for _ , accTxs := range txs {
		*heads = append(*heads , accTxs[0])
		acc , _ := Sender(signer , accTxs[0])
		txs[acc] = accTxs[1:]
	}
	return &TransactionForProducing{
		txs:txs,
		heads:*heads,
		signer:signer,
	}
}

func (t *TransactionForProducing)Peek()*Transaction{
	if len(t.heads) == 0{
		return nil
	}
	return t.heads[0]
}

func (t *TransactionForProducing)Pop(){
	if len(t.heads) > 0 {
		t.heads = t.heads[1:]
	}
}

func (t *TransactionForProducing)Shift(){
	acc , _ := Sender(t.signer , t.heads[0])
	if txs , ok := t.txs[acc];ok && len(txs) > 0{
		t.heads[0],t.txs[acc] = txs[0] , txs[1:]
	}else{
		t.Pop()
	}
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int           { return len(s) }
func (s TxByNonce) Less(i, j int) bool { return s[i].Data.AccountNonce < s[j].Data.AccountNonce }
func (s TxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }


// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	to         *types.Address
	from       types.Address
	nonce      uint64
	amount     *big.Int
	data       []byte
	checkNonce bool
}

func NewMessage(from types.Address, to *types.Address, nonce uint64, amount *big.Int, data []byte, checkNonce bool) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		data:       data,
		checkNonce: checkNonce,
	}
}

func (m Message) From() types.Address { return m.from }
func (m Message) To() *types.Address  { return m.to }
func (m Message) Value() *big.Int      { return m.amount }
func (m Message) Nonce() uint64        { return m.nonce }
func (m Message) Data() []byte         { return m.data }
func (m Message) CheckNonce() bool     { return m.checkNonce }
