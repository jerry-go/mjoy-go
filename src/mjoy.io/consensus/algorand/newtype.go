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
// @File: types.go
// @Date: 2018/06/12 11:01:51
////////////////////////////////////////////////////////////////////////////////

package algorand

import (
	"mjoy.io/common/types"
	"mjoy.io/common"
	"mjoy.io/core/blockchain/block"
	"bytes"
	"github.com/tinylib/msgp/msgp"
)

// for algorand1, fill block header ConsensusData filed
// id = "algorand1"
// para = Q(r) = hash(Sig(Q(r-1)), r) where r = block number
var (
	ConsensusDataId = "algorand1"
)
//go:generate msgp
// credentialData is the data for generating credential
// credential = sig(credentialData)
type CredentialData struct {
	Round         types.BigInt
	Step          types.BigInt
	Quantity      types.Hash    //the seed of round r
}
func (s *CredentialData)GetMsgp()[]byte{
	var buf bytes.Buffer
	err := msgp.Encode(&buf, s)
	if err != nil{
		return nil
	}

	return buf.Bytes()
}

type CredentialSig struct {
	Round         types.BigInt
	Step          types.BigInt
	R             types.BigInt
	S             types.BigInt
	V             types.BigInt
}
func (s *CredentialSig)GetMsgp()[]byte{
	var buf bytes.Buffer
	err := msgp.Encode(&buf, s)
	if err != nil{
		return nil
	}

	return buf.Bytes()
}

type SignatureVal struct {
	R             types.BigInt
	S             types.BigInt
	V             types.BigInt
}
func (s *SignatureVal)GetMsgp()[]byte{
	var buf bytes.Buffer
	err := msgp.Encode(&buf, s)
	if err != nil{
		return nil
	}

	return buf.Bytes()
}

func (c *CredentialData) Hash() types.Hash {
	hash, err := common.MsgpHash(c)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

// step1 (Block Proposal) message
// m(r,1) = (Br, esig(H(Br)), σr1)
type M1 struct {
	Block         *block.Block
	Esig          []byte
	Credential    *CredentialSig
}

// step2 (The First Step of the Graded Consensus Protocol GC) message
// step3 (The Second Step of GC) message
// step2 and step3 message has the same structure
// m(r,2) = (ESIG(v′), σr2)
type M23 struct {
	//hash is v′, the hash of the next block
	Hash          types.Hash
	Esig          []byte
	Credential    *CredentialSig
}

// step4 and step other message
// m(r,s) = (ESIG(b), ESIG(v′), σrs)
type MCommon struct {
	//B is the BBA⋆ input b, 0 or 1
	B             uint
	EsigB         []byte
	//hash is v′, the hash of the next block
	Hash          types.Hash
	EsigV         []byte
	Credential    *CredentialSig
}

