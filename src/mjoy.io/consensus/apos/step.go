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
// @File: step.go
// @Date: 2018/06/21 10:38:21
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"mjoy.io/common/types"
	"mjoy.io/core/blockchain/block"
	"sync"
	"time"
)

type step interface {
	setCtx(ctx *stepCtx)         // set the context of step
	getTTL() time.Duration       // get the ttl of step
	timerHandle()                //timerout handle when time's up
	dataHandle(data interface{}) //data handle when data in
	stopHandle()                 //deal the last work
}

// the routine of step
type stepRoutine struct {
	inputCh chan interface{}
	stopCh  chan interface{}
	timer   *time.Timer
	s       step
	wg      *sync.WaitGroup
}

func newStepRoutine() *stepRoutine {
	return &stepRoutine{
		make(chan interface{}),
		make(chan interface{}),
		nil,
		nil,
		&sync.WaitGroup{},
	}
}

func (sr *stepRoutine) setStep(s step) {
	sr.s = s
}

func (sr *stepRoutine) reset() {
	sr.inputCh = make(chan interface{})
	sr.stopCh = make(chan interface{})
	sr.timer = nil
	sr.s = nil
	sr.wg = &sync.WaitGroup{}
}

//func (sr *stepRoutine)sendMsg(dataPack,*Round) error{
//
//}
func (sr *stepRoutine) sendMsg(dp dataPack) error {
	sr.inputCh <- dp
	return nil
}

// run the routine of step
func (sr *stepRoutine) run(s step) {
	sr.s = s

	run := func() {
		sr.wg.Add(1)
		defer func() {
			sr.wg.Done()
			sr.reset()
		}()

		// start timer
		timeDelay := sr.s.getTTL()
		sr.timer = time.NewTimer(timeDelay)
		defer sr.timer.Stop()

		for {
			select {
			case data := <-sr.inputCh:
				sr.s.dataHandle(data)

			case <-sr.timer.C:
				sr.s.timerHandle()

			case <-sr.stopCh:
				sr.s.stopHandle()
				return
			}
		}
	}

	// go routine
	go run()
}

// stop routine and wait until the routine is closed
func (sr *stepRoutine) stop() {
	close(sr.stopCh)
	sr.wg.Wait()
}

//stepCtx contains all functions the stepObj will use
type stepCtx struct {
	getStep   func() int // get the number of step in the round
	getRound  func() int
	stopStep  func() // stop the step
	stopRound func() // stop all the step in the round, and end the round

	//getCredential func() signature
	//getEphemeralSig func(signed []byte) signature
	esig                  func(pEphemeralSign *EphemeralSign) error
	sendInner             func(pack dataPack) error
	propagateMsg          func(dataPack) error
	getCredential         func() *CredentialSign
	setRound              func(*Round)
	makeEmptyBlockForTest func(cs *CredentialSign) *block.Block
	getEmptyBlockHash     func() types.Hash
	getEphemeralSig       func(signed []byte) Signature
	getProducerNewBlock   func(data *block.ConsensusData) *block.Block
	//getPrivKey

	//gilad
	commonCoin func(round , step , t uint64)uint64  //x
	writeRet func(data *VoteData)                   //x
	sortition func(hash types.Hash , t,w,W uint64)uint64
	verifyBlock func(b *block.Block)bool
	verifySort func(cret CredentialSign , w, W,t uint64)uint64
	getCredentialByStep   func(step uint64)*CredentialSign
	getAccountMonney func (address types.Address , round uint64)uint64
	getTotalMonney func(round uint64)uint64
	getBpThreshold func()uint64
	getVoteThreshold func()uint64

	startVoteTimer func(delay int)
	makeBlockConsensusData func(bp *BlockProposal) *block.ConsensusData

	setBpResult func(hash types.Hash)
	setReductionResult func(hash types.Hash)
	setBbaResult  func(hash types.Hash)
	setFinalResult func(hash types.Hash)



}
