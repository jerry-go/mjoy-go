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
	"errors"
	"math/big"
)
//go:generate msgp
func (cs *CredentialSign) validate() error{
	leader := false
	if 1 == cs.Step{
		leader = true
	}
	hash := cs.hash()

	//verify right
	if isPotVerifier(hash.Bytes(), leader) == false {
		return errors.New("credential has no right to verify")
	}

	//verify signature
	if _, err := cs.sender(); err != nil {
		return errors.New(fmt.Sprintf("verify CredentialSig fail: %s", err))
	}

	return nil
}

type msgCredentialSig struct {
	cs    *CredentialSign
	*message.MsgPriv
}

func NewMsgCredential(c *CredentialSign) *msgCredentialSig{
	msgCs := &msgCredentialSig{
		cs:      c,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgCs)
	return msgCs
}



func (c *msgCredentialSig) DataHandle(data interface{}) {
	fmt.Println("msgBlockProposal data handle")
}

func (c *msgCredentialSig) StopHandle() {
	fmt.Printf("stop ...\n")
}

// step1 (Block Proposal) message
// m(r,1) = (Br, esig(H(Br)), σr1)
type BlockProposal struct {
	Block         *block.Block
	Esig          *EphemeralSign
	Credential    *CredentialSign
}

func (bp *BlockProposal) validate() error{
	//verify Credential
	if err := bp.Credential.validate(); err != nil {
		return err
	}

	//verify ephemeral signature
	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = bp.Block.Hash().Bytes()
	if _, err := bp.Esig.sender(); err != nil {
		return errors.New(fmt.Sprintf("BP verify ephemeral signature fail: %s", err))
	}

	return nil
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


func (bp *msgBlockProposal) DataHandle(data interface{}) {
	fmt.Println("msgBlockProposal data handle")
}

func (bp *msgBlockProposal) StopHandle() {
	fmt.Printf("stop ...\n")
}


// step2 (The First Step of the Graded Consensus Protocol GC) message
// step3 (The Second Step of GC) message
// step2 and step3 message has the same structure
// m(r,2) = (ESIG(v′), σr2),v′= H(Bℓr) OR emptyHash{}
type GradedConsensus struct {
	//hash is v′, the hash of the next block
	Hash          types.Hash    //the Br's hash
	Esig          *EphemeralSign     //the signature of somebody's ephemeral secret key
	Credential    *CredentialSign
}

func (gc *GradedConsensus) validate() error{
	//verify Credential
	if err := gc.Credential.validate(); err != nil {
		return err
	}

	//verify ephemeral signature
	gc.Esig.round = gc.Credential.Round
	gc.Esig.step = gc.Credential.Step
	gc.Esig.val = gc.Hash.Bytes()
	if _, err := gc.Esig.sender(); err != nil {
		return errors.New(fmt.Sprintf("GC verify ephemeral signature fail: %s", err))
	}

	return nil
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

func (gc *msgGradedConsensus) DataHandle(data interface{}) {
	fmt.Println("msgGradedConsensus data handle")
}

func (gc *msgGradedConsensus) StopHandle() {
	fmt.Printf("stop ...\n")
}

// step4 and step other message
// m(r,s) = (ESIG(b), ESIG(v′), σrs)
type BinaryByzantineAgreement struct {
	//B is the BBA⋆ input b, 0 or 1
	B             uint
	EsigB         *EphemeralSign
	//hash is v′, the hash of the next block
	Hash          types.Hash
	EsigV         *EphemeralSign
	Credential    *CredentialSign
}

func (bba *BinaryByzantineAgreement) validate() error{
	//verify Credential
	if err := bba.Credential.validate(); err != nil {
		return err
	}

	if bba.B > 1 {
		return errors.New(fmt.Sprintf("B value %d is not right in apos protocal", bba.B))
	}

	//for step m + 3
	if Config().maxBBASteps + 3 == int(bba.Credential.Step) {
		// for step m +3, b must be 1 and v must be Hash(empty block(qr = last qr))
		if bba.B != 1 {
			logger.Info("bba m + 3 step message'b is not equal 1", bba.B)
			return errors.New("bba m + 3 step message'b is not equal 1")
		}
		// todo verify empty block hash, need get right empty block
		//if v.apos.makeEmptyBlockForTest().Hash() != msg.Hash {
		//	logger.Info("m + 3 message hash is not empty block hash", err)
		//	return errors.New("m + 3 message hash is not empty block hash")
		//}
	}

	//verify B ephemeral signature
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	if _, err := bba.EsigB.sender(); err != nil {
		return errors.New(fmt.Sprintf("BBA B verify ephemeral signature fail: %s", err))
	}

	//verify V ephemeral signature
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	if _, err := bba.EsigV.sender(); err != nil {
		return errors.New(fmt.Sprintf("BBA B verify ephemeral signature fail: %s", err))
	}

	return nil
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

func (bba *BinaryByzantineAgreement) DataHandle(data interface{}) {
	fmt.Println("BinaryByzantineAgreement data handle")
}

func (bba *BinaryByzantineAgreement) StopHandle() {
	fmt.Printf("stop ...\n")
}
