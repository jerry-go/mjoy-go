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
// @File: steps.go
// @Date: 2018/06/15 11:35:15
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"mjoy.io/common/types"
	"time"
	"sync"
	"container/heap"
	"math/big"
)

//steps handle

//step 1:Block Proposal
type step1BlockProposal struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int            //which step the obj stay
	apos *Apos
	round *Round
	pCredential *CredentialSig
	lock sync.RWMutex

}
func makeStep1Obj(pApos *Apos , pCredential *CredentialSig  , step int)*step1BlockProposal{
	s := new(step1BlockProposal)
	s.apos = pApos

	s.msgIn = make(chan dataPack , 100)

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

func (this *step1BlockProposal)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()

	//new a M1 data
	m1 := new(M1)
	m1.Block = this.apos.makeEmptyBlockForTest()
	m1.Credential = this.pCredential
	m1.Esig = this.apos.commonTools.ESIG(m1.Block.Hash())
	//fill struct members
	//todo: should using interface
	this.apos.outMsger.SendInner(m1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Out M1",COLOR_SHORT_RESET)
	//todo:what should we dealing



}

//step 2:First step of GC
type step2FirstStepGC struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int
	smallestLBr *M1 //this node regard smallestLBr as the smallest credential's block info
	round *Round
	apos *Apos
	pCredential *CredentialSig
	lock sync.RWMutex

}

func makeStep2Obj(pApos *Apos , pCredential *CredentialSig , step int)*step2FirstStepGC{
	s := new(step2FirstStepGC)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)

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

func (this *step2FirstStepGC)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time

	delayT := time.Duration(Config().verifyDelay + Config().blockDelay)

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			//time to work now,send all
			if this.smallestLBr == nil {
				//mean we do not receive a M1 from now on
				return
			}
			m2 := new(M23)
			m2.Credential = this.pCredential
			m2.Hash = this.smallestLBr.Block.Hash()
			sigBytes := this.apos.commonTools.ESIG(m2.Hash)
			m2.Esig = append(m2.Esig , sigBytes...)

			this.apos.outMsger.SendInner(m2)
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"Out M2",COLOR_SHORT_RESET)
			return
		case data:=<-this.msgIn:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M1",COLOR_SHORT_RESET)
			m1 := new(M1)
			m1 = data.(*M1)

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
				this.apos.outMsger.PropagateMsg(m1)
			}

		case <-this.exit:
			return
		}
	}
}


//step 3:Second Step of GC
type step3SecondStepGC struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it

	exit chan int
	step int
	apos *Apos
	round *Round
	pCredential *CredentialSig
	//all M2 have received
	allM2Index map[types.Hash]map[CredentialSigForKey]bool
	lock sync.RWMutex

}

func makeStep3Obj(pApos *Apos , pCredential *CredentialSig , step int)*step3SecondStepGC{
	s := new(step3SecondStepGC)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)
	s.step = step
	s.pCredential = pCredential
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
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

func (this *step3SecondStepGC)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time
	delayT := time.Duration(3*Config().verifyDelay + Config().blockDelay)
	//log.Println("timeDelay:",3*Config().verifyDelay + Config().blockDelay)
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
				sigBytes := this.apos.commonTools.ESIG(m3.Hash)
				m3.Esig = append(m3.Esig , sigBytes...)
				this.apos.outMsger.SendInner(m3)
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"Out M3",COLOR_SHORT_RESET)
			}(this)

			return
		case data:=<-this.msgIn:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M2",COLOR_SHORT_RESET)
			func(this *step3SecondStepGC){
				this.lock.Lock()
				defer this.lock.Unlock()

				m2 := new(M23)
				m2 = data.(*M23)
				if m2 == nil{
					return
				}
				//add to IndexMap
				var subIndex map[CredentialSigForKey]bool
				subIndex = this.allM2Index[m2.Hash]
				if subIndex == nil {
					this.allM2Index[m2.Hash] = make(map[CredentialSigForKey]bool)
					subIndex = this.allM2Index[m2.Hash]
				}
				sigKey := *m2.Credential.ToCredentialSigKey()
				if _ , ok := subIndex[sigKey];!ok{
					subIndex[sigKey] = true
				}

			}(this)
			continue

		case <-this.exit:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)
			return
		}
	}
}


//step 4:First Step of BBA*
type step4FirstStepBBA struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int
	apos *Apos
	round *Round
	pCredential *CredentialSig
	//all M3 have received
	allM2Index map[types.Hash]map[CredentialSigForKey]bool
	lock sync.RWMutex
}

func makeStep4Obj(pApos *Apos , pCredential *CredentialSig  , step int)*step4FirstStepBBA{
	s := new(step4FirstStepBBA)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)
	s.step = step
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
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

func (this *step4FirstStepBBA)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time
	delayT := time.Duration(5*Config().verifyDelay + Config().blockDelay)

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
				//todo:should be H(Be)
				v := types.Hash{}
				g := 0
				if maxLen >= (2*total/3){
					v = maxHash
					g = 2
				}else if maxLen >= (1*total/3){
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

				//b big.Int
				h := types.BytesToHash(big.NewInt(int64(m4.B)).Bytes())
				sigBytes := this.apos.commonTools.ESIG(h)
				m4.EsigB = append(m4.EsigB,sigBytes...)
				//v
				sigBytes = this.apos.commonTools.ESIG(m4.Hash)
				m4.EsigV = append(m4.EsigV , sigBytes...)
				this.apos.outMsger.SendInner(m4)
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"Out M4",COLOR_SHORT_RESET)

			}(this)

			return
		case data:=<-this.msgIn:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M3",COLOR_SHORT_RESET)
			func(this *step4FirstStepBBA){
				this.lock.Lock()
				defer this.lock.Unlock()

				m3 := new(M23)
				m3 = data.(*M23)
				if m3 == nil{
					return
				}
				//add to IndexMap
				var subIndex map[CredentialSigForKey]bool
				subIndex = this.allM2Index[m3.Hash]
				if subIndex == nil {
					this.allM2Index[m3.Hash] = make(map[CredentialSigForKey]bool)
					subIndex = this.allM2Index[m3.Hash]
				}
				sigKey := *m3.Credential.ToCredentialSigKey()
				if _ , ok := subIndex[sigKey];!ok{
					subIndex[sigKey] = true
				}

			}(this)
			continue

		case <-this.exit:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)
			return
		}
	}
}

//merge 5 6 7

//step 5 6 7:Coin-Fixed-To-x step of BBA*
type step567CoinGenFlipBBA struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it

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

func makeStep567Obj(pApos *Apos , pCredential *CredentialSig  , step int)*step567CoinGenFlipBBA{
	s := new(step567CoinGenFlipBBA)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)
	s.allMxIndex = make(map[types.Hash]*binaryStatus)
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

func (this *step567CoinGenFlipBBA)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time

	delayT := time.Duration((2*this.step -3)*Config().verifyDelay + Config().blockDelay)

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
									*cHeap = append(*cHeap , &CredentialSigStatus{c:*c.ToCredentialSig() , v:0})
								}
								allCredential = bStatus.export1Credential()
								// 1 credential
								for _,c := range allCredential{
									*cHeap = append(*cHeap , &CredentialSigStatus{c:*c.ToCredentialSig() , v:1})
								}

							}
							heap.Init(cHeap)
							little := heap.Pop(cHeap).(CredentialSigStatus)

							mx.B = uint(little.v)
						}
					}
				}
				h := types.BytesToHash(big.NewInt(int64(mx.B)).Bytes())
				sigBytes := this.apos.commonTools.ESIG(h)
				mx.EsigB = append(mx.EsigB , sigBytes...)
				//v
				sigBytes = this.apos.commonTools.ESIG(mx.Hash)
				mx.EsigV = append(mx.EsigV , sigBytes...)

				this.apos.outMsger.SendInner(mx)
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"Out M",mx.Credential.Step.IntVal.Int64(),COLOR_SHORT_RESET)
			}(this)

			return
		case data:=<-this.msgIn:
			func(this *step567CoinGenFlipBBA){
				this.lock.Lock()
				defer this.lock.Unlock()

				m6 := new(MCommon)
				m6 = data.(*MCommon)
				if m6 == nil{
					return
				}
				//add to IndexMap
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M",m6.Credential.Step.IntVal.Int64(),COLOR_SHORT_RESET)
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
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)

			return
		}
	}
}

//m+3


//step 5 6 7:Coin-Fixed-To-x step of BBA*
type stepm3LastBBA struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int
	stepIndex int
	apos *Apos
	round *Round
	pCredential *CredentialSig
	lock sync.RWMutex

	//all M6 have received

}

func makeStepm3Obj(pApos *Apos , pCredential *CredentialSig , step int)*stepm3LastBBA{
	s := new(stepm3LastBBA)
	s.apos = pApos


	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)
	s.step = step


	s.pCredential = pCredential
	return s
}

func (this *stepm3LastBBA)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.round = pRound
	this.msgIn <- data
	return nil
}

func (this *stepm3LastBBA)stop(){
	this.exit<-1
}

func (this *stepm3LastBBA)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time
	delayT := time.Duration((2*Config().maxBBASteps + 3)*Config().verifyDelay + Config().blockDelay)

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			func(this *stepm3LastBBA){
				this.lock.Lock()
				defer this.lock.Unlock()


				m3 := new(MCommon)
				//todo:should be H(Be)
				m3.Hash = types.Hash{}
				m3.B = 1
				m3.Credential = this.pCredential
				h := types.BytesToHash(big.NewInt(int64(m3.B)).Bytes())
				sigBytes := this.apos.commonTools.ESIG(h)
				m3.EsigB = append(m3.EsigB , sigBytes...)
				//v
				sigBytes = this.apos.commonTools.ESIG(m3.Hash)
				m3.EsigV = append(m3.EsigV , sigBytes...)
				this.apos.outMsger.SendInner(m3)

			}(this)

			return
		case <-this.msgIn:

			continue

		case <-this.exit:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)
			return
		}
	}
}












