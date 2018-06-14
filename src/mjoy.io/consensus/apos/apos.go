package apos

import (
	"sync"
	"mjoy.io/utils/crypto"
	"mjoy.io/common/types"
	"math/big"
	"github.com/tinylib/msgp/msgp"
	"bytes"
	"fmt"
	"reflect"
)

/*
Instructions For Apos:
The Type-Apos manage the main loop of Apos consensus,
handle the Condition0,Condition1 and m+3,
and Output Apos-SystemParam to the sub-goroutine
*/

type Apos struct {
	allStepObj map[int]stepInterface
	algoParam *algoParam
	systemParam interface{} //the difference of algoParam and systemParam is that algoParam show the Apos
							//running status,but the systemParam show the Mjoy runing
	mainStep int
	commonTools CommonTools
	outMsger OutMsger

	//all goroutine send msg to Apos by this Chan
	allMsgBridge chan dataPack

	roundCtx      *Round

	stop bool
	lock sync.RWMutex
}

type diffCredentialSig struct {
	difficulty    *big.Int
	credential    *CredentialSig
}

//round context
type Round struct {
	round          types.BigInt
	apos           *Apos

	credentials    map[int]*CredentialSig

	allStepObj     map[int]stepInterface

	smallestLBr     *M1
	lock sync.RWMutex
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


func (this *Round)StartVerify() {
	for i, credential := range this.credentials {
		stepObj := this.apos.stepsFactory(i, credential)
		this.addStepObj(i, stepObj)
		go stepObj.run()
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
	//todo more verify need

	//send this msg to step1 goroutine
	if stepObj, ok := this.allStepObj[2]; ok {
		stepObj.sendMsg(msg, this)
	}

}

func (this *Round)ReceiveM23(msg *M23) {
	//verify msg
	if msg.Credential.Round.IntVal.Cmp(&this.round.IntVal) != 0 {
		logger.Warn("verify fail, M23 msg is not in current round", msg.Credential.Round.IntVal.Uint64(), this.round.IntVal.Uint64())
		return
	}

	if msg.Credential.Step.IntVal.Uint64() != 2  && msg.Credential.Step.IntVal.Uint64() != 3 {
		logger.Warn("verify fail, M23 msg step is not 2 or 3", msg.Credential.Round.IntVal.Uint64(), msg.Credential.Step.IntVal.Uint64())
		return
	}
	//todo more verify need

	//send this msg to step3 or step4 goroutine
	if msg.Credential.Step.IntVal.Uint64() == 2 {
		if stepObj, ok := this.allStepObj[3]; ok {
			stepObj.sendMsg(msg, this)
		}
	} else if msg.Credential.Step.IntVal.Uint64() == 3 {
		if stepObj, ok := this.allStepObj[4]; ok {
			stepObj.sendMsg(msg, this)
		}
	}
}

func (this *Round)ReceiveMCommon(msg *MCommon) {
	//verify msg
	if msg.Credential.Round.IntVal.Cmp(&this.round.IntVal) != 0 {
		logger.Warn("verify fail, MCommon msg is not in current round", msg.Credential.Round.IntVal.Uint64(), this.round.IntVal.Uint64())
		return
	}

	step := msg.Credential.Step.IntVal.Uint64()
	if step < 4{
		logger.Warn("verify fail, MCommon msg step is not right", msg.Credential.Round.IntVal.Uint64(), msg.Credential.Step.IntVal.Uint64())
		return
	}
	//todo more verify need

	//send this msg to step other goroutine
	if stepObj, ok := this.allStepObj[int(step) + 1]; ok {
		stepObj.sendMsg(msg, this)
	}

	//condition 0 and condition 1

}

func (this *Round)CommonProcess() {
	for{
		select {
		// receive message
		case outData := <-this.apos.outMsger.GetDataMsg():
			switch v := outData.(type) {
			case *CredentialSig:
				fmt.Println(v)
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
		}
	}
}

func (this *Round)Run(){
	// make verifiers Credential
	this.GenerateCredentials()

	// broadcast Credentials
	this.BroadcastCredentials()

	this.StartVerify()

}



//Create Apos
func NewApos(msger OutMsger ,cmTools CommonTools )*Apos{
	a := new(Apos)
	a.allStepObj = make(map[int]stepInterface)
	a.outMsger = msger
	a.commonTools = cmTools
	a.allMsgBridge = make(chan dataPack , 10000)
	a.reset()

	return a
}


func (this *Apos)msgRcv(){
	for{
		select {
		case outData:=<-this.outMsger.GetMsg():
			if this.stop{   //other logic deal has stop
				continue
			}
			_ = outData
			//todo:add msg to buffer
        case stepData := <-this.allMsgBridge:
        	if this.stop{   //other logic deal has stop
        		continue
	        }
        	_= stepData
        	//todo:add msg to buffer
		}
	}
}

//this is the main loop of Apos
func (this *Apos)Run(){
	go this.msgRcv()
	for{
		//Apos status reset
		this.reset()

		//1.create credential
		//2.send credential
		//3.create goroutine
		for i:=1;i<this.algoParam.maxSteps;i++{
			pCredential := this.makeCredential(i)
			broadCastCredential := false
			if i == 1 {
				broadCastCredential = this.judgeLeader(pCredential)
			}else{
				broadCastCredential = this.judgeVerifier(pCredential, i)
			}

			if broadCastCredential {
				var buf bytes.Buffer
				err := msgp.Encode(&buf , pCredential)
				if err != nil {
					logger.Errorf("broadCastCredential msgp:" , err.Error())
				}
				packData := PackConsensusData(i , 0 , buf.Bytes())
				if packData != nil {
					this.outMsger.SendMsg(packData)
				}


			}


			//here should make goroutine and run it
			exist := this.existStepObj(i)
			if !exist{
				stepObj := this.stepsFactory(i,pCredential)
				this.addStepObj(i,stepObj)
				go stepObj.run()
			}

		}

		//Condition 0 ,Condition 1 and m+3 judge


	}
}
//inform stepObj to stop running
func (this *Apos)broadCastStop(){
	for _,v := range this.allStepObj {
		v.stop()
	}
}

//reset the status of Apos
func (this *Apos)reset(){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.broadCastStop()
	this.allStepObj = make(map[int]stepInterface)
	this.algoParam = newAlgoParam()
	this.mainStep = 0
	this.stop = false
}
func (this *Apos)addStepObj(step int , stepObj stepInterface){
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.allStepObj[step]; !ok {
		this.allStepObj[step] = stepObj
	}
}

func (this *Apos)existStepObj(step int)bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if _, ok := this.allStepObj[step]; ok {
		return true
	}
	return false
}



//Create The Credential
func (this *Apos)makeCredential(s int)*CredentialSig{
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

func (this *Apos)judgeLeader(pCredentialSig *CredentialSig)bool{
	srcBytes := []byte{}
	srcBytes = append(srcBytes , pCredentialSig.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)
	difficulty := BytesToDifficulty(h)

	if difficulty.Cmp(this.algoParam.leaderDifficulty) > 0 {
		return true
	}

	return false
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
		stepObj = makeStep1Obj(this,pCredential,this.allMsgBridge,step)
	case 2:
		stepObj = makeStep2Obj(this,pCredential,this.allMsgBridge,step)
	case 3:
		stepObj = makeStep3Obj(this,pCredential,this.allMsgBridge,step)
	case 4:
		stepObj = makeStep4Obj(this,pCredential,this.allMsgBridge,step)
	default:
		if step > this.algoParam.maxSteps{
			stepObj = nil
		}else if (step >= 5 && step <= (this.algoParam.m + 2)) {
			stepObj = makeStep567Obj(this,pCredential,this.allMsgBridge,step)
		}else{
			stepObj = nil
		}
	}
	return
}

















