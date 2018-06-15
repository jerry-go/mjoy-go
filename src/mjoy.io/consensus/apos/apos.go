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
	"fmt"
	"reflect"
	"errors"
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

type Apos struct {
	algoParam *algoParam
	systemParam interface{} //the difference of algoParam and systemParam is that algoParam show the Apos
							//running status,but the systemParam show the Mjoy runing
	mainStep int
	commonTools CommonTools
	outMsger OutMsger

	//all goroutine send msg to Apos by this Chan
	allMsgBridge chan dataPack


	roundCtx      *Round
	validate      *MsgValidator
	roundOverCh   chan interface{}

	stop bool
	lock sync.RWMutex
}

type diffCredentialSig struct {
	difficulty    *big.Int
	credential    *CredentialSig
}

type VoteInfo struct {
	sum           int
}

// Potential Leader used for judge End condition 0 and 1
type PotentialLeader struct {
	m1            *M1
	diff          *big.Int
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

type peerMsgs struct {
	msgCommons     map[int]*MCommon
	msg23s         map[int]*M23

	//0 :default honesty peer. 1: malicious peer
	honesty        uint
}

//round context
type Round struct {
	round          types.BigInt
	//condition 0 and condition 1 end number tH = 2n/3 + 1
	targetNum      int

	apos           *Apos

	credentials    map[int]*CredentialSig

	allStepObj     map[int]stepInterface

	smallestLBr    *M1
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
	this.round = types.BigInt{*big.NewInt(int64(round))}
	this.apos = apos
	this.roundOverCh = roundOverCh

	this.targetNum = this.apos.algoParam.targetNum
	this.maxLeaderNum = this.apos.algoParam.maxLeaderNum
	this.credentials = make(map[int]*CredentialSig)
	this.allStepObj = make(map[int]stepInterface)
	this.leaders = make(map[types.Hash]*PotentialLeader)
	this.quitCh = make(chan interface{} , 1)
	this.msgs = make(map[types.Address]*peerMsgs)
}

func (this *Round)setSmallestBrM1(m *M1){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.smallestLBr = m
}

func (this *Round)addStepObj(step int , stepObj stepInterface){
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
func (this *Round)GenerateCredentials() {
	for i := 1; i < this.apos.algoParam.maxSteps; i++{
		credential := this.apos.makeCredential(i)
		isVerfier := this.apos.judgeVerifier(credential, i)
		if isVerfier {
			this.credentials[i] = credential
		}
	}
}

func (this *Round)BroadcastCredentials() {
	for i, credential := range this.credentials {
		logger.Info("SendCredential round", this.round.IntVal.Uint64(), "step", i)
		this.apos.outMsger.SendCredential(credential)
	}
}


func (this *Round)StartVerify(wg *sync.WaitGroup) {
	for i, credential := range this.credentials {
		stepObj := this.apos.stepsFactory(i, credential)
		if stepObj != nil {
			this.addStepObj(i, stepObj)
			go stepObj.run(wg)
		}
	}
}


// hear M0 is the Credential message
func (this *Round)ReceiveM0(msg *CredentialSig) {
	//verify msg
	err := this.apos.validate.ValidateCredential(msg)
	if err != nil {
		logger.Info("verify m0 fail", err)
		return
	}
	//todo msg filter need

	//Propagate message via p2p
	this.apos.outMsger.PropagateCredential(msg)
}

func (this *Round)MinDifficultM1() *PotentialLeader{
	curDiff := maxUint256
	cur := &PotentialLeader{diff:maxUint256}
	for _,pleader := range this.leaders {
		if curDiff.Cmp(pleader.diff) > 0 {
			curDiff = pleader.diff
			cur = pleader
		}
	}
	return cur
}


func (this *Round)SaveM1(msg *M1) {
	difficulty := GetDifficulty(msg.Credential)
	if this.curLeaderNum >= this.maxLeaderNum {
		if difficulty.Cmp(this.curLeaderDiff) <= 0 {
			logger.Info("can not save m1 because of difficulty is not catch", difficulty)
			return
		}
		delete(this.leaders, this.curLeader)
		pleader := this.MinDifficultM1()
		this.curLeaderDiff = pleader.diff
		this.curLeader = pleader.m1.Block.Hash()
		this.curLeaderNum--
	}

	if this.curLeaderNum == 0 || difficulty.Cmp(this.curLeaderDiff) < 0 {
		this.curLeaderDiff = difficulty
		this.curLeader = msg.Block.Hash()
	}

	hash := msg.Block.Hash()
	if _, ok := this.leaders[hash]; !ok {
		pleader := &PotentialLeader{msg, difficulty,make(map[uint]*VoteInfo)}
		this.leaders[hash] = pleader
		this.curLeaderNum++
	}
}

func (this *Round)ReceiveM1(msg *M1) {
	//verify msg
	if msg.Credential.Round.IntVal.Cmp(&this.round.IntVal) != 0 {
		logger.Warn("verify fail, M1 msg is not in current round", msg.Credential.Round.IntVal.Uint64(), this.round.IntVal.Uint64())
		return
	}

	if msg.Credential.Step.IntVal.Uint64() != 1 {
		logger.Warn("verify fail, M1 msg step is not 1", msg.Credential.Round.IntVal.Uint64(), msg.Credential.Step.IntVal.Uint64())
		return
	}

	err := this.apos.validate.ValidateM1(msg)
	if err != nil {
		logger.Info("verify m1 fail", err)
		return
	}
	//todo msg filter need

	//send this msg to step2 goroutine
	if stepObj, ok := this.allStepObj[2]; ok {
		stepObj.sendMsg(msg, this)
	}

	// for M1 Propagate process will in stepObj

	this.SaveM1(msg)
}


func (this *Round)filterM23(msg *M23) error {
	address, err := this.apos.getSender(msg.Credential)
	if err != nil {
		return err
	}
	step := msg.Credential.Step.IntVal.Uint64()

	if peerM23s, ok := this.msgs[address]; ok {
		if peerM23s.honesty == 1 {
			return errors.New("not honesty peer")
		}

		if m23, ok := peerM23s.msg23s[int(step)]; ok {
			if m23.Hash == msg.Hash {
				return errors.New("duplicate message m23")
			} else {
				peerM23s.honesty = 1
				return errors.New("receive different  vote message m23, it must a malicious peer")
			}
		} else {
			peerM23s.msg23s[int(step)] = msg
		}
	} else {
		ps := &peerMsgs{
			msgCommons: make(map[int]*MCommon),
			msg23s: make(map[int]*M23),
			honesty: 0,
		}
		ps.msg23s[int(step)] = msg
		this.msgs[address] = ps
	}
	return nil
}

func (this *Round)ReceiveM23(msg *M23) {
	//verify msg
	if msg.Credential.Round.IntVal.Cmp(&this.round.IntVal) != 0 {
		logger.Warn("verify fail, M23 msg is not in current round", msg.Credential.Round.IntVal.Uint64(), this.round.IntVal.Uint64())
		return
	}

	step := msg.Credential.Step.IntVal.Uint64()
	if step != 2  && step != 3 {
		logger.Warn("verify fail, M23 msg step is not 2 or 3", msg.Credential.Round.IntVal.Uint64(), step)
		return
	}
	err := this.apos.validate.ValidateM23(msg)
	if err != nil {
		logger.Info("verify m23 fail", err)
		return
	}

	if err = this.filterM23(msg); err != nil {
		logger.Info("filter m23 fail", err)
		return
	}

	//send this msg to step3 or step4 goroutine
	if stepObj, ok := this.allStepObj[int(step) + 1]; ok {
		stepObj.sendMsg(msg, this)
	}
}

func (this *Round)EndCondition(voteNum int, b uint) int {
	if voteNum >= this.targetNum {
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

func (this *Round)SaveMCommon(msg *MCommon) int {
	hash := msg.Hash
	if pleader, ok := this.leaders[hash]; ok {
		step := msg.Credential.Step.IntVal.Uint64()
		b := msg.B
		voteNum := 0
		if ((step + 1 -2 ) % 3 == 0) && (0 == b) {
			voteNum = pleader.AddVoteNumber(uint(step), b)
		}

		if ((step + 1 -2 ) % 3 == 1) && (1 == b) {
			voteNum = pleader.AddVoteNumber(uint(step), b)
		}

		if int(step) == this.apos.algoParam.maxSteps + 3 {
			b = 2
			voteNum = pleader.AddVoteNumber(uint(step), b)
		}
		return this.EndCondition(voteNum, b)
	}
	return IDLE
}

func (this *Round)ReceiveMCommon(msg *MCommon) {
	//verify msg
	if msg.Credential.Round.IntVal.Cmp(&this.round.IntVal) != 0 {
		logger.Warn("verify fail, MCommon msg is not in current round", msg.Credential.Round.IntVal.Uint64(), this.round.IntVal.Uint64())
		return
	}

	step := msg.Credential.Step.IntVal.Uint64()
	if step < 4{
		logger.Warn("verify fail, MCommon msg step is not right", msg.Credential.Round.IntVal.Uint64(), step)
		return
	}
	err := this.apos.validate.ValidateMCommon(msg)
	if err != nil {
		logger.Info("verify msg common fail", err)
		return
	}
	//todo msg filter need

	//send this msg to step other goroutine
	if stepObj, ok := this.allStepObj[int(step) + 1]; ok {
		stepObj.sendMsg(msg, this)
	}

	//condition 0 and condition 1
	ret := this.SaveMCommon(msg)

	if ret != IDLE {
		//end condition 0, 1 or maxstep
		this.broadCastStop()
		//todo need import block to block chain


		this.quitCh <- 1
	}
}

func (this *Round)CommonProcess() {
	for{
		select {
		// receive message
		case outData := <-this.apos.outMsger.GetDataMsg():
			switch v := outData.(type) {
			case *CredentialSig:
				//fmt.Println(v)
				this.ReceiveM0(v)
			case *M1:
				//fmt.Println(v)
				this.ReceiveM1(v)
			case *M23:
				//fmt.Println(v)
				this.ReceiveM23(v)
			case *MCommon:
				//fmt.Println(v)
				this.ReceiveMCommon(v)
			default:
				logger.Warn("invalid message type ",reflect.TypeOf(v))
			}
		case <-this.quitCh:
			logger.Info("round exit ")
			return
		}
	}
}


func (this *Round)Run(){
	wg := sync.WaitGroup{}
	// make verifiers Credential
	this.GenerateCredentials()

	// broadcast Credentials
	this.BroadcastCredentials()

	this.StartVerify(&wg)

	this.CommonProcess()
	wg.Wait()
	this.roundOverCh<-1 //inform the caller,the mission complete
}


//Create Apos
func NewApos(msger OutMsger ,cmTools CommonTools)*Apos{
	a := new(Apos)
	a.outMsger = msger
	a.commonTools = cmTools
	a.allMsgBridge = make(chan dataPack , 10000)

	a.reset()

	a.validate = NewMsgValidator(a.algoParam, a)

	return a
}



//this is the main loop of Apos
func (this *Apos)Run(){

	//start round
	this.roundOverCh<-1
	for{
		select {
		case <-this.roundOverCh:
			this.roundCtx = newRound(this.commonTools.GetNextRound(),this,this.roundOverCh)
			go this.roundCtx.Run()
		}
	}
}


//reset the status of Apos
func (this *Apos)reset(){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.algoParam = newAlgoParam()
	this.mainStep = 0
	this.stop = false
}


func (this *Apos)getSender(cs *CredentialSig) (types.Address, error){
	cd := CredentialData{cs.Round,cs.Step, this.commonTools.GetQr_k(1)}
	sig := &SignatureVal{&cs.R, &cs.S, &cs.V}
	return this.commonTools.Sender(cd.Hash(), sig)
}

//Create The Credential
func (this *Apos)makeCredential(s int) *CredentialSig{
	//create the credential and check i is the potential leader or verifier
	//r := this.commonTools.GetNowBlockNum()
	//k := this.algoParam.GetK()
	//get Qr_k
	r := this.commonTools.GetNowBlockNum()
	k := 1

	Qr_k := this.commonTools.GetQr_k(k)
	str := fmt.Sprintf("%d%d%s",r,k,Qr_k.Hex())
	//get sig
	R,S,V := this.commonTools.SIG([]byte(str))

	//if endFloat <= this.algoParam
	c := new(CredentialSig)
	c.Round = types.BigInt{IntVal:*big.NewInt(int64(r))}
	c.Step = types.BigInt{IntVal:*big.NewInt(int64(s))}
	c.R = types.BigInt{IntVal:*R}
	c.S = types.BigInt{IntVal:*S}
	c.V = types.BigInt{IntVal:*V}

	return c


}


func (this *Apos)judgeVerifier(pCredentialSig *CredentialSig, setp int) bool{
	srcBytes := []byte{}
	srcBytes = append(srcBytes , pCredentialSig.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)
	difficulty := BytesToDifficulty(h)

	verifierDifficulty := new(big.Int)
	if 1 == setp {
		verifierDifficulty = this.algoParam.leaderDifficulty
	} else {
		verifierDifficulty = this.algoParam.verifierDifficulty
	}
	if difficulty.Cmp(verifierDifficulty) > 0 {
		return true
	}

	return false
}



func (this *Apos)stepsFactory(step int , pCredential *CredentialSig)(stepObj stepInterface){
	stepObj = nil
	switch step {
	case 1:
		stepObj = makeStep1Obj(this,pCredential,step)
	case 2:
		stepObj = makeStep2Obj(this,pCredential,step)
	case 3:
		stepObj = makeStep3Obj(this,pCredential,step)
	case 4:
		stepObj = makeStep4Obj(this,pCredential,step)

	default:
		if step > this.algoParam.maxSteps{
			stepObj = nil
		}else if (step >= 5 && step <= (this.algoParam.m + 2)) {
			stepObj = makeStep567Obj(this,pCredential,step)
		}else if (step == (this.algoParam.m+3)){
			stepObj = makeStepm3Obj(this,pCredential,step)
		}else{
			stepObj = nil
		}
	}
	return
}

















