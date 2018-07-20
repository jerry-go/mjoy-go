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
// @File: count_vote.go
// @Date: 2018/07/19 13:34:19
////////////////////////////////////////////////////////////////////////////////

package apos

import (
	"mjoy.io/common/types"
	"fmt"
	"time"
)

type stepVotes struct {
	counts      map[types.Hash]uint
	//flag for vote result
	isFinish    bool
	value       types.Hash
}

func newStepVotes() *stepVotes {
	sv := new(stepVotes)
	sv.counts = make(map[types.Hash]uint)
	return sv
}

type countVote struct {
	voteRecord  map[int]*stepVotes
	msgCh       chan *ByzantineAgreementStar
	stopCh      chan interface{}
	timer       *time.Timer
	timerStep   uint
	emptyBlock  types.Hash

	sendVoteResult func(s int, hash types.Hash)


}

func newCountVote(sendVoteResult func(s int, hash types.Hash), emptyBlock  types.Hash) *countVote {
	cv := new(countVote)
	cv.init()
	cv.sendVoteResult = sendVoteResult
	cv.emptyBlock = emptyBlock
	return cv
}

func (cv *countVote) init() {
	cv.voteRecord = make(map[int]*stepVotes)
	cv.msgCh = make(chan *ByzantineAgreementStar, 1)
	cv.stopCh = make(chan interface{}, 1)
}

//this function should be called by BP handle
func (cv *countVote) startTimer(delay int) {
	delayDuration := time.Second * time.Duration(delay)
	cv.timer = time.NewTimer(delayDuration)
	cv.timerStep = STEP_REDUCTION_1
}

func (cv *countVote) run() {
	for {
		select {
		// receive message
		case voteMsg := <-cv.msgCh:
			step, hash, complete := cv.processMsg(voteMsg)
			if complete {
				cv.countSuccess(step, hash)
			}
		//timeout message
		case <-cv.timer.C:
			fmt.Println("timeout", cv.timerStep)
			cv.timeoutHandle()
		case <-cv.stopCh:
			fmt.Println("countVote run exit", cv.timerStep)
			cv.timer.Stop()
			return
		}
	}
}

func (cv *countVote) timeoutHandle() {
	timeoutStep := cv.timerStep
	resetTimer := true
	switch {
	case timeoutStep == STEP_REDUCTION_1:
		cv.timerStep = STEP_REDUCTION_2
	case timeoutStep == STEP_REDUCTION_2:
		cv.timerStep = 1
	case timeoutStep <= Config().maxStep:
		cv.timerStep++
	default:
		//ignore
		resetTimer = false
	}

	delay := time.Second * time.Duration(Config().delayStep)
	if resetTimer {
		cv.timer.Reset(delay)
	}
	cv.sendVoteResult(int(timeoutStep), TimeOut)

}

func (cv *countVote) countSuccess(step int, hash types.Hash) {

	delay := time.Second * time.Duration(Config().delayStep)
	cv.timer.Reset(delay)
	cv.sendVoteResult(step, hash)

	switch {
	case step == STEP_REDUCTION_1:
		cv.timerStep = STEP_REDUCTION_2
	case step == STEP_REDUCTION_2:
		cv.timerStep = 1
	case step <= int(Config().maxStep):
		bbaIdex := step % 3
		if bbaIdex == 1 && hash != cv.emptyBlock {
			//bba complete: block hash
			cv.timerStep = STEP_FINAL
		} else if bbaIdex == 2 && hash == cv.emptyBlock {
			//bba complete: empty block hash
			cv.timerStep = STEP_FINAL
		}
		cv.timerStep++
	default:
		cv.timerStep = STEP_IDLE
	}
}

func (cv *countVote) addVotes(ba *ByzantineAgreementStar) (types.Hash, uint) {
	hash := ba.Hash
	step := int(ba.Credential.Step)
	votes := ba.Credential.votes
	if sv, ok:= cv.voteRecord[step]; !ok {
		svNew := newStepVotes()
		cv.voteRecord[step] = svNew
		svNew.counts[hash] = votes
		return hash, votes
	} else {
		if hashVote, ok := sv.counts[hash]; ok {
			hashVote += votes
			return hash, hashVote
		} else {
			sv.counts[hash] = votes
			return hash, votes
		}
	}
}

func (cv *countVote) processMsg(ba *ByzantineAgreementStar) (int, types.Hash, bool) {
	step := int(ba.Credential.Step)

	hash, votes := cv.addVotes(ba)

	sv := cv.voteRecord[step]

	//check this step whether is finish
	if sv.isFinish {
		return step, types.Hash{}, false
	}

	if votes > getThreshold(step) {
		sv.isFinish = true
		return step, hash, true
	}
	return step, hash, false
}
