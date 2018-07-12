package apos

import (
	"math/big"
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"
)

type StepCtxInterface interface {
	ESIG(types.Hash) []byte
	SendInner(dataPack) error
	PropagateMsg(dataPack) error
	GetCredential() *CredentialSig
	SetRound(*Round)
	makeEmptyBlockForTest() *block.Block
}

type StepContext struct {
	esig                  func(hash types.Hash) []byte
	sendInner             func(pack dataPack) error
	propagateMsg          func(dataPack) error
	getCredential         func() *CredentialSig
	setRound              func(*Round)
	makeEmptyBlockForTest func(cs *CredentialSig) *block.Block
}

func makeStepContext() *StepContext {
	s := new(StepContext)
	return s
}

type stepCtxData struct {
	apos        *Apos
	round       *Round
	pCredential *CredentialSig
}

func makeStepCtxData(apos *Apos, pCredential *CredentialSig) *stepCtxData {
	s := new(stepCtxData)
	s.apos = apos
	s.pCredential = pCredential
	return s
}

func (this *stepCtxData) ESIG(h types.Hash) (R, S, V *big.Int) {

	signature := MakeEmptySignature()
	sig := this.apos.commonTools.SigHash(h)
	R, S, V, err := signature.FillBySig(sig)
	if err != nil {
		logger.Error("signature.fillBySig wrong:", err.Error())
	}
	return R, S, V
}

func (this *stepCtxData) SendInner(dp dataPack) error {
	return this.apos.outMsger.SendInner(dp)
}

func (this *stepCtxData) GetCredential() *CredentialSig {
	pC := new(CredentialSig)
	pC.Round = this.pCredential.Round
	pC.Step = this.pCredential.Step
	pC.R = this.pCredential.R
	pC.S = this.pCredential.S
	pC.V = this.pCredential.V
	return pC
}

func (this *stepCtxData) SetRound(pRound *Round) {
	this.round = pRound
}

func (this *stepCtxData) makeEmptyBlockForTest(cs *CredentialSign) *block.Block {
	return this.apos.makeEmptyBlockForTest(cs)
}

func (this *stepCtxData) PropagateMsg(dp dataPack) error {
	return this.apos.outMsger.PropagateMsg(dp)
}
