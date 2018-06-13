package algorand

import (
	"sync"
	"mjoy.io/utils/crypto"
	"mjoy.io/common/types"
	"math/big"
	"github.com/tinylib/msgp/msgp"
	"bytes"
)

/*
Instructions For Algorand:
The Type-Algorand manage the main loop of algorand consensus,
handle the Condition0,Condition1 and m+3,
and Output Algorand-SystemParam to the sub-goroutine
*/

type Algorand struct {
	all map[int]stepInterface
	algoParam *algoParam
	systemParam interface{} //the difference of algoParam and systemParam is that algoParam show the Algorand
							//running status,but the systemParam show the Mjoy runing
	mainStep int
	commonTools CommonTools
	outMsger OutMsger

	//all goroutine send msg to Algorand by this Chan
	allMsgBridge chan []byte

	mu sync.RWMutex
}

//Create Algorand
func NewAlgorand(msger OutMsger ,cmTools CommonTools )*Algorand{
	a := new(Algorand)
	a.all = make(map[int]stepInterface)
	a.outMsger = msger
	a.commonTools = cmTools
	a.allMsgBridge = make(chan []byte , 10000)
	a.reset()

	return a
}

//this is the main loop of Algorand
func (this *Algorand)Run(){

	for{
		//Algorand status reset
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
				broadCastCredential = this.judgeVerifier(pCredential)
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






		}

	}
}

//reset the status of algorand
func (this *Algorand)reset(){
	this.all = make(map[int]stepInterface)
	this.algoParam = newAlgoParam()
	this.mainStep = 0
}


//Create The Credential
func (this *Algorand)makeCredential(s int)*CredentialSig{
	//create the credential and check i is the potential leader or verifier
	//r := this.commonTools.GetNowBlockNum()
	//k := this.algoParam.GetK()
	//get Qr_k
	r := this.commonTools.GetNowBlockNum()
	k := 1

	Qr_k := this.commonTools.GetQr_k(k)
	//get sig
	R,S,V := this.commonTools.SIG(r,1,Qr_k)

	//if endFloat <= this.algoParam
	c := new(CredentialSig)
	c.Round = types.BigInt{IntVal:*big.NewInt(int64(r))}
	c.Step = types.BigInt{IntVal:*big.NewInt(int64(s))}
	c.R = types.BigInt{IntVal:*R}
	c.S = types.BigInt{IntVal:*S}
	c.V = types.BigInt{IntVal:*V}

	return c


}

func (this *Algorand)judgeLeader(pCredentialSig *CredentialSig)bool{
	srcBytes := []byte{}
	srcBytes = append(srcBytes , pCredentialSig.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)
	endFloat , err := BytesToFloat(h)
	if err != nil {
		logger.Errorf("judgeLeader:%s" , err.Error())
		return false
	}

	if endFloat <= this.algoParam.pLeader {

		return true
	}

	return false
}

func (this *Algorand)judgeVerifier(pCredentialSig *CredentialSig)bool{
	srcBytes := []byte{}
	srcBytes = append(srcBytes , pCredentialSig.R.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.S.IntVal.Bytes()...)
	srcBytes = append(srcBytes , pCredentialSig.V.IntVal.Bytes()...)

	h := crypto.Keccak256(srcBytes)
	endFloat , err := BytesToFloat(h)
	if err != nil {
		logger.Errorf("judgeLeader:%s" , err.Error())
		return false
	}
	if endFloat <= this.algoParam.pVerifier {

		return true
	}

	return false
}


func (this *Algorand)stepsFactory(step int , pCredential *CredentialSig)(stepObj stepInterface){
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
		}else if (step >= 5 && step <= (this.algoParam.m + 2) )&&((step - 2)%3 == 0%3) {
			stepObj = makeStep5Obj(this,pCredential,this.allMsgBridge,step)
		}else if (step >= 6 && step <= (this.algoParam.m + 2) )&&((step - 2)%3 == 1%3) {
			stepObj = makeStep6Obj(this,pCredential,this.allMsgBridge,step)
		}else if (step >= 7 && step <= (this.algoParam.m + 2) )&&((step - 2)%3 == 2%3) {
			stepObj = makeStep7Obj(this,pCredential,this.allMsgBridge,step)
		}else{
			stepObj = nil
		}


	}
	return
}

















