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
	"sync"
	"mjoy.io/utils/crypto"
	"mjoy.io/common/types"
	"math/big"
	"reflect"
	"errors"
	"mjoy.io/core/blockchain/block"
	"time"
)

const (
	IDLE  = iota
	ENDCONDITION0
	ENDCONDITION1
	ENDMAX
)

/*
Instructions For Apos:
The Type-Apos manage the main loop of Apos consensus,
handle the Condition0,Condition1 and m+3,
and Output Apos-SystemParam to the sub-goroutine
*/

type VoteInfo struct {
	sum           int
}

// Potential Leader used for judge End condition 0 and 1
type PotentialLeader struct {
	bp            *BlockProposal
	stepMsg       map[uint]*VoteInfo
}

func (this *PotentialLeader)AddVoteNumber(step uint, b uint) int {
	if _, ok := this.stepMsg[step]; !ok {
		vi := &VoteInfo{1}
		this.stepMsg[step] = vi
		return 1
	}
	this.stepMsg[step].sum++
	return this.stepMsg[step].sum
}

type peerBba struct {
	bba         *BinaryByzantineAgreement
	//B 1: msg.b == 0; 2: msg.b == 1; 3 means that receive two messages with different B
	B           uint
}


type peerMsgs struct {
	msgBbas        map[int]*peerBba
	msgGcs         map[int]*GradedConsensus
	msgCs          map[int]*CredentialSign

	//0 :default honesty peer. 1: malicious peer
	honesty        uint
}

//round context
type Round struct {
	round          uint64
	//condition 0 and condition 1 end number tH = 2n/3 + 1
	targetNum      int

	apos           *Apos

	credentials    map[int]*CredentialSign

	allStepObj     map[int]*stepRoutine

	smallestLBr    *BlockProposal
	lock           sync.RWMutex

	leaders        map[types.Hash]*PotentialLeader
	maxLeaderNum   int
	curLeaderNum   int
	curLeaderDiff  *big.Int
	curLeader      types.Hash

	msgs           map[types.Address]*peerMsgs

	quitCh         chan interface{}
	roundOverCh    chan interface{}
}

func newRound(round int , apos *Apos , roundOverCh chan interface{})*Round{
	r := new(Round)
	r.init(round,apos,roundOverCh)
	return r
}

func (this *Round)init(round int , apos *Apos , roundOverCh chan interface{}){
	this.round = uint64(round)
	this.apos = apos
	this.roundOverCh = roundOverCh

	// this.maxLeaderNum = this.apos.algoParam.maxLeaderNum
	this.credentials = make(map[int]*CredentialSign)
	this.allStepObj = make(map[int]*stepRoutine)
	this.leaders = make(map[types.Hash]*PotentialLeader)
	this.quitCh = make(chan interface{} , 1)
	this.msgs = make(map[types.Address]*peerMsgs)
}

func (this *Round)setSmallestBrM1(bp *BlockProposal){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.smallestLBr = bp
}

func (this *Round)addStepObj(step int , stepObj *stepRoutine){
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.allStepObj[step]; !ok {
		this.allStepObj[step] = stepObj
	}
}

//inform stepObj to stop running
func (this *Round)broadCastStop(){
	for _,v := range this.allStepObj {
		v.stop()
	}
}

// Generate valid Credentials in current round
func (this *Round)generateCredentials() {
	for i := 1; i <= Config().maxBBASteps + 3; i++{
		credential := this.apos.makeCredential(i)
		isVerfier := this.apos.judgeVerifier(credential, i)
		//logger.Info("GenerateCredential step:",i,"  isVerifier:",isVerfier)
		if isVerfier {
			logger.Info("GenerateCredential step:",i,"  isVerifier:",isVerfier)
			this.credentials[i] = credential
		}
	}
}

func (this *Round)broadcastCredentials() {
	for i, credential := range this.credentials {
		_ = i
		//logger.Info("SendCredential round", this.round.IntVal.Uint64(), "step", i)
		this.apos.outMsger.SendCredential(credential)
	}
}


func (this *Round)startVerify(wg *sync.WaitGroup) {
	for i, credential := range this.credentials {
		stepRoutineObj,stepObj := this.apos.stepsFactory(i, credential)
		if stepRoutineObj != nil && stepObj != nil{
			this.addStepObj(i, stepRoutineObj)
			go stepRoutineObj.run(stepObj)
		}
	}
}

func (this *Round)filterMsgCs(msg *CredentialSign) error {
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
				return errors.New("duplicate message Credential Signature")
			}
		} else {
			peermsgs.msgCs[int(step)] = msg
		}
	} else {
		ps := &peerMsgs{
			msgBbas: make(map[int]*peerBba),
			msgGcs: make(map[int]*GradedConsensus),
			msgCs: make(map[int]*CredentialSign),
			honesty: 0,
		}
		ps.msgCs[int(step)] = msg
		this.msgs[address] = ps
	}
	return nil
}
// hear M0 is the Credential message
func (this *Round)receiveMsgCs(msg *CredentialSign) {
	logger.Info("Receive message CredentialSign")
	if err := this.filterMsgCs(msg); err != nil {
		logger.Info("filter m0 fail", err)
		return
	}
	//Propagate message via p2p
	this.apos.outMsger.PropagateCredential(msg)
}

func (this *Round)saveBp(msg *BlockProposal) {
	hash := msg.Block.Hash()
	if _, ok := this.leaders[hash]; !ok {
		pleader := &PotentialLeader{msg,make(map[uint]*VoteInfo)}
		this.leaders[hash] = pleader
		this.curLeaderNum++
	}
}

func (this *Round)receiveMsgBp(msg *BlockProposal) {
	//verify msg
	if msg.Credential.Round != this.round {
		logger.Warn("verify fail, BlockProposal msg is not in current round", msg.Credential.Round, this.round)
		return
	}

	//send this msg to step2 goroutine
	if stepObj, ok := this.allStepObj[2]; ok {
		stepObj.sendMsg(msg)
	}

	// for M1 Propagate process will in stepObj

	this.saveBp(msg)
}


func (this *Round)filterMsgGc(msg *GradedConsensus) error {
	address, err := msg.Credential.sender()
	if err != nil {
		return err
	}
	step := msg.Credential.Step

	if peerMsgGcs, ok := this.msgs[address]; ok {
		if peerMsgGcs.honesty == 1 {
			return errors.New("not honesty peer")
		}

		if gc, ok := peerMsgGcs.msgGcs[int(step)]; ok {
			if gc.Hash == msg.Hash {
				return errors.New("duplicate message m23")
			} else {
				peerMsgGcs.honesty = 1
				return errors.New("receive different vote message m23, it must a malicious peer")
			}
		} else {
			peerMsgGcs.msgGcs[int(step)] = msg
		}
	} else {
		ps := &peerMsgs{
			msgBbas: make(map[int]*peerBba),
			msgGcs: make(map[int]*GradedConsensus),
			msgCs: make(map[int]*CredentialSign),
			honesty: 0,
		}
		ps.msgGcs[int(step)] = msg
		this.msgs[address] = ps
	}
	return nil
}

func (this *Round)receiveMsgGc(msg *GradedConsensus) {
	//verify msg
	if msg.Credential.Round != this.round {
		logger.Warn("verify fail, GradedConsensus msg is not in current round", msg.Credential.Round, this.round)
		return
	}

	if err := this.filterMsgGc(msg); err != nil {
		logger.Info("filter m23 fail", err)
		return
	}

	//send this msg to step3 or step4 goroutine
	step := msg.Credential.Step
	if stepObj, ok := this.allStepObj[int(step) + 1]; ok {
		stepObj.sendMsg(msg)
	}

	//Propagate message via p2p
	this.apos.outMsger.PropagateMsg(msg)
	logger.Info("Propagate Graded Consensus message via p2p")
}


func (this *Round)filterMsgBba(msg *BinaryByzantineAgreement) error {
	address, err := msg.Credential.sender()
	if err != nil {
		return err
	}
	step := msg.Credential.Step

	if peerMsgBbas, ok := this.msgs[address]; ok {
		if peerMsgBbas.honesty == 1 {
			return errors.New("not honesty peer")
		}

		if peerbba, ok := peerMsgBbas.msgBbas[int(step)]; ok {
			if peerbba.bba.Hash == msg.Hash && (peerbba.B == 3 || peerbba.B == msg.B + 1){
				return errors.New("duplicate bba message")
			} else if (peerbba.bba.Hash == msg.Hash) {
				// for bba message, player j can send different B value
				peerbba.B = 3
				logger.Info("receive different vote bba message!", msg.B)
				return nil
			} else {
				peerMsgBbas.honesty = 1
				return errors.New("receive different hash in BBA message, it must a malicious peer")
			}
		} else {
			msgPeer := &peerBba{msg, msg.B + 1}
			peerMsgBbas.msgBbas[int(step)] = msgPeer
		}
	} else {
		ps := &peerMsgs{
			msgBbas: make(map[int]*peerBba),
			msgGcs: make(map[int]*GradedConsensus),
			msgCs: make(map[int]*CredentialSign),
			honesty: 0,
		}
		msgPeer := &peerBba{msg, msg.B + 1}
		ps.msgBbas[int(step)] = msgPeer
		this.msgs[address] = ps
	}
	return nil
}

func (this *Round)endCondition(voteNum int, b uint) int {

	//if voteNum >= this.targetNum {
	if isAbsHonest(voteNum, false) {
		logger.Info("end condition ", b, "vote number", voteNum)
		if 0 == b{
			return ENDCONDITION0
		} else if 1 == b {
			return ENDCONDITION1
		} else {
			return ENDMAX
		}
	} else {
		return IDLE
	}
}

func (this *Round)saveMsgBba(msg *BinaryByzantineAgreement) int {
	hash := msg.Hash
	if pleader, ok := this.leaders[hash]; ok {
		step := msg.Credential.Step
		b := msg.B
		voteNum := 0
		if ((step + 1 -2 ) % 3 == 0) && (0 == b) {
			voteNum = pleader.AddVoteNumber(uint(step), b)
		}

		if ((step + 1 -2 ) % 3 == 1) && (1 == b) {
			voteNum = pleader.AddVoteNumber(uint(step), b)
		}

		if int(step) == Config().maxBBASteps + 3 {
			b = 2
			voteNum = pleader.AddVoteNumber(uint(step), b)
		}
		logger.Info("save bba message: leader", hash.String(), "step", step, "vote result", b, "vote number sum", voteNum)
		return this.endCondition(voteNum, b)
	}
	return IDLE
}

func (this *Round)receiveMsgBba(msg *BinaryByzantineAgreement) {
	//verify msg
	if msg.Credential.Round != this.round {
		logger.Warn("verify fail, bba msg is not in current round", msg.Credential.Round, this.round)
		return
	}

	if err := this.filterMsgBba(msg); err != nil {
		logger.Info("filter bba message fail:", err)
		return
	}

	step := msg.Credential.Step

	//send this msg to step other goroutine
	if stepObj, ok := this.allStepObj[int(step) + 1]; ok {
		stepObj.sendMsg(msg)
	}
	//Propagate message via p2p
	this.apos.outMsger.PropagateMsg(msg)

	//condition 0 and condition 1
	ret := this.saveMsgBba(msg)

	if ret != IDLE {
		//end condition 0, 1 or maxstep
		this.broadCastStop()
		//todo need import block to block chain

		logger.Info("OK Consensus....ret:",ret)
		this.quitCh <- 1
	}
}

func (this *Round)commonProcess() {
	for{
		select {
		// receive message
		case outData := <-this.apos.outMsger.GetDataMsg():
			switch v := outData.(type) {
			case *CredentialSign:
				this.receiveMsgCs(v)
			case *BlockProposal:
				this.receiveMsgBp(v)
			case *GradedConsensus:
				this.receiveMsgGc(v)
			case *BinaryByzantineAgreement:
				this.receiveMsgBba(v)
			default:
				logger.Warn("invalid message type ",reflect.TypeOf(v))
			}
		case <-this.quitCh:
			logger.Info("round exit ")
			return
		}
	}
}


func (this *Round)run(){
	wg := sync.WaitGroup{}
	// make verifiers Credential
	this.generateCredentials()

	// broadcast Credentials
	this.broadcastCredentials()

	this.startVerify(&wg)

	this.commonProcess()
	wg.Wait()
	this.roundOverCh<-1 //inform the caller,the mission complete
}

type Apos struct {
	systemParam interface{} //the difference of algoParam and systemParam is that algoParam show the Apos
	//running status,but the systemParam show the Mjoy runing
	mainStep int
	commonTools CommonTools
	outMsger OutMsger

	//all goroutine send msg to Apos by this Chan
	allMsgBridge chan dataPack


	roundCtx      *Round

	roundOverCh   chan interface{}
	aposStopCh    chan interface{}  //for test if apos just deal once
	stop bool
	lock sync.RWMutex
}


//Create Apos
func NewApos(msger OutMsger ,cmTools CommonTools)*Apos{
	a := new(Apos)
	a.outMsger = msger
	a.commonTools = cmTools
	a.allMsgBridge = make(chan dataPack , 10000)
	a.roundOverCh = make(chan interface{} , 1)
	a.aposStopCh = make(chan interface{} , 1)
	a.outMsger = MsgTransfer()

	a.reset()

	return a
}

func (this *Apos)makeEmptyBlockForTest(cs *CredentialSign)*block.Block{
	header := &block.Header{Number:types.NewBigInt(*big.NewInt(int64(this.commonTools.GetNowBlockNum()))),Time:types.NewBigInt(*big.NewInt(time.Now().Unix()))}
	//chainId := big.NewInt(100)
	//signer := block.NewBlockSigner(chainId)
	srcBytes := []byte{}
	srcBytes = append(srcBytes , cs.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes , cs.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes , cs.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)
	header.ConsensusData.Id = ConsensusDataId
	header.ConsensusData.Para = h
	R,S,V := this.commonTools.SIG(header.HashNoSig())
	header.R = &types.BigInt{*R}
	header.S = &types.BigInt{*S}
	header.V = &types.BigInt{*V}

	b := block.NewBlock(header , nil , nil)
	return b
}

func (this *Apos)SetOutMsger(outMsger OutMsger){
	this.outMsger = outMsger
}

//this is the main loop of Apos
func (this *Apos)Run(){

	//start round
	//this.roundOverCh<-1
	this.roundCtx = newRound(this.commonTools.GetNextRound(),this,this.roundOverCh)
	go this.roundCtx.run()
	logger.Info("Apos is running.....")
	for{
		select {
		case <-this.roundOverCh:
			logger.Info("round overs ", this.roundCtx.round)
			this.aposStopCh<-1
			return //if apos deal once ,stop it
			this.roundCtx = newRound(this.commonTools.GetNextRound(),this,this.roundOverCh)
			go this.roundCtx.run()
		}
	}
}


//reset the status of Apos
func (this *Apos)reset(){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.mainStep = 0
	this.stop = false
}


func (this *Apos)getSender(cs *CredentialSig) (types.Address, error){
	cd := CredentialData{cs.Round,cs.Step, this.commonTools.GetQr_k(1)}
	sig := &SignatureVal{&cs.R, &cs.S, &cs.V}
	return this.commonTools.Sender(cd.Hash(), sig)
}

//Create The Credential
//todo need private key for sign
func (this *Apos)makeCredential(s int) *CredentialSign{
	//create the credential and check i is the potential leader or verifier
	//r := this.commonTools.GetNowBlockNum()
	//k := this.algoParam.GetK()
	//get Qr_k
	r := this.commonTools.GetNowBlockNum()
	k := 1

	Qr_k := this.commonTools.GetQr_k(k)
	//str := fmt.Sprintf("%d%d%s",r,k,Qr_k.Hex())
	//get sig
	cd := CredentialData{*types.NewBigInt(*big.NewInt(int64(r))),*types.NewBigInt(*big.NewInt(int64(s))), Qr_k}

	//R,S,V := this.commonTools.SIG(types.BytesToHash([]byte(str)))
	R,S,V := this.commonTools.SIG(cd.Hash())

	//if endFloat <= this.algoParam
	c := new(CredentialSign)
	c.Round = uint64(r)
	c.Step = uint64(s)
	c.R = &types.BigInt{IntVal:*R}
	c.S = &types.BigInt{IntVal:*S}
	c.V = &types.BigInt{IntVal:*V}

	return c


}

func (this *Apos)StopCh()chan interface{}{
	return this.aposStopCh
}

func (this *Apos)judgeVerifier(cs *CredentialSign, setp int) bool{
	h := cs.Signature.hash()
	leader := false
	if 1 == setp {
		leader = true
	}
	return isPotVerifier(h.Bytes(), leader)
}

//before

//func (this *Apos)stepsFactory(step int , pCredential *CredentialSig)(stepObj stepInterface){
//	stepObj = nil
//	switch step {
//	case 1:
//		stepObj = makeStep1Obj(this,pCredential,step)
//	case 2:
//		stepObj = makeStep2Obj(this,pCredential,step)
//	case 3:
//		stepObj = makeStep3Obj(this,pCredential,step)
//	case 4:
//		stepObj = makeStep4Obj(this,pCredential,step)
//
//	default:
//		if step > Config().maxBBASteps + 3{
//			stepObj = nil
//		}else if (step >= 5 && step <= (Config().maxBBASteps + 2)) {
//			stepObj = makeStep567Obj(this,pCredential,step)
//		}else if (step == (Config().maxBBASteps + 3)){
//			stepObj = makeStepm3Obj(this,pCredential,step)
//		}else{
//			stepObj = nil
//		}
//	}
//	return
//}

//now
func (this *Apos)stepsFactory(step int , pCredential *CredentialSign)(stepRoutineObj *stepRoutine , stepObj step){

	stepCtx := makeStepCtx()
	//GetCredential
	stepCtx.getCredential = func() *CredentialSign {
		return pCredential
	}
	stepCtx.esig = this.commonTools.ESIG
	stepCtx.sendInner = this.outMsger.SendInner
	stepCtx.propagateMsg = this.outMsger.PropagateMsg
	stepCtx.getStep = func()int{
		return step
	}
	stepRoutineObj = nil   // if anyWrong ,return a nil obj
	stepObj = nil


	switch step {
	case 1:

		stepRoutineObj = newStepRoutine()
		stepCtx.makeEmptyBlockForTest = this.makeEmptyBlockForTest

		delayT := time.Duration(3)

		stepObj = makeStepObj1(delayT*time.Second , stepRoutineObj.stopCh)
		stepObj.setCtx(stepCtx)


	case 2:
		stepRoutineObj = newStepRoutine()

		delayT := time.Duration(Config().verifyDelay + Config().blockDelay)
		if LessTimeDelayFlag{
			delayT = time.Duration(LessTimeDelayCnt)
		}


		stepObj = makeStepObj2(delayT*time.Second , stepRoutineObj.stopCh)
		stepObj.setCtx(stepCtx)


	case 3:
		stepRoutineObj = newStepRoutine()

		delayT := time.Duration(3*Config().verifyDelay + Config().blockDelay)
		if LessTimeDelayFlag{
			delayT = time.Duration(LessTimeDelayCnt)
		}

		stepObj = makeStepObj3(delayT*time.Second , stepRoutineObj.stopCh)
		stepObj.setCtx(stepCtx)

	case 4:
		stepRoutineObj = newStepRoutine()

		delayT := time.Duration(5*Config().verifyDelay + Config().blockDelay)
		if LessTimeDelayFlag{
			delayT = time.Duration(LessTimeDelayCnt)
		}

		stepObj = makeStepObj4(delayT*time.Second , stepRoutineObj.stopCh)
		stepObj.setCtx(stepCtx)


	default:

		if step > Config().maxBBASteps + 3{
			stepRoutineObj = nil

		}else if (step >= 5 && step <= (Config().maxBBASteps + 2)) {
			stepRoutineObj = newStepRoutine()

			delayT := time.Duration((2*step -3)*Config().verifyDelay + Config().blockDelay)
			if LessTimeDelayFlag{
				delayT = time.Duration(LessTimeDelayCnt)
			}

			stepObj = makeStepObj567(delayT*time.Second , stepRoutineObj.stopCh)
			stepObj.setCtx(stepCtx)


		}else if (step == (Config().maxBBASteps + 3)){
			//stepObj = makeStepm3Obj(this,pCredential,step)
			stepRoutineObj = newStepRoutine()

			delayT := time.Duration((2*Config().maxBBASteps + 3)*Config().verifyDelay + Config().blockDelay)
			if LessTimeDelayFlag{
				delayT = time.Duration(LessTimeDelayCnt)
			}

			stepObj = makeStepObjm3(delayT*time.Second , stepRoutineObj.stopCh)
			stepObj.setCtx(stepCtx)

		}else{
			stepRoutineObj = nil
			stepObj = nil
		}
	}
	return
}














