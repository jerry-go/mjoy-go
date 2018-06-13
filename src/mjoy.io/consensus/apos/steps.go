package apos

import (
	"fmt"
	"mjoy.io/common/types"
	"time"
)

//steps handle

//step 1:Block Proposal
type step1BlockProposal struct {
	msgIn chan interface{}   //Data: Out ---- > In , we should create it
	msgOut chan interface{}  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int            //which step the obj stay
	apos *Apos

}
func makeStep1Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan interface{} , step int)*step1BlockProposal{
	s := new(step1BlockProposal)
	s.apos = pApos


	s.msgIn = make(chan interface{} , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)

	s.step = step
	return s
}
func (this *step1BlockProposal)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step1BlockProposal)stop(){
	this.exit<-1
}

func (this *step1BlockProposal)run(){

	for{
		//new a M1 data
		m1 := new(M1)
		//fill struct members
		//packdata and send out
		data := m1.GetMsgp()
		PackConsensusData(1,1,data)

		//todo:what should we dealing
		fmt.Println("For test")
	}
}

//step 2:First step of GC
type step2FirstStepGC struct {
	msgIn chan interface{}   //Data: Out ---- > In , we should create it
	msgOut chan interface{}  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	smallestLBr *M1 //this node regard smallestLBr as the smallest credential's block info
	apos *Apos

}

func makeStep2Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan interface{} , step int)*step2FirstStepGC{
	s := new(step2FirstStepGC)
	s.apos = pApos


	s.msgIn = make(chan interface{} , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step2FirstStepGC)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step2FirstStepGC)stop(){
	this.exit<-1
}

func (this *step2FirstStepGC)run(){
	//this step ,we should wait the time
	delayT := time.Duration(this.apos.algoParam.timeDelayY + this.apos.algoParam.timeDelayA)

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			//time to work now,send all

			m2 := new(M23)
			m2.Credential = this.smallestLBr.Credential
			m2.Hash = this.smallestLBr.Block.Hash()

			R,S,V := this.apos.commonTools.SIG(m2.Hash[:])
			sigVal := new(SignatureVal)
			sigVal.R = types.BigInt{*R}
			sigVal.S = types.BigInt{*S}
			sigVal.V = types.BigInt{*V}

			m2.Esig = sigVal.GetMsgp()

			data := m2.GetMsgp()
			consensusPacket := PackConsensusData(2,1,data)
			this.msgOut<-consensusPacket
			return
		case data:=<-this.msgIn:
			m1 := new(M1)
			*m1 = data.(M1)

			if m1 == nil{
				continue
			}
			if this.smallestLBr == nil {
				this.smallestLBr = m1
				continue
			}
			//compare M1 before and M1 current
			r := this.smallestLBr.Credential.Cmp(m1.Credential)
			if r > 0 {
				//exchange smallestLBr and m1
				this.smallestLBr = m1
			}

		case <-this.exit:
			return
		}
	}
}


//step 3:Second Step of GC
type step3SecondStepGC struct {
	msgIn chan interface{}   //Data: Out ---- > In , we should create it
	msgOut chan interface{}  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

	//all M2 have received
	allM2Index map[types.Hash]map[CredentialSig]bool


}

func makeStep3Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan interface{} , step int)*step3SecondStepGC{
	s := new(step3SecondStepGC)
	s.apos = pApos


	s.msgIn = make(chan interface{} , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	s.allM2Index = make(map[types.Hash]map[CredentialSig]bool)
	return s
}

func (this *step3SecondStepGC)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step3SecondStepGC)stop(){
	this.exit<-1
}

func (this *step3SecondStepGC)run(){
	//this step ,we should wait the time
	delayT := time.Duration(3*this.apos.algoParam.timeDelayY + this.apos.algoParam.timeDelayA)

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			//time to work now,send all
			total:=0
			max := 0
			maxHash := types.Hash{}
			for hash,supporter := range this.allM2Index{
				_ = hash
				currentLen := len(supporter)
				total += len(supporter)

			}

			return
		case data:=<-this.msgIn:
			m2 := new(M23)
			*m2 = data.(M23)
			if m2 == nil{
				continue
			}
			//add to IndexMap
			var subIndex map[CredentialSig]bool
			subIndex = this.allM2Index[m2.Hash]
			if subIndex == nil {
				this.allM2Index[m2.Hash] = make(map[CredentialSig]bool)
				subIndex = this.allM2Index[m2.Hash]
			}

			if _ , ok := subIndex[*m2.Credential];!ok{
				subIndex[*m2.Credential] = true
			}
			continue

		case <-this.exit:
			return
		}
	}
}


//step 4:First Step of BBA*
type step4FirstStepBBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos
	//all M3 have received
	allM2Index map[types.Hash]map[CredentialSig]bool

}

func makeStep4Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step4FirstStepBBA{
	s := new(step4FirstStepBBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step4FirstStepBBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step4FirstStepBBA)stop(){
	this.exit<-1
}

func (this *step4FirstStepBBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}


//step 5<= s <= m+2 ,s-2 mod 3 == 0 mod 3:Coin-Fixed-To-0 step of BBA*
type step5CoinFixedTo0BBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep5Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step5CoinFixedTo0BBA{
	s := new(step5CoinFixedTo0BBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step5CoinFixedTo0BBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step5CoinFixedTo0BBA)stop(){
	this.exit<-1
}

func (this *step5CoinFixedTo0BBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}

//step 6<= s <= m+2 ,s-2 mod 3 == 1 mod 3:Coin-Fixed-To-1 step of BBA*
type step6CoinFixedTo1BBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep6Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step6CoinFixedTo1BBA{
	s := new(step6CoinFixedTo1BBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step6CoinFixedTo1BBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step6CoinFixedTo1BBA)stop(){
	this.exit<-1
}

func (this *step6CoinFixedTo1BBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}

//step 7<= s <= m+2 ,s-2 mod 3 == 2 mod 3:Coin-Fixed-To-1 step of BBA*
type step7CoinGenFlipBBA struct {
	msgIn chan []byte   //Data: Out ---- > In , we should create it
	msgOut chan []byte  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos

}

func makeStep7Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan []byte , step int)*step7CoinGenFlipBBA{
	s := new(step7CoinGenFlipBBA)
	s.apos = pApos


	s.msgIn = make(chan []byte , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	return s
}

func (this *step7CoinGenFlipBBA)sendMsg(data []byte)error{
	//todo:
	return nil
}

func (this *step7CoinGenFlipBBA)stop(){
	this.exit<-1
}

func (this *step7CoinGenFlipBBA)run(){
	for{
		//todo:what should we dealing
		fmt.Println("For test")
	}
}














