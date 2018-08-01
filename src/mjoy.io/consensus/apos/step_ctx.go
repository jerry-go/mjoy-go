package apos

import (
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"
)


//stepCtx contains all functions the stepObj will use
type stepCtx struct {
	getStep   func() int // get the number of step in the round
	getRound  func() uint64
	stopStep  func() // stop the step
	stopRound func() // stop all the step in the round, and end the round

	//getCredential func() signature
	//getEphemeralSig func(signed []byte) signature
	esig                  func(pEphemeralSign *EphemeralSign) error
	sendInner             func(pack dataPack) error
	propagateMsg          func(dataPack) error
	getCredential         func() *CredentialSign
	setRound              func(*Round)
	makeEmptyBlockForTest func(cs *CredentialSign) *block.Block
	getEmptyBlockHash     func() types.Hash
	getEphemeralSig       func(signed []byte) Signature
	getProducerNewBlock   func(data *block.ConsensusData) *block.Block
	//getPrivKey

	//gilad
	commonCoin func(round , step , t uint64)uint64  //x
	writeRet func(data *VoteData)                   //x
	sortition func(hash types.Hash , t,w,W uint64)uint64
	verifyBlock func(b *block.Block)bool
	verifySort func(cret CredentialSign , w, W,t uint64)uint64
	getCredentialByStep   func(step uint64)*CredentialSign
	getAccountMonney func (address types.Address , round uint64)uint64
	getTotalMonney func(round uint64)uint64
	getBpThreshold func()uint64
	getVoteThreshold func()uint64

	startVoteTimer func(delay int)
	makeBlockConsensusData func(bp *BlockProposal) *block.ConsensusData

	setBpResult func(hash types.Hash)
	setReductionResult func(hash types.Hash)
	setBbaResult  func(hash types.Hash)
	setFinalResult func(hash types.Hash)

}