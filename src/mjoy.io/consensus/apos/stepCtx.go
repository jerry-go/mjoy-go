package apos

import (
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"
)

type StepCtxInterface interface {
	ESIG( types.Hash)([]byte)
	SendInner( dataPack)error
	PropagateMsg(dataPack)error
	GetCredential()*CredentialSig
	SetRound(*Round)
	makeEmptyBlockForTest()*block.Block
}



type stepCtxData struct {
	apos        *Apos
	round       *Round
	pCredential *CredentialSig
}

func makeStepCtxData(apos *Apos , pCredential *CredentialSig)*stepCtxData{
	s := new(stepCtxData)
	s.apos = apos
	s.pCredential = pCredential
	return s
}



func (this *stepCtxData) ESIG(h types.Hash) ([]byte) {
	return this.apos.commonTools.ESIG(h)
}

func (this *stepCtxData) SendInner(dp dataPack) error {
	return this.apos.outMsger.SendInner(dp)
}

func (this *stepCtxData) GetCredential() *CredentialSig {
	pC := new(CredentialSig)
	pC.Round    = this.pCredential.Round
	pC.Step     = this.pCredential.Step
	pC.R        = this.pCredential.R
	pC.S        = this.pCredential.S
	pC.V        = this.pCredential.V
	return pC
}

func (this *stepCtxData)SetRound(pRound *Round){
	this.round = pRound
}

func (this *stepCtxData)makeEmptyBlockForTest()*block.Block{
	return this.apos.makeEmptyBlockForTest()
}

func (this *stepCtxData)PropagateMsg(dp dataPack)error{
	return this.apos.outMsger.PropagateMsg(dp)
}