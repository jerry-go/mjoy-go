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
// @File: message.go
// @Date: 2018/06/22 14:40:22
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"mjoy.io/core/blockchain/block"
	"mjoy.io/common/types"
	"mjoy.io/consensus/message"
	"fmt"
)
//go:generate msgp

type msgCredentialSig struct {
	cs    *CredentialSig
	*message.MsgPriv
}

func NewMsgCredentialSig(cs *CredentialSig) *msgCredentialSig{
	msgCs := &msgCredentialSig{
		cs:      cs,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgCs)
	return msgCs
}

func (tm *msgCredentialSig) DataHandle(data interface{}) {
	fmt.Println("msgBlockProposal data handle")
}

func (tm *msgCredentialSig) StopHandle() {
	fmt.Printf("stop ...\n")
}

// step1 (Block Proposal) message
// m(r,1) = (Br, esig(H(Br)), σr1)
type BlockProposal struct {
	Block         *block.Block
	Esig          []byte
	Credential    *CredentialSig
}
type msgBlockProposal struct {
	bp    *BlockProposal
	*message.MsgPriv
}

// new a message
func NewMsgBlockProposal(bp *BlockProposal) *msgBlockProposal{
	msgBp := &msgBlockProposal{
		bp:      bp,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgBp)
	return msgBp
}

func (tm *msgBlockProposal) DataHandle(data interface{}) {
	fmt.Println("msgBlockProposal data handle")
}

func (tm *msgBlockProposal) StopHandle() {
	fmt.Printf("stop ...\n")
}


// step2 (The First Step of the Graded Consensus Protocol GC) message
// step3 (The Second Step of GC) message
// step2 and step3 message has the same structure
// m(r,2) = (ESIG(v′), σr2),v′= H(Bℓr) OR emptyHash{}
type GradedConsensus struct {
	//hash is v′, the hash of the next block
	Hash          types.Hash    //the Br's hash
	Esig          []byte        //the signature of somebody's ephemeral secret key
	Credential    *CredentialSig
}

type msgGradedConsensus struct {
	gc    *GradedConsensus
	*message.MsgPriv
}

func NewMsgGradedConsensus(gc *GradedConsensus) *msgGradedConsensus{
	msgGc := &msgGradedConsensus{
		gc:      gc,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgGc)
	return msgGc
}

func (tm *msgGradedConsensus) DataHandle(data interface{}) {
	fmt.Println("msgGradedConsensus data handle")
}

func (tm *msgGradedConsensus) StopHandle() {
	fmt.Printf("stop ...\n")
}

// step4 and step other message
// m(r,s) = (ESIG(b), ESIG(v′), σrs)
type BinaryByzantineAgreement struct {
	//B is the BBA⋆ input b, 0 or 1
	B             uint
	EsigB         []byte
	//hash is v′, the hash of the next block
	Hash          types.Hash
	EsigV         []byte
	Credential    *CredentialSig
}

type msgBinaryByzantineAgreement struct {
	bba    *BinaryByzantineAgreement
	*message.MsgPriv
}

func NewMsgBinaryByzantineAgreement(bba *BinaryByzantineAgreement) *msgBinaryByzantineAgreement{
	msgBba := &msgBinaryByzantineAgreement{
		bba:      bba,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgBba)
	return msgBba
}

func (tm *BinaryByzantineAgreement) DataHandle(data interface{}) {
	fmt.Println("BinaryByzantineAgreement data handle")
}

func (tm *BinaryByzantineAgreement) StopHandle() {
	fmt.Printf("stop ...\n")
}
