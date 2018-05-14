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
// @File: receipt.go
// @Date: 2018/05/08 15:18:08
////////////////////////////////////////////////////////////////////////////////

package transaction

import (
	"bytes"
	"fmt"
	"mjoy.io/common/types"
	"mjoy.io/common/types/util"
	"mjoy.io/utils/bloom"
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

var (
	receiptStatusFailedMsgp     = []byte{}
	receiptStatusSuccessfulMsgp = []byte{0x01}
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint(1)
)

// Receipt represents the results of a transaction.
type Receipt struct {
	// Consensus fields
	PostState         []byte      `json:"root"`
	Status            uint        `json:"status"`
	Bloom             types.Bloom `json:"logsBloom"         gencodec:"required"`
	Logs              []*Log      `json:"logs"              gencodec:"required"`

	// Implementation fields (don't reorder!)
	TxHash          types.Hash    `json:"transactionHash"   gencodec:"required"`
	ContractAddress types.Address `json:"contractAddress"`
}




type ReceiptMsgp struct {
	PostStateOrStatus []byte
	Bloom             types.Bloom
	Logs              []*Log
}


// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewReceipt(root []byte, failed bool) *Receipt {
	r := &Receipt{PostState: util.CopyBytes(root)}
	if failed {
		r.Status = ReceiptStatusFailed
	} else {
		r.Status = ReceiptStatusSuccessful
	}
	return r
}


//when decoding,use it to set Status depend on postStateOrStatus
func (r *Receipt) SetStatus(postStateOrStatus []byte) error {
	switch {
	case bytes.Equal(postStateOrStatus, receiptStatusSuccessfulMsgp):
		r.Status = ReceiptStatusSuccessful
	case bytes.Equal(postStateOrStatus, receiptStatusFailedMsgp):
		r.Status = ReceiptStatusFailed
	case len(postStateOrStatus) == len(types.Hash{}):
		r.PostState = postStateOrStatus
	default:
		return fmt.Errorf("invalid receipt status %x", postStateOrStatus)
	}
	return nil
}
//when encoding,use it to set PostStateOrStatus
func (r *Receipt) StatusEncoding() []byte {
	if len(r.PostState) == 0 {
		if r.Status == ReceiptStatusFailed {
			return receiptStatusFailedMsgp
		}
		return receiptStatusSuccessfulMsgp
	}
	return r.PostState
}

// String implements the Stringer interface.
func (r *Receipt) String() string {
	if len(r.PostState) == 0 {
		return fmt.Sprintf("receipt{status=%d   bloom=%x logs=%v}", r.Status,  r.Bloom, r.Logs)
	}
	return fmt.Sprintf("receipt{med=%x   bloom=%x logs=%v}", r.PostState,   r.Bloom, r.Logs)
}

// ReceiptForStorage is a wrapper around a Receipt that flattens and parses the
// entire content of a receipt, as opposed to only the consensus fields originally.
type ReceiptForStorage struct {
	PostState         []byte      `json:"root"`
	Status            uint        `json:"status"`
	Bloom             types.Bloom `json:"logsBloom"         gencodec:"required"`
	Logs              []*Log      `json:"logs"              gencodec:"required"`

	// Implementation fields (don't reorder!)
	TxHash          types.Hash    `json:"transactionHash"   gencodec:"required"`
	ContractAddress types.Address `json:"contractAddress"`
}



// Receipts is a wrapper around a Receipt array to implement DerivableList.
type Receipts []*Receipt

// Len returns the number of receipts in this list.
func (r Receipts) Len() int { return len(r) }

func (r Receipts)GetMsgp(i int)[]byte{
	var buf bytes.Buffer
	err := msgp.Encode(&buf, r[i])
	if err != nil{
		return nil
	}
	return buf.Bytes()
}

func CreateBloom(receipts Receipts) types.Bloom {
	bloomIn := []bloom.BloomByte{}
	for _, receipt := range receipts {
		for _, log := range receipt.Logs {
			bloomIn = append(bloomIn, log.Address)
			for _, b := range log.Topics {
				bloomIn = append(bloomIn, b)
			}
		}
	}
	return bloom.CreateBloom(bloomIn)
}

//type Receipts_s [][]*Receipt
type Receipts_s struct{
	Receipts_s [] Receipts
}