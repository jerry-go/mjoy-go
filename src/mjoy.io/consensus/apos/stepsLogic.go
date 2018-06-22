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
	"fmt"
)

//steps handle

//step 1:Block Proposal
type step1BlockProposalLogic struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int            //which step the obj stay
	lock sync.RWMutex

	stepCtx *StepContext

}
func makeStep1ObjLogic(step int , ctx *StepContext)*step1BlockProposalLogic{
	s := new(step1BlockProposalLogic)

	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)

	s.step = step

	s.stepCtx = ctx
	return s
}
func (this *step1BlockProposalLogic)sendMsg(data dataPack, pRound *Round)error{
	this.lock.Lock()
	defer this.lock.Unlock()
	this.stepCtx.SetRound(pRound)
	this.msgIn <- data
	//todo:
	return nil
}

func (this *step1BlockProposalLogic)stop(){
	this.exit<-1
}

func (this *step1BlockProposalLogic)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()

	//new a M1 data
	m1 := new(M1)
	m1.Block = this.stepCtx.makeEmptyBlockForTest()

	m1.Credential = this.stepCtx.GetCredential()
	m1.Esig = this.stepCtx.ESIG(m1.Block.Hash())
	//fill struct members
	//todo: should using interface
	this.stepCtx.SendInner(m1)
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Out M1",COLOR_SHORT_RESET)
	//todo:what should we dealing



}

//step 2:First step of GC
type step2FirstStepGCLogic struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int
	smallestLBr *M1 //this node regard smallestLBr as the smallest credential's block info
	lock sync.RWMutex
	stepCtx *StepContext
}

func makeStep2ObjLogic(step int , stepCtx *StepContext)*step2FirstStepGCLogic{
	s := new(step2FirstStepGCLogic)
	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)
	s.step = step
	s.stepCtx = stepCtx
	return s
}

func (this *step2FirstStepGCLogic)sendMsg(data dataPack , pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.stepCtx.SetRound(pRound)
	this.msgIn <- data
	return nil
}

func (this *step2FirstStepGCLogic)stop(){
	this.exit<-1
}

func (this *step2FirstStepGCLogic)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time

	delayT := time.Duration(Config().verifyDelay + Config().blockDelay)
	if LessTimeDelayFlag{
		delayT = time.Duration(LessTimeDelayCnt)
	}
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
			m2.Credential = this.stepCtx.GetCredential()
			m2.Hash = this.smallestLBr.Block.Hash()
			sigBytes := this.stepCtx.ESIG(m2.Hash)
			m2.Esig = append(m2.Esig , sigBytes...)

			this.stepCtx.SendInner(m2)
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"Out M2",COLOR_SHORT_RESET)
			return
		case data:=<-this.msgIn:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"In M1",COLOR_SHORT_RESET)
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
				this.stepCtx.PropagateMsg(m1)
			}

		case <-this.exit:
			return
		}
	}
}


//step 3:Second Step of GC
type step3SecondStepGCLogic struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it

	exit chan int
	step int
	//all M2 have received
	allM2Index map[types.Hash]map[CredentialSigForKey]bool
	lock sync.RWMutex
	stepCtx *StepContext

}

func makeStep3ObjLogic(step int , stepCtx *StepContext)*step3SecondStepGCLogic{
	s := new(step3SecondStepGCLogic)
	s.stepCtx = stepCtx


	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)
	s.step = step
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
	return s
}

func (this *step3SecondStepGCLogic)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.stepCtx.SetRound(pRound)
	this.msgIn <- data
	return nil
}

func (this *step3SecondStepGCLogic)stop(){
	this.exit<-1
}

func (this *step3SecondStepGCLogic)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time
	delayT := time.Duration(3*Config().verifyDelay + Config().blockDelay)
	if LessTimeDelayFlag{
		delayT = time.Duration(LessTimeDelayCnt)
	}
	//log.Println("timeDelay:",3*Config().verifyDelay + Config().blockDelay)
	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			func(this *step3SecondStepGCLogic){
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

				if maxLen * 3 > 2*total{
					v = maxHash
				}
				//pack m3 Data

				m3 := new(M23)
				m3.Credential = this.stepCtx.GetCredential()

				m3.Hash = v
				sigBytes := this.stepCtx.ESIG(m3.Hash)
				m3.Esig = append(m3.Esig , sigBytes...)
				this.stepCtx.SendInner(m3)
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"Out M3 ",v.String(),COLOR_SHORT_RESET)
			}(this)

			return
		case data:=<-this.msgIn:
			//logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M2",COLOR_SHORT_RESET)
			func(this *step3SecondStepGCLogic){
				this.lock.Lock()
				defer this.lock.Unlock()

				m2 := new(M23)
				m2 = data.(*M23)
				if m2 == nil{
					return
				}
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"In M2",m2.Hash.String(),COLOR_SHORT_RESET)
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
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)
			return
		}
	}
}


//step 4:First Step of BBA*
type step4FirstStepBBALogic struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int
	//all M3 have received
	allM2Index map[types.Hash]map[CredentialSigForKey]bool
	lock sync.RWMutex
	stepCtx *StepContext
}

func makeStep4ObjLogic(step int , stepCtx *StepContext)*step4FirstStepBBALogic{
	s := new(step4FirstStepBBALogic)


	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)
	s.step = step
	s.allM2Index = make(map[types.Hash]map[CredentialSigForKey]bool)
	s.stepCtx = stepCtx
	return s
}

func (this *step4FirstStepBBALogic)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.stepCtx.SetRound(pRound)
	this.msgIn <- data
	return nil
}

func (this *step4FirstStepBBALogic)stop(){
	this.exit<-1
}

func (this *step4FirstStepBBALogic)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time
	delayT := time.Duration(5*Config().verifyDelay + Config().blockDelay)
	if LessTimeDelayFlag{
		delayT = time.Duration(LessTimeDelayCnt)
	}

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			func(this *step4FirstStepBBALogic){
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
				if maxLen * 3 > 2 * total{
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"Step4  Do :maxLen * 3 > 2 * total,g=2",COLOR_SHORT_RESET)
					v = maxHash
					g = 2
				}else if maxLen * 3 > total{
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"Step4  Do :maxLen * 3 > total,g=1",COLOR_SHORT_RESET)
					v = maxHash
					g = 1
				}else{
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,"Step4  Do :Else,g=0",COLOR_SHORT_RESET)
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
				m4.Credential = this.stepCtx.GetCredential()

				//b big.Int
				h := types.BytesToHash(big.NewInt(int64(m4.B)).Bytes())
				sigBytes := this.stepCtx.ESIG(h)
				m4.EsigB = append(m4.EsigB,sigBytes...)
				//v
				sigBytes = this.stepCtx.ESIG(m4.Hash)
				m4.EsigV = append(m4.EsigV , sigBytes...)
				this.stepCtx.SendInner(m4)
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"Out M4",m4.Hash.String(),m4.B,COLOR_SHORT_RESET)

			}(this)

			return
		case data:=<-this.msgIn:
			func(this *step4FirstStepBBALogic){
				this.lock.Lock()
				defer this.lock.Unlock()

				m3 := new(M23)
				m3 = data.(*M23)
				if m3 == nil{
					return
				}
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"In M3",m3.Hash.String(),COLOR_SHORT_RESET)
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
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)
			return
		}
	}
}

//merge 5 6 7

//step 5 6 7:Coin-Fixed-To-x step of BBA*
type step567CoinGenFlipBBALogic struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it

	exit chan int
	step int
	stepIndex int
	lock sync.RWMutex

	//all M6 have received
	allMxIndex map[types.Hash]*binaryStatus
	stepCtx *StepContext
}

func makeStep567ObjLogic(step int , stepCtx *StepContext)*step567CoinGenFlipBBALogic{
	s := new(step567CoinGenFlipBBALogic)


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

	s.stepCtx = stepCtx
	return s
}

func (this *step567CoinGenFlipBBALogic)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.stepCtx.SetRound(pRound)

	this.msgIn <- data
	return nil
}

func (this *step567CoinGenFlipBBALogic)stop(){
	this.exit<-1
}

func (this *step567CoinGenFlipBBALogic)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time

	delayT := time.Duration((2*this.step -3)*Config().verifyDelay + Config().blockDelay)
	if LessTimeDelayFlag{
		delayT = time.Duration(LessTimeDelayCnt)
	}

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			func(this *step567CoinGenFlipBBALogic){
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
				if max0Len * 3 > 2 * total {
					mx.Hash = maxHash
					mx.Credential = this.stepCtx.GetCredential()
					mx.B = 0
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
						"StepCommon  Do :max0Len * 3 > 2 * total,B=0",
						COLOR_SHORT_RESET)
				}else if max1Len * 3 > 2 * total {
					mx.Hash = maxHash
					mx.Credential = this.stepCtx.GetCredential()
					mx.B = 1
					logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
						"StepCommon  Do :max0Len * 3 > 2 * total,B=1",
						COLOR_SHORT_RESET)
				}else {
					mx.Hash = maxHash
					mx.Credential = this.stepCtx.GetCredential()
					switch this.stepIndex {
					case 5:
						mx.B = 0
						logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
							"StepCommon  Do :else %5,B=0",
							COLOR_SHORT_RESET)
					case 6:
						mx.B = 1
						logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
							"StepCommon  Do :else %6,B=1",
							COLOR_SHORT_RESET)
					case 7:
						{
							cHeap := new(CredentialSigStatusHeap)
							allCnt := 0
							for _,bStatus := range this.allMxIndex{
								allCredential := bStatus.export0Credential()
								// 0 credential
								for _,c := range allCredential{
									allCnt++
									*cHeap = append(*cHeap , &CredentialSigStatus{c:*c.ToCredentialSig() , v:0})
								}
								allCredential = bStatus.export1Credential()
								// 1 credential
								for _,c := range allCredential{
									allCnt++
									*cHeap = append(*cHeap , &CredentialSigStatus{c:*c.ToCredentialSig() , v:1})
								}

							}
							fmt.Println("...................All Mx Index:" , allCnt)
							heap.Init(cHeap)
							little := heap.Pop(cHeap).(*CredentialSigStatus)

							mx.B = uint(little.v)

							logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX,
								"StepCommon  Do :else %7,B=",mx.B,
								COLOR_SHORT_RESET)
						}
					}
				}
				h := types.BytesToHash(big.NewInt(int64(mx.B)).Bytes())
				sigBytes := this.stepCtx.ESIG(h)
				mx.EsigB = append(mx.EsigB , sigBytes...)
				//v
				sigBytes = this.stepCtx.ESIG(mx.Hash)
				mx.EsigV = append(mx.EsigV , sigBytes...)

				this.stepCtx.SendInner(mx)
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"Out M",mx.Credential.Step.IntVal.Int64(),COLOR_SHORT_RESET)
			}(this)

			return
		case data:=<-this.msgIn:
			func(this *step567CoinGenFlipBBALogic){
				this.lock.Lock()
				defer this.lock.Unlock()

				m6 := new(MCommon)
				m6 = data.(*MCommon)
				if m6 == nil{
					return
				}
				//add to IndexMap
				//logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.pCredential.Step.IntVal.Int64(),"In M",m6.Credential.Step.IntVal.Int64(),COLOR_SHORT_RESET)
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"In M3 Hash:",m6.Hash.String(),"B:",m6.B,COLOR_SHORT_RESET)
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
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)

			return
		}
	}
}

//m+3


//step 5 6 7:Coin-Fixed-To-x step of BBA*
type stepm3LastBBALogic struct {
	msgIn chan dataPack   //Data: Out ---- > In , we should create it
	exit chan int
	step int
	stepIndex int
	lock sync.RWMutex
	stepCtx *StepContext
	//all M6 have received

}

func makeStepm3ObjLogic(step int , stepCtx *StepContext)*stepm3LastBBALogic{
	s := new(stepm3LastBBALogic)


	s.msgIn = make(chan dataPack , 100)

	s.exit = make(chan int , 1)
	s.step = step

	s.stepCtx = stepCtx
	return s
}

func (this *stepm3LastBBALogic)sendMsg(data dataPack, pRound *Round)error{
	//todo:
	this.lock.Lock()
	defer this.lock.Unlock()
	this.stepCtx.SetRound(pRound)
	this.msgIn <- data
	return nil
}

func (this *stepm3LastBBALogic)stop(){
	this.exit<-1
}

func (this *stepm3LastBBALogic)run(wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()
	//this step ,we should wait the time
	delayT := time.Duration((2*Config().maxBBASteps + 3)*Config().verifyDelay + Config().blockDelay)
	if LessTimeDelayFlag{
		delayT = time.Duration(LessTimeDelayCnt)
	}

	timerT := time.Tick(delayT*time.Second)
	for{
		select {
		case <-timerT:
			func(this *stepm3LastBBALogic){
				this.lock.Lock()
				defer this.lock.Unlock()


				m3 := new(MCommon)
				//todo:should be H(Be)
				m3.Hash = types.Hash{}
				m3.B = 1
				m3.Credential = this.stepCtx.GetCredential()
				h := types.BytesToHash(big.NewInt(int64(m3.B)).Bytes())
				sigBytes := this.stepCtx.ESIG(h)
				m3.EsigB = append(m3.EsigB , sigBytes...)
				//v
				sigBytes = this.stepCtx.ESIG(m3.Hash)
				m3.EsigV = append(m3.EsigV , sigBytes...)
				this.stepCtx.SendInner(m3)

			}(this)

			return
		case <-this.msgIn:

			continue

		case <-this.exit:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX,"[A]Step:",this.stepCtx.GetCredential().Step.IntVal.Int64(),"  Exit",COLOR_SHORT_RESET)
			return
		}
	}
}












