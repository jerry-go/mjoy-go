package apos

import (
	"fmt"
	"mjoy.io/common/types"
	"time"
	"sync"
	"container/heap"
)

//steps handle

//step 1:Block Proposal
type step1BlockProposal struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	msgOut chan dataPack  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int            //which step the obj stay
	apos *Apos
	round *Round
	pCredential *CredentialSig
	lock sync.RWMutex

}
func makeStep1Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan dataPack , step int)*step1BlockProposal{
	s := new(step1BlockProposal)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.pCredential = pCredential
	s.step = step
	return s
}
func (this *step1BlockProposal)sendMsg(data dataPack, pRound *Round)error{
	this.lock.Lock()
	defer this.lock.Unlock()
	this.round = pRound
	this.msgIn <- data
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

		this.msgOut <- m1
		//todo:what should we dealing
		fmt.Println("For test")
	}
}

//step 2:First step of GC
type step2FirstStepGC struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	msgOut chan dataPack  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	smallestLBr *M1 //this node regard smallestLBr as the smallest credential's block info
	round *Round
	apos *Apos
	pCredential *CredentialSig
	lock sync.RWMutex

}

func makeStep2Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan dataPack , step int)*step2FirstStepGC{
	s := new(step2FirstStepGC)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	s.pCredential = pCredential
	return s
}

func (this *step2FirstStepGC)sendMsg(data dataPack , pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.round = pRound
	this.msgIn <- data
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
			m2.Credential = this.pCredential
			m2.Hash = this.smallestLBr.Block.Hash()

			R,S,V := this.apos.commonTools.SIG(m2.Hash[:])
			sigVal := new(SignatureVal)
			sigVal.R = types.BigInt{*R}
			sigVal.S = types.BigInt{*S}
			sigVal.V = types.BigInt{*V}

			m2.Esig = sigVal.GetMsgp()


			this.msgOut<-m2
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
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	msgOut chan dataPack  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos
	round *Round
	pCredential *CredentialSig
	//all M2 have received
	allM2Index map[types.Hash]map[CredentialSig]bool
	lock sync.RWMutex

}

func makeStep3Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan dataPack , step int)*step3SecondStepGC{
	s := new(step3SecondStepGC)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step
	s.allM2Index = make(map[types.Hash]map[CredentialSig]bool)
	return s
}

func (this *step3SecondStepGC)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.round = pRound
	this.msgIn <- data
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
			func(this *step3SecondStepGC){
				this.lock.Lock()
				defer this.lock.Unlock()


				//time to work now,send all
				total:=0
				maxLen := 0
				maxHash := types.Hash{}


				for hash,supporter := range this.allM2Index{

					currentLen := len(supporter)
					if currentLen > maxLen{
						maxLen = currentLen
						maxHash = hash
					}
					total += currentLen

				}
				v := types.Hash{}

				if maxLen >= (2*total/3){
					v = maxHash
				}
				//pack m3 Data

				m3 := new(M23)
				m3.Credential = this.pCredential
				m3.Hash = v

				R,S,V := this.apos.commonTools.SIG(m3.Hash[:])
				sigVal := new(SignatureVal)
				sigVal.R = types.BigInt{*R}
				sigVal.S = types.BigInt{*S}
				sigVal.V = types.BigInt{*V}

				m3.Esig = sigVal.GetMsgp()

				this.msgOut<-m3

			}(this)

			return
		case data:=<-this.msgIn:
			func(this *step3SecondStepGC){
				this.lock.Lock()
				defer this.lock.Unlock()

				m2 := new(M23)
				*m2 = data.(M23)
				if m2 == nil{
					return
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

			}(this)
			continue

		case <-this.exit:
			return
		}
	}
}


//step 4:First Step of BBA*
type step4FirstStepBBA struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	msgOut chan dataPack  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	apos *Apos
	round *Round
	pCredential *CredentialSig
	//all M3 have received
	allM2Index map[types.Hash]map[CredentialSig]bool
	lock sync.RWMutex
}

func makeStep4Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan dataPack , step int)*step4FirstStepBBA{
	s := new(step4FirstStepBBA)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step

	s.pCredential = pCredential
	return s
}

func (this *step4FirstStepBBA)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.round = pRound
	this.msgIn <- data
	return nil
}

func (this *step4FirstStepBBA)stop(){
	this.exit<-1
}

func (this *step4FirstStepBBA)run(){
	//this step ,we should wait the time
	delayT := time.Duration(3*this.apos.algoParam.timeDelayY + this.apos.algoParam.timeDelayA)

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			func(this *step4FirstStepBBA){
				this.lock.Lock()
				defer this.lock.Unlock()


				//time to work now,send all
				total:=0
				maxLen := 0
				maxHash := types.Hash{}


				for hash,supporter := range this.allM2Index{

					currentLen := len(supporter)
					if currentLen > maxLen{
						maxLen = currentLen
						maxHash = hash
					}
					total += currentLen

				}
				v := types.Hash{}
				g := 0
				if maxLen >= (2*total/3){
					v = maxHash
					g = 2
				}else if maxLen >= (2*total/3){
					v = maxHash
					g = 1
				}else{
					v = types.Hash{}
					g = 0
				}

				b := 0

				if g == 2 {
					b = 0
				}else{
					b = 1
				}

				//pack m4 Data
				m4 := new(MCommon)
				m4.Hash = v
				m4.B = uint(b)
				m4.Credential = this.pCredential
				//b
				str := fmt.Sprintf("%d",m4.B)
				R,S,V := this.apos.commonTools.SIG([]byte(str))
				sigVal := new(SignatureVal)
				sigVal.R = types.BigInt{*R}
				sigVal.S = types.BigInt{*S}
				sigVal.V = types.BigInt{*V}
				m4.EsigB = sigVal.GetMsgp()
				//v
				str = fmt.Sprintf("%s",m4.Hash.Hex())
				R,S,V = this.apos.commonTools.SIG([]byte(str))
				sigVal.R = types.BigInt{*R}
				sigVal.S = types.BigInt{*S}
				sigVal.V = types.BigInt{*V}
				m4.EsigV = sigVal.GetMsgp()

				this.msgOut<-m4

			}(this)

			return
		case data:=<-this.msgIn:
			func(this *step4FirstStepBBA){
				this.lock.Lock()
				defer this.lock.Unlock()

				m3 := new(M23)
				*m3 = data.(M23)
				if m3 == nil{
					return
				}
				//add to IndexMap
				var subIndex map[CredentialSig]bool
				subIndex = this.allM2Index[m3.Hash]
				if subIndex == nil {
					this.allM2Index[m3.Hash] = make(map[CredentialSig]bool)
					subIndex = this.allM2Index[m3.Hash]
				}

				if _ , ok := subIndex[*m3.Credential];!ok{
					subIndex[*m3.Credential] = true
				}

			}(this)
			continue

		case <-this.exit:
			return
		}
	}
}

//merge 5 6 7

//step 5 6 7:Coin-Fixed-To-x step of BBA*
type step567CoinGenFlipBBA struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	msgOut chan dataPack  //Data: In ----- > Out, out caller should give it us
	exit chan int
	step int
	stepIndex int
	apos *Apos
	round *Round
	pCredential *CredentialSig
	lock sync.RWMutex

	//all M6 have received
	allMxIndex map[types.Hash]*binaryStatus
}

func makeStep567Obj(pApos *Apos , pCredential *CredentialSig , outMsgChan chan dataPack , step int)*step567CoinGenFlipBBA{
	s := new(step567CoinGenFlipBBA)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)
	s.msgOut = outMsgChan

	s.exit = make(chan int , 1)
	s.step = step


	if ((step - 2)%3 == 0%3) {
		s.stepIndex = 5
	}else if  ((step - 2)%3 == 1%3) {
		s.stepIndex = 6
	}else if  ((step - 2)%3 == 2%3) {
		s.stepIndex = 7
	}

	s.pCredential = pCredential
	return s
}

func (this *step567CoinGenFlipBBA)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.round = pRound
	this.msgIn <- data
	return nil
}

func (this *step567CoinGenFlipBBA)stop(){
	this.exit<-1
}

func (this *step567CoinGenFlipBBA)run(){
	//this step ,we should wait the time
	delayT := time.Duration(3*this.apos.algoParam.timeDelayY + this.apos.algoParam.timeDelayA)

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			func(this *step567CoinGenFlipBBA){
				this.lock.Lock()
				defer this.lock.Unlock()


				//time to work now,send all
				total:=0
				maxLen := 0
				maxHash := types.Hash{}

				max1Len := 0
				max0Len := 0


				for hash,bStatus := range this.allMxIndex{
					//add total cnt
					currentTotalLen := bStatus.getTotalCnt()

					if currentTotalLen > maxLen{
						maxLen = currentTotalLen
						maxHash = hash
					}

					total += currentTotalLen
				}

				maxBStatus := this.allMxIndex[maxHash]
				if maxBStatus == nil {
					return
				}
				max0Len = maxBStatus.getCnt(0)
				max1Len = maxBStatus.getCnt(1)




				mx := new(MCommon)
				//check 2/3 0 and 2/3 1
				if max0Len >= (2*total/3) {
					mx.Hash = maxHash
					mx.Credential = this.pCredential
					mx.B = 0
				}else if max1Len >= (2*total/3){
					mx.Hash = maxHash
					mx.Credential = this.pCredential
					mx.B = 1
				}else {
					mx.Hash = maxHash
					mx.Credential = this.pCredential
					switch this.stepIndex {
					case 5:
						mx.B = 0
					case 6:
						mx.B = 1
					case 7:
						{
							cHeap := new(CredentialSigStatusHeap)
							for _,bStatus := range this.allMxIndex{
								allCredential := bStatus.export0Credential()
								// 0 credential
								for _,c := range allCredential{
									*cHeap = append(*cHeap , &CredentialSigStatus{c:c , v:0})
								}
								allCredential = bStatus.export1Credential()
								// 1 credential
								for _,c := range allCredential{
									*cHeap = append(*cHeap , &CredentialSigStatus{c:c , v:1})
								}

							}
							heap.Init(cHeap)
							little := heap.Pop(cHeap).(CredentialSigStatus)

							mx.B = uint(little.v)
						}
					}
				}

				str := fmt.Sprintf("%d",mx.B)
				R,S,V := this.apos.commonTools.SIG([]byte(str))
				sigVal := new(SignatureVal)
				sigVal.R = types.BigInt{*R}
				sigVal.S = types.BigInt{*S}
				sigVal.V = types.BigInt{*V}
				mx.EsigB = sigVal.GetMsgp()
				//v
				str = fmt.Sprintf("%s",mx.Hash.Hex())
				R,S,V = this.apos.commonTools.SIG([]byte(str))
				sigVal.R = types.BigInt{*R}
				sigVal.S = types.BigInt{*S}
				sigVal.V = types.BigInt{*V}
				mx.EsigV = sigVal.GetMsgp()


				this.msgOut<-mx

			}(this)

			return
		case data:=<-this.msgIn:
			func(this *step567CoinGenFlipBBA){
				this.lock.Lock()
				defer this.lock.Unlock()

				m6 := new(MCommon)
				*m6 = data.(MCommon)
				if m6 == nil{
					return
				}
				//add to IndexMap

				var subIndex *binaryStatus
				subIndex = this.allMxIndex[m6.Hash]
				if subIndex == nil {
					subIndex = makeBinaryStatus()
					this.allMxIndex[m6.Hash] = subIndex
				}
				//check sig
				//set status
				subIndex.setToStatus(*m6.Credential , int(m6.B))

			}(this)
			continue

		case <-this.exit:
			return
		}
	}
}












