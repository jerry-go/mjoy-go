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
	"errors"
	"fmt"
	"math/big"
	"mjoy.io/common"
	"mjoy.io/common/types"
	"mjoy.io/consensus/message"
	"mjoy.io/utils/event"
	"reflect"
	"sync"
	"mjoy.io/core/blockchain/block"
)

const (
	STEP_BP = iota + 0xffff
	STEP_REDUCTION_1
	STEP_REDUCTION_2
	STEP_FINAL
)

//go:generate msgp
func (cs *CredentialSign) validate() (types.Address, error) {
	//leader := false
	//if 1 == cs.Step{
	//	leader = true
	//}
	//hash := cs.Signature.hash()

	//verify right
	//if isPotVerifier(hash.Bytes(), leader) == false {
	//	return types.Address{}, errors.New("credential has no right to verify")
	//}

	//verify signature
	sender, err := cs.sender()
	if err != nil {
		return types.Address{}, errors.New(fmt.Sprintf("verify CredentialSig fail: %s", err))
	}

	return sender, nil
}

type msgCredentialSig struct {
	cs *CredentialSign
	*message.MsgPriv
}

func NewMsgCredential(c *CredentialSign) *msgCredentialSig {
	msgCs := &msgCredentialSig{
		cs:      c,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgCs)
	return msgCs
}

func (c *msgCredentialSig) DataHandle(data interface{}) {
	logger.Debug("msgCredentialSig data handle")
	if _, err := c.cs.validate(); err != nil {
		logger.Info("message CredentialSig validate error:", err)
		return
	}
	MsgTransfer().Send2Apos(c.cs)
}

func (c *msgCredentialSig) StopHandle() {
	logger.Debug("msgCredentialSig stop ...")
}

// step1 (Block Proposal) message
// m(r,1) = (Br, esig(H(Br)), σr1)
type BlockProposal struct {
	Block      *block.Block
	Esig       *EphemeralSign
	Credential *CredentialSign
}

func newBlockProposal() *BlockProposal {
	b := new(BlockProposal)
	b.Esig = new(EphemeralSign)
	return b
}

func (bp *BlockProposal) validate() error {
	//verify step
	if bp.Credential.Step != 1 {
		return errors.New(fmt.Sprintf("Block Proposal step is not 1: %d", bp.Credential.Step))
	}

	//verify Credential
	cretSender, err := bp.Credential.validate()
	if err != nil {
		return err
	}

	//verify ephemeral signature
	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = bp.Credential.Step
	bp.Esig.val = bp.Block.Hash().Bytes()
	sender, err := bp.Esig.sender()
	if err != nil {
		return errors.New(fmt.Sprintf("BP verify ephemeral signature fail: %s", err))
	}
	if cretSender != sender {
		logger.Debug("Block Proposal Ephemeral signature address is not equal to Credential signature address", sender.Hex(), cretSender.Hex())
		return errors.New("sender's address between Credential and Ephemeral is not equal")
	}

	//todo block validate

	return nil
}

type msgBlockProposal struct {
	bp *BlockProposal
	*message.MsgPriv
}

// new a message
func NewMsgBlockProposal(bp *BlockProposal) *msgBlockProposal {
	msgBp := &msgBlockProposal{
		bp:      bp,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgBp)
	return msgBp
}

func (bp *msgBlockProposal) DataHandle(data interface{}) {
	logger.Debug("msgBlockProposal data handle")
	if err := bp.bp.validate(); err != nil {
		logger.Info("message BlockProposal validate error:", err)
		return
	}
	MsgTransfer().Send2Apos(bp.bp)
}

func (bp *msgBlockProposal) StopHandle() {
	logger.Debug("msgBlockProposal stop ...")
}

// step2 (The First Step of the Graded Consensus Protocol GC) message
// step3 (The Second Step of GC) message
// step2 and step3 message has the same structure
// m(r,2) = (ESIG(v′), σr2),v′= H(Bℓr) OR emptyHash{}
type GradedConsensus struct {
	//hash is v′, the hash of the next block
	Hash       types.Hash     //the Br's hash
	Esig       *EphemeralSign //the signature of somebody's ephemeral secret key
	Credential *CredentialSign
}

func newGradedConsensus() *GradedConsensus {
	g := new(GradedConsensus)
	g.Esig = new(EphemeralSign)

	return g
}

func (gc *GradedConsensus) validate() error {
	step := gc.Credential.Step
	if step != 2 && step != 3 {
		return errors.New(fmt.Sprintf("Graded Consensus step is not 2 or 3: %d", gc.Credential.Step))
	}
	//verify Credential
	cretSender, err := gc.Credential.validate()
	if err != nil {
		return err
	}

	//verify ephemeral signature
	gc.Esig.round = gc.Credential.Round
	gc.Esig.step = gc.Credential.Step
	gc.Esig.val = gc.Hash.Bytes()
	sender, err := gc.Esig.sender()
	if err != nil {
		return errors.New(fmt.Sprintf("GC verify ephemeral signature fail: %s", err))
	}
	if cretSender != sender {
		logger.Debug("Graded Consensus Ephemeral signature address is not equal to Credential signature address", sender.Hex(), cretSender.Hex())
		return errors.New("sender's address between Credential and Ephemeral is not equal")
	}

	return nil
}

func (gc *GradedConsensus) GcHash() types.Hash {
	hash, err := common.MsgpHash(gc)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

type msgGradedConsensus struct {
	gc *GradedConsensus
	*message.MsgPriv
}

func NewMsgGradedConsensus(gc *GradedConsensus) *msgGradedConsensus {
	msgGc := &msgGradedConsensus{
		gc:      gc,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgGc)
	return msgGc
}

func (gc *msgGradedConsensus) DataHandle(data interface{}) {
	logger.Debug("msgGradedConsensus data handle")
	if err := gc.gc.validate(); err != nil {
		logger.Info("message GradedConsensus validate error:", err)
		return
	}
	MsgTransfer().Send2Apos(gc.gc)
}

func (gc *msgGradedConsensus) StopHandle() {
	logger.Debug("msgGradedConsensus stop ...")
}

// step4 and step other message
// m(r,s) = (ESIG(b), ESIG(v′), σrs)
type BinaryByzantineAgreement struct {
	//B is the BBA⋆ input b, 0 or 1
	B     uint
	EsigB *EphemeralSign
	//hash is v′, the hash of the next block
	Hash       types.Hash
	EsigV      *EphemeralSign
	Credential *CredentialSign
}

func newBinaryByzantineAgreement() *BinaryByzantineAgreement {
	b := new(BinaryByzantineAgreement)
	b.EsigB = new(EphemeralSign)
	b.EsigV = new(EphemeralSign)

	return b
}

func (bba *BinaryByzantineAgreement) validate() error {
	//verify step
	if bba.Credential.Step < 4 {
		return errors.New(fmt.Sprintf("Binary Byzantine Agreement step is not right: %d", bba.Credential.Step))
	}
	//verify Credential
	cretSender, err := bba.Credential.validate()
	if err != nil {
		return err
	}

	if bba.B > 1 {
		return errors.New(fmt.Sprintf("B value %d is not right in apos protocal", bba.B))
	}

	//for step m + 3
	if Config().maxBBASteps+3 == int(bba.Credential.Step) {
		// for step m +3, b must be 1 and v must be Hash(empty block(qr = last qr))
		if bba.B != 1 {
			logger.Info("bba m + 3 step message'b is not equal 1", bba.B)
			return errors.New("bba m + 3 step message'b is not equal 1")
		}
		//verify empty block hash
		if gCommonTools != nil {
			if gCommonTools.MakeEmptyBlock(makeEmptyBlockConsensusData(bba.Credential.Round)).Hash() != bba.Hash {
				logger.Info("m + 3 message hash is not empty block hash", err)
				return errors.New("m + 3 message hash is not empty block hash")
			}
		}
	}

	//verify B ephemeral signature
	bba.EsigB.round = bba.Credential.Round
	bba.EsigB.step = bba.Credential.Step
	bba.EsigB.val = big.NewInt(int64(bba.B)).Bytes()
	bSender, err := bba.EsigB.sender()
	if err != nil {
		return errors.New(fmt.Sprintf("BBA B verify ephemeral signature fail: %s", err))
	}

	if cretSender != bSender {
		logger.Debug("BinaryByzantineAgreement Ephemeral B signature address is not equal to Credential signature address", bSender.Hex(), cretSender.Hex())
		return errors.New("sender's address between Credential and B Ephemeral is not equal")
	}

	//verify V ephemeral signature
	bba.EsigV.round = bba.Credential.Round
	bba.EsigV.step = bba.Credential.Step
	bba.EsigV.val = bba.Hash.Bytes()
	hashSender, err := bba.EsigV.sender()
	if err != nil {
		return errors.New(fmt.Sprintf("BBA V verify ephemeral signature fail: %s", err))
	}
	if cretSender != hashSender {
		logger.Debug("BinaryByzantineAgreement Ephemeral V signature address is not equal to Credential signature address", hashSender.Hex(), cretSender.Hex())
		return errors.New("sender's address between Credential and V Ephemeral is not equal")
	}

	return nil
}

func (bba *BinaryByzantineAgreement) BbaHash() types.Hash {
	hash, err := common.MsgpHash(bba)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

type msgBinaryByzantineAgreement struct {
	bba *BinaryByzantineAgreement
	*message.MsgPriv
}

func NewMsgBinaryByzantineAgreement(bba *BinaryByzantineAgreement) *msgBinaryByzantineAgreement {
	msgBba := &msgBinaryByzantineAgreement{
		bba:     bba,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgBba)
	return msgBba
}

func (bba *msgBinaryByzantineAgreement) DataHandle(data interface{}) {
	logger.Debug("BinaryByzantineAgreement data handle")
	if err := bba.bba.validate(); err != nil {
		logger.Info("message BinaryByzantineAgreement validate error:", err)
		return
	}
	MsgTransfer().Send2Apos(bba.bba)
}

func (bba *msgBinaryByzantineAgreement) StopHandle() {
	logger.Debug("msgBinaryByzantineAgreement stop ...")
}


type ByzantineAgreementStar struct {
	Hash       types.Hash      //voted block's hash.
	Esig       *EphemeralSign  //the signature of somebody's ephemeral secret key
	Credential *CredentialSign
}

func newByzantineAgreementStar() *ByzantineAgreementStar {
	b := new(ByzantineAgreementStar)
	b.Esig = new(EphemeralSign)
	return b
}

func (ba *ByzantineAgreementStar) validate() error {
	//verify step
	if ba.Credential.Step < 1 || uint(ba.Credential.Step) > Config().maxStep{
		return errors.New(fmt.Sprintf("Byzantine Agreement Star step is not right: %d", ba.Credential.Step))
	}
	//verify Credential
	cretSender, err := ba.Credential.validate()
	if err != nil {
		return err
	}

	//verify ephemeral signature
	ba.Esig.round = ba.Credential.Round
	ba.Esig.step = ba.Credential.Step
	ba.Esig.val = ba.Hash.Bytes()
	sender, err := ba.Esig.sender()
	if err != nil {
		return errors.New(fmt.Sprintf("BA* verify ephemeral signature fail: %s", err))
	}

	if cretSender != sender {
		logger.Debug("BA* Ephemeral hash signature address is not equal to Credential signature address", sender.Hex(), cretSender.Hex())
		return errors.New("sender's address between Credential and Hash Ephemeral is not equal")
	}

	return nil
}

func (ba *ByzantineAgreementStar) BaHash() types.Hash {
	hash, err := common.MsgpHash(ba)
	if err != nil {
		return types.Hash{}
	}
	return hash
}

type msgByzantineAgreementStar struct {
	ba *ByzantineAgreementStar
	*message.MsgPriv
}

func NewMsgByzantineAgreementStar(ba *ByzantineAgreementStar) *msgByzantineAgreementStar {
	msgBba := &msgByzantineAgreementStar{
		ba:     ba,
		MsgPriv: message.NewMsgPriv(),
	}
	message.Msgcore().Handle(msgBba)
	return msgBba
}

func (ba *msgByzantineAgreementStar) DataHandle(data interface{}) {
	logger.Debug("msgByzantineAgreementStar data handle", ba.ba.Credential.Round, ba.ba.Credential.Step)
	if err := ba.ba.validate(); err != nil {
		logger.Info("message ByzantineAgreementStar validate error:", err)
		return
	}
	MsgTransfer().Send2Apos(ba.ba)
}

func (bba *msgByzantineAgreementStar) StopHandle() {
	logger.Debug("msgByzantineAgreementStar stop ...")
}

//message transfer between msg and Apos
type msgTransfer struct {
	receiveSubChan     chan dataPack
	somebodyGetSubChan bool
	receiveChan        chan dataPack //receive message from BBa, Gc, Bp and etc.
	sendChan           chan dataPack

	csFeed  event.Feed
	bpFeed  event.Feed
	gcFeed  event.Feed
	bbaFeed event.Feed
	scope   event.SubscriptionScope
}

// about MsgTransfer singleton
var (
	msgTransferInstance *msgTransfer
	msgTransferOnce     sync.Once
)

// get the MsgTransfer singleton
func MsgTransfer() *msgTransfer {
	msgTransferOnce.Do(func() {
		msgTransferInstance = &msgTransfer{
			receiveChan:        make(chan dataPack, 10),
			receiveSubChan:     make(chan dataPack, 10),
			sendChan:           make(chan dataPack, 10),
			somebodyGetSubChan: false,
		}
	})
	return msgTransferInstance
}

func (mt *msgTransfer) BroadCast(msg []byte) error {
	return nil
}

func (mt *msgTransfer) GetMsg() <-chan dataPack {
	return mt.receiveChan
}

func (mt *msgTransfer) GetDataMsg() <-chan dataPack {
	return mt.receiveChan
}

//return the chan sub chan,just for test
func (mt *msgTransfer) GetSubDataMsg() <-chan dataPack {
	mt.somebodyGetSubChan = true
	return mt.receiveSubChan
}

func (mt *msgTransfer) SendCredential(c *CredentialSign) error {
	return nil
}

func (mt *msgTransfer) PropagateCredential(c *CredentialSign) error {
	//logger.Debug("PropagateCredential", c.Round, c.Step)
	go mt.csFeed.Send(CsEvent{c})
	return nil
}

func (mt *msgTransfer) sendInner(data dataPack) {
	mt.receiveChan <- data

	//send the data to receiveSubCh
	if mt.somebodyGetSubChan {
		mt.receiveSubChan <- data
	}
}

func (mt *msgTransfer) SendInner(data dataPack) error {
	//todo here need to validate process??
	logger.Debug("SendInner type:", reflect.TypeOf(data))
	go mt.sendInner(data)

	return nil
}

func (mt *msgTransfer) PropagateMsg(data dataPack) error {
	logger.Debug("msgTransfer PropagateMsg in, data type:", reflect.TypeOf(data))
	switch v := data.(type) {
	case *CredentialSign:
		go mt.csFeed.Send(CsEvent{v})
	case *BlockProposal:
		go mt.bpFeed.Send(BpEvent{v})
	case *GradedConsensus:
		go mt.gcFeed.Send(GcEvent{v})
	case *BinaryByzantineAgreement:
		go mt.bbaFeed.Send(BbaEvent{v})
	default:
		logger.Warn("in PropagateMsg invalid message type ", reflect.TypeOf(v))
	}
	return nil
}

func (mt *msgTransfer) Send2Apos(data dataPack) {
	mt.receiveChan <- data
}

func (mt *msgTransfer) SubscribeCsEvent(ch chan<- CsEvent) event.Subscription {
	return mt.scope.Track(mt.csFeed.Subscribe(ch))
}
func (mt *msgTransfer) SubscribeBpEvent(ch chan<- BpEvent) event.Subscription {
	return mt.scope.Track(mt.bpFeed.Subscribe(ch))
}
func (mt *msgTransfer) SubscribeGcEvent(ch chan<- GcEvent) event.Subscription {
	return mt.scope.Track(mt.gcFeed.Subscribe(ch))
}
func (mt *msgTransfer) SubscribeBbaEvent(ch chan<- BbaEvent) event.Subscription {
	return mt.scope.Track(mt.bbaFeed.Subscribe(ch))
}

type CsEvent struct{ Cs *CredentialSign }
type BpEvent struct{ Bp *BlockProposal }
type GcEvent struct{ Gc *GradedConsensus }
type BbaEvent struct{ Bba *BinaryByzantineAgreement }
