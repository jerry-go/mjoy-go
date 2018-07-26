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
// @File: apos.go
// @Date: 2018/06/15 11:35:15
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"

	"mjoy.io/utils/crypto"
	"reflect"
	"sync"
	"time"

)

type peerMsgs struct {
	msgCs   map[int]*CredentialSign
	msgBas  map[int]*ByzantineAgreementStar

	//0 :default honesty peer. 1: malicious peer
	honesty uint
}

type pqCredential struct {
	pq          priorityQueue
	credentials map[string]*CredentialSign
}

type mainStepOutput struct {
	bp        types.Hash
	reduction types.Hash
	bba       types.Hash
	final     types.Hash

	mu        sync.Mutex
}

func (this *mainStepOutput) setBpResult(bp types.Hash) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.bp = bp
}
func (this *mainStepOutput) setReductionResult(reduction types.Hash) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.reduction = reduction
}
func (this *mainStepOutput) setBbaResult(bba types.Hash) bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.bba = bba
	nullHash := types.Hash{}

	if this.final != nullHash {
		return true
	}
	return false
}
func (this *mainStepOutput) setFinalResult(final types.Hash) bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.final = final
	nullHash := types.Hash{}
	if this.bba != nullHash {
		return true
	}
	return false
}

//round context
type Round struct {
	round uint64
	//condition 0 and condition 1 end number tH = 2n/3 + 1
	targetNum int

	apos *Apos

	credentials map[int]*CredentialSign

	smallestLBr *BlockProposal
	lock        sync.RWMutex

	emptyBlock    *block.Block
	maxLeaderNum  int
	curLeaderNum  int
	curLeaderDiff *big.Int
	curLeader     types.Hash


	msgs map[types.Address]*peerMsgs
	csPq map[int]*pqCredential

	quitCh      chan *block.Block
	roundOverCh chan interface{}

	bpObj   *BpObj
	voteObj  *VoteObj


	//version 1.1
	mainStepRlt   mainStepOutput
	parentHash    types.Hash
	countVote     *countVote
}

//gilad tools



func CalculatePriority(hash types.Hash , w , W ,t uint64 )uint64{
	pri := new(big.Int).SetBytes(hash[0:6]).Int64() % 13

	return uint64(pri)
}

func (this *Round)startVoteTimer(delay int){
	this.countVote.startTimer(delay)
}

func (this *Round)makeBlockConsensusData(bp *BlockProposal) *block.ConsensusData{
	return makeBlockConsensusData(bp, this.apos.commonTools)
}

func (this *Round)getCredentialByStep (step uint64)*CredentialSign{
	this.lock.RLock()
	defer this.lock.RUnlock()

	if c , ok := this.credentials[int(step)];ok {
		return c
	}
	return nil
}

func (this *Round)commonCoin (round , step , t uint64)uint64{
	return (round + step) % 2
}

func (this *Round)getAccountMonney(address types.Address , round uint64)uint64{
	monney := new(big.Int).SetBytes(address[2:6])	//4bytes
	return monney.Uint64() + round
}

func (this *Round)getTotalMonney(round uint64)uint64{
	monney := big.NewInt(0xffffffffff)
	return monney.Uint64() + round
}

func (this *Round)getBpThreshold()uint64{
	return uint64(Config().tProposer)
}

func (this *Round)getVoteThreshold()uint64{
	return uint64(Config().tStepThreshold)
}


func (this *Round)verifyBlock(b *block.Block)bool{
	lastHash := this.apos.commonTools.GetNowBlockHash()

	//here we just compare the parent hash is right or not
	if lastHash.Equal(&b.B_header.ParentHash){
		return true
	}
	return false
}


func (this *Round)sortition (hash types.Hash , t,w,W uint64 )uint64{
	//no need take VRFsk
	return CalculatePriority(hash , w,W,t)
}

func (this *Round)verifySort(cret CredentialSign , w, W,t uint64)uint64{

	//credential

	_ , err := cret.sender()
	if err != nil {
		return 0
	}
	//here should call interface
	return CalculatePriority(cret.Signature.Hash() , w , W , t)

}



func newRound(round int, parentHash types.Hash,apos *Apos, roundOverCh chan interface{}) *Round {
	r := new(Round)
	r.init(round, apos, roundOverCh)
	r.parentHash = parentHash
	return r
}

func (this *Round) getEmptyBlockHash() types.Hash {
	return this.emptyBlock.Hash()
}

func (this *Round) init(round int, apos *Apos, roundOverCh chan interface{}) {
	this.round = uint64(round)
	this.apos = apos
	this.roundOverCh = roundOverCh

	// this.maxLeaderNum = this.apos.algoParam.maxLeaderNum
	this.credentials = make(map[int]*CredentialSign)
	emptyBlock := this.apos.commonTools.MakeEmptyBlock(makeEmptyBlockConsensusData(this.round))
	this.emptyBlock = emptyBlock


	this.quitCh = make(chan *block.Block, 1)

	this.msgs = make(map[types.Address]*peerMsgs)
	this.csPq = make(map[int]*pqCredential)


	//step ctx init
	// step context
}

func (this *Round) setBpResult(hash types.Hash) {

	logger.Info("round", this.round ,"setBpResult", hash.String())

	this.mainStepRlt.setBpResult(hash)
}
func (this *Round) setReductionResult(hash types.Hash) {
	logger.Info("round", this.round, "setReductionResult", hash.String())
	this.mainStepRlt.setReductionResult(hash)
}

func (this *Round) setBbaResult(hash types.Hash) {
	logger.Info("round", this.round, "setBbaResult", hash.String())
	complete := this.mainStepRlt.setBbaResult(hash)
	if complete {
		if hash == this.mainStepRlt.final {
			logger.Info("Final consensus!!!")
		} else {
			logger.Info("Tentative consensus!!!")
		}
		consensusBlock := this.bpObj.getExistBlock(hash)
		if consensusBlock == nil {
			logger.Error(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "SetFianl getExistBlock == nil , hash=" , this.mainStepRlt.bba.Hex(), COLOR_SHORT_RESET)
		}
		if hash == this.getEmptyBlockHash() {
			consensusBlock = this.emptyBlock
		}
		if consensusBlock != nil {
			this.quitCh <- consensusBlock
		} else {
			logger.Error(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "SetFinalResult Get a nil block" , this.mainStepRlt.bba.Hex(), COLOR_SHORT_RESET)
			//todo need download this block based on hash
		}
	}
}

func (this *Round) setFinalResult(hash types.Hash) {
	logger.Info("round", this.round, "setFinalResult", hash.String())
	complete :=this.mainStepRlt.setFinalResult(hash)
	if complete {
		if hash == this.mainStepRlt.bba {
			logger.Info("Final consensus!!!")
		} else {
			logger.Info("Tentative consensus!!!")
		}
		consensusBlock := this.bpObj.getExistBlock(this.mainStepRlt.bba)
		if consensusBlock == nil {
			logger.Error(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "SetFianl getExistBlock == nil , hash=" , this.mainStepRlt.bba.Hex(), COLOR_SHORT_RESET)
		}
		if this.mainStepRlt.bba == this.getEmptyBlockHash() {
			consensusBlock = this.emptyBlock
		}
		if consensusBlock != nil {
			this.quitCh <- consensusBlock
		} else {
			logger.Error(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "SetFinalResult Get a nil block" , this.mainStepRlt.bba.Hex(), COLOR_SHORT_RESET)
			//todo need download this block based on hash
		}
	}else{
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "setFinalResult,but not complete......", COLOR_SHORT_RESET)
	}
}




//inform stepObj to stop running
func (this *Round) broadCastStop() {
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "In BroadCastStop...", COLOR_SHORT_RESET)
	this.bpObj.stop()
	this.voteObj.stop()
	this.countVote.stop()
}

// Generate valid Credentials in current round
func (this *Round) generateCredentials() {
	for i := 1; i < int(Config().maxStep); i++ {
		credential := this.apos.makeCredential(i)
		isVerfier := this.apos.judgeVerifier(credential, i)
		//logger.Info("GenerateCredential step:",i,"  isVerifier:",isVerfier)
		if isVerfier {
			logger.Info("GenerateCredential step:", i, "  votes:", credential.votes)
			this.credentials[i] = credential
			this.apos.commonTools.CreateTmpPriKey(i)
		}
	}

	for i := STEP_BP; i < STEP_IDLE; i++ {
		credential := this.apos.makeCredential(i)
		isVerfier := this.apos.judgeVerifier(credential, i)
		//logger.Info("GenerateCredential step:",i,"  isVerifier:",isVerfier)
		if isVerfier {
			logger.Info("GenerateCredential step:", i, "  votes:", credential.votes)
			this.credentials[i] = credential
			this.apos.commonTools.CreateTmpPriKey(i)
		}
	}
}

func (this *Round) broadcastCredentials() {
	for i, credential := range this.credentials {
		_ = i
		logger.Info("SendCredential round", this.round, "step", i)
		this.apos.outMsger.SendInner(credential)
	}
}

func (this *Round) startStepObjs(wg *sync.WaitGroup) {

	stepCtx := &stepCtx{}

	stepCtx.setBpResult = this.setBpResult
	stepCtx.setReductionResult = this.setReductionResult
	stepCtx.setBbaResult = this.setBbaResult
	stepCtx.setFinalResult = this.setFinalResult

	//ctx for new step obj

	stepCtx.sortition = this.sortition
	stepCtx.verifyBlock = this.verifyBlock
	stepCtx.verifySort = this.verifySort
	stepCtx.getCredentialByStep = this.getCredentialByStep
	stepCtx.getAccountMonney = this.getAccountMonney
	stepCtx.getTotalMonney = this.getTotalMonney
	stepCtx.getBpThreshold = this.getBpThreshold
	stepCtx.getVoteThreshold = this.getVoteThreshold
	stepCtx.startVoteTimer = this.startVoteTimer
	stepCtx.commonCoin = this.commonCoin
	stepCtx.getProducerNewBlock = this.apos.commonTools.GetProducerNewBlock
	stepCtx.makeBlockConsensusData = this.makeBlockConsensusData


	roundRt := this.round
	stepCtx.getRound = func() uint64 {
		return roundRt
	}

	stepCtx.esig = this.apos.commonTools.Esig
	stepCtx.sendInner = this.apos.outMsger.SendInner
	stepCtx.propagateMsg = this.apos.outMsger.PropagateMsg
	stepCtx.getEmptyBlockHash = this.getEmptyBlockHash


	this.bpObj = makeBpObj(stepCtx)
	this.voteObj = makeVoteObj(stepCtx)


	sendVoteData := func(step int, hash types.Hash) {

		this.voteObj.SendVoteData(this.round, uint64(step), hash)
	}
	this.countVote = newCountVote(sendVoteData, this.emptyBlock.Hash())
	if this.countVote == nil {
		logger.Error("this.countVote == nil...........")
	}
	go this.bpObj.run()
	go this.voteObj.run()
	go this.countVote.run()

}


func (this *Round) filterMsgCs(msg *CredentialSign) error {

	address, err := msg.sender()
	if err != nil {
		return err
	}
	step := msg.Step

	if peermsgs, ok := this.msgs[address]; ok {
		if peermsgs.honesty == 1 {
			return errors.New("not honesty peer")
		}
		if mCs, ok := peermsgs.msgCs[int(step)]; ok {
			if mCs.Step == step {
				return errors.New("duplicate message Credential")
			}
		} else {
			peermsgs.msgCs[int(step)] = msg
		}
	} else {
		ps := &peerMsgs{
			msgBas:  make(map[int]*ByzantineAgreementStar),
			msgCs:   make(map[int]*CredentialSign),
			honesty: 0,
		}
		ps.msgCs[int(step)] = msg
		this.msgs[address] = ps
	}
	return nil
}


// process the Credential message
func (this *Round) receiveMsgCs(msg *CredentialSign) {
	logger.Debug("Receive message CredentialSign [r:s]:", msg.Round, msg.Step)
	if msg.Round != this.round {
		logger.Warn("verify fail, Credential msg is not in current round,want:", msg.Round, "  but:",this.round)
		return
	}

	//duplicate message check
	if err := this.filterMsgCs(msg); err != nil {
		logger.Info("filter Credential fail", err)
		return
	}
	//Propagate message via p2p
	this.apos.outMsger.PropagateMsg(msg)
}


func (this *Round) receiveMsgBp(msg *BlockProposal) {
	//verify msg
	if msg.Credential.Round != this.round {
		logger.Warn("verify fail, BlockProposal msg is not in current round", msg.Credential.Round, this.round)
		return
	}

	this.bpObj.sendMsg(msg)
	// for BP Propagate process will in stepObj

}

func (this *Round) filterMsgBa(msg *ByzantineAgreementStar) error {
	address, err := msg.Credential.sender()
	if err != nil {
		return err
	}
	step := msg.Credential.Step

	if peerMsgBas, ok := this.msgs[address]; ok {
		if peerMsgBas.honesty == 1 {
			return errors.New("not honesty peer")
		}
		if peerba, ok := peerMsgBas.msgBas[int(step)]; ok {
			if peerba.Hash == msg.Hash {
				return errors.New("duplicate message ByzantineAgreementStar")
			} else {
				peerMsgBas.honesty = 1
				return errors.New("receive different hash in BA message, it must a malicious peer")
			}
		} else {
			peerMsgBas.msgBas[int(step)] = msg
		}
	} else {
		ps := &peerMsgs{
			msgBas:  make(map[int]*ByzantineAgreementStar),
			msgCs:   make(map[int]*CredentialSign),
			honesty: 0,
		}
		ps.msgBas[int(step)] = msg
		this.msgs[address] = ps
	}
	return nil
}
func (this *Round) receiveMsgBaStar(msg *ByzantineAgreementStar) {
	//verify msg
	if msg.Credential.Round != this.round {
		logger.Warn("verify fail, ba msg is not in current round", msg.Credential.Round, this.round)
		return
	}

	if msg.Credential.ParentHash != this.parentHash {
		logger.Warn("verify fail, ba msg is not in current block chain", msg.Credential.ParentHash.String(), this.parentHash.String())
		return
	}

	if err := this.filterMsgBa(msg); err != nil {
		logger.Info("filter ba message fail:", err)
		return
	}


	this.countVote.sendMsg(msg)


	//Propagate message via p2p
	this.apos.outMsger.PropagateMsg(msg)
}

func (this *Round) commonProcess() {
	for {
		select {
		// receive message
		case outData := <-this.apos.outMsger.GetDataMsg():
			switch v := outData.(type) {
			case *CredentialSign:
				this.receiveMsgCs(v)
			case *BlockProposal:
				this.receiveMsgBp(v)
			case *ByzantineAgreementStar:
				this.receiveMsgBaStar(v)
			default:
				logger.Warn("invalid message type ", reflect.TypeOf(v))
			}
		case consensusBlock := <-this.quitCh:
			fmt.Println("CommonProcess end block:", consensusBlock)
			bs := block.Blocks{}
			bs = append(bs, consensusBlock)
			_, err := this.apos.commonTools.InsertChain(bs)
			fmt.Println("InsertOneBlock    ErrStatus:", err)

			logger.Info("round exit ")
			this.broadCastStop()
			return
		}
	}
}

func (this *Round) run() {

	wg := sync.WaitGroup{}
	logger.Debug("run()......step1")
	// make verifiers Credential
	this.generateCredentials()

	// broadcast Credentials
	this.broadcastCredentials()

	this.startStepObjs(&wg)
	logger.Debug("run()......step2")
	this.commonProcess()
	wg.Wait()
	this.roundOverCh <- 1 //inform the caller,the mission complete
}

type Apos struct {
	systemParam interface{} //the difference of algoParam and systemParam is that algoParam show the Apos
	//running status,but the systemParam show the Mjoy runing
	mainStep    int
	commonTools CommonTools
	outMsger    OutMsger

	//all goroutine send msg to Apos by this Chan
	allMsgBridge chan dataPack

	roundCtx *Round

	roundOverCh chan interface{}
	aposStopCh  chan interface{} //for test if apos just deal once
	stop        bool
	lock        sync.RWMutex
}

//Create Apos
func NewApos(msger OutMsger, cmTools CommonTools) *Apos {
	logger.Debug("NewApos....................")
	a := new(Apos)
	//a.outMsger = msger
	a.commonTools = cmTools
	gCommonTools = cmTools
	a.allMsgBridge = make(chan dataPack, 10000)
	a.roundOverCh = make(chan interface{}, 1)
	a.aposStopCh = make(chan interface{}, 1)
	a.outMsger = MsgTransfer()

	a.reset()

	return a
}

func (this *Apos) SetPriKey(priKey *ecdsa.PrivateKey) {
	this.commonTools.SetPriKey(priKey)
}

func (this *Apos) makeEmptyBlockForTest(cs *CredentialSign) *block.Block {
	header := &block.Header{Number: types.NewBigInt(*big.NewInt(int64(this.commonTools.GetNextRound()))), Time: types.NewBigInt(*big.NewInt(time.Now().Unix())),
		ParentHash: this.commonTools.GetNowBlockHash()}
	//chainId := big.NewInt(100)
	//signer := block.NewBlockSigner(chainId)
	srcBytes := []byte{}
	srcBytes = append(srcBytes, cs.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes, cs.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes, cs.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)
	header.ConsensusData.Id = ConsensusDataId
	header.ConsensusData.Para = h
	signature := MakeEmptySignature()

	sig := this.commonTools.SigHash(header.HashNoSig())
	if sig == nil {
		logger.Error("sig == nil")
		return nil
	}
	R, S, V, err := signature.FillBySig(sig)
	if err != nil {
		logger.Error("makeEmptyBlockForTest Err:", err.Error())
		return nil
	}

	header.R = &types.BigInt{*R}
	header.S = &types.BigInt{*S}
	header.V = &types.BigInt{*V}

	b := block.NewBlock(header, nil, nil)
	return b
}

func (this *Apos) SetOutMsger(outMsger OutMsger) {
	this.outMsger = outMsger
}

func SetTestConfig() {
	//set config
	Config().R = 1000
	Config().tProposer = 1
	Config().tStep = 2000
	Config().tStepThreshold = 1
	Config().tFinal = 10000
	Config().tFinalThreshold = 1
	Config().maxStep = 15
	Config().delayPriority = 5
	Config().delayStep = 5
	Config().delayBlock = 10
	Config().delayStepVar = 5
}

//this is the main loop of Apos
func (this *Apos) Run() {
	SetTestConfig()

	//start round
	//this.roundOverCh<-1
	fmt.Println("Apos Run round:", this.commonTools.GetNextRound())
	this.roundCtx = newRound(this.commonTools.GetNextRound(), this.commonTools.GetNowBlockHash(), this, this.roundOverCh)

	go this.roundCtx.run()
	logger.Info("Apos is running.....")
	for {
		select {
		case <-this.roundOverCh:
			//logger.Info("round overs ", this.roundCtx.round)
			//this.aposStopCh<-1
			//return //if apos deal once ,stop it
			logger.Debug("Apos New Round Running...............")
			this.roundCtx = newRound(this.commonTools.GetNextRound(), this.commonTools.GetNowBlockHash(), this, this.roundOverCh)
			go this.roundCtx.run()
		}
	}
}

//reset the status of Apos
func (this *Apos) reset() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.mainStep = 0
	this.stop = false
}

//Create The Credential
func (this *Apos) makeCredential(s int) *CredentialSign {
	r := this.commonTools.GetNextRound()
	c := new(CredentialSign)
	c.Signature.init()
	c.Round = uint64(r)
	c.Step = uint64(s)
	c.ParentHash = this.commonTools.GetNowBlockHash()


	err := this.commonTools.Sig(c)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	c.votes = uint(CalculatePriority(c.Signature.Hash() , 0,0,0))
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_GREEN+COLOR_SUFFIX , "***Credential Votes Show:  Round:",c.Round , " Step:" , c.Step , "  Votes:" , c.votes , COLOR_SHORT_RESET)

	return c
}

func (this *Apos) StopCh() chan interface{} {
	return this.aposStopCh
}

func (this *Apos) judgeVerifier(cs *CredentialSign, setp int) bool {

	//h := cs.Signature.Hash()
	//leader := false
	//if 1 == setp {
	//	leader = true
	//}
	//return isPotVerifier(h.Bytes(), leader)
	if cs.votes > 0 {
		return true
	} else {
		return false
	}
}
