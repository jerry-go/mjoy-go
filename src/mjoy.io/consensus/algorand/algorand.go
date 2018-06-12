package algorand

import (
	"sync"
	"mjoy.io/utils/crypto"
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
	mu sync.RWMutex
}

//Create Algorand
func NewAlgorand(msger OutMsger ,cmTools CommonTools )*Algorand{
	a := new(Algorand)
	a.all = make(map[int]stepInterface)
	a.outMsger = msger
	a.commonTools = cmTools

	a.reset()

	return a
}

//this is the main loop of Algorand
func Run(){

}

//reset the status of algorand
func (this *Algorand)reset(){
	this.all = make(map[int]stepInterface)
	this.algoParam = newAlgoParam()
	this.mainStep = 0
}


//Create The Credential
func (this *Algorand)makeCredential(){
	//create the credential and check i is the potential leader or verifier
	r := this.commonTools.GetNowBlockNum()
	k := this.algoParam.GetK()
	//get Qr_k
	Qr_k := this.commonTools.GetQr_k(k)
	//get sig
	R,S,V := this.commonTools.SIG(r,1,Qr_k)
	srcBytes := []byte{}
	srcBytes = append(srcBytes , R.Bytes()...)
	srcBytes = append(srcBytes , S.Bytes()...)
	srcBytes = append(srcBytes , V.Bytes()...)

	h := crypto.Keccak256(srcBytes)


	endFloat , err := BytesToFloat(h)
	if err != nil {
		panic(err)
	}

	_ = endFloat
	//if endFloat <= this.algoParam

}



















