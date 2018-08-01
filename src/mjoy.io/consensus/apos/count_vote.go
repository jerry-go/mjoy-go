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
	"time"
)

type stepVotes struct {
	counts      map[types.Hash]float64
	//flag for vote result
	isFinish    bool
	value       types.Hash
}

func newStepVotes() *stepVotes {
	sv := new(stepVotes)
	sv.counts = make(map[types.Hash]float64)
	return sv
}

type countVote struct {
	voteRecord  map[int]*stepVotes
	msgCh       chan *ByzantineAgreementStar
	stopCh      chan interface{}
	timer       *time.Timer
	timerStep   uint
	emptyBlock  types.Hash

	bbaFinish   bool

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
	cv.timerStep = STEP_IDLE
	cv.timer = time.NewTimer(0)
}

//this function should be called by BP handle
func (cv *countVote) startTimer(delay int) {
	delayDuration := time.Second * time.Duration(delay)
	cv.timer.Reset(delayDuration)
	cv.timerStep = STEP_REDUCTION_1
}

func (cv *countVote) run() {
	for {
		select {
		// receive message
		case voteMsg := <-cv.msgCh:
			step, hash, complete := cv.processMsg(voteMsg)
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX  , "CountVote Get Msg:" , "step:" , step , "  hash:" , hash.Hex() , COLOR_SHORT_RESET)
			if complete {
				logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "CountVote Complete Step:" , step , COLOR_SHORT_RESET)
				cv.countSuccess(step, hash)
			}
		//timeout message
		case <-cv.timer.C:
			logger.Debug("countVote timeout, step", cv.timerStep)
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "countVote timeout, step", cv.timerStep , COLOR_SHORT_RESET)
			cv.timeoutHandle()
		case <-cv.stopCh:
			logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "countVote run exit", cv.timerStep, COLOR_SHORT_RESET)
			cv.timer.Stop()
			return
		}
	}
}

func (cv *countVote) getNextTimerStep(step int) int {
	timeoutStep := step
	for {
		switch {
		case timeoutStep == STEP_REDUCTION_1:
			timeoutStep = STEP_REDUCTION_2
		case timeoutStep == STEP_REDUCTION_2:
			timeoutStep = 1
		case timeoutStep < int(Config().maxStep):
			timeoutStep++
			if timeoutStep == int(Config().maxStep) {
				timeoutStep = STEP_IDLE
			}
		case timeoutStep == STEP_FINAL:
			timeoutStep = STEP_IDLE
		default:
			//ignore
			timeoutStep = STEP_IDLE
		}
		if timeoutStep < int(Config().maxStep) && cv.bbaFinish{
			timeoutStep = STEP_FINAL
		}
		if sv, ok:= cv.voteRecord[timeoutStep]; ok {
			if sv.isFinish == true {
				continue
			}
		}
		break
	}
	cv.timerStep = uint(timeoutStep)
	return timeoutStep
}

func (cv *countVote) timeoutHandle() {
	timeoutStep := int(cv.timerStep)

	if timeoutStep == STEP_IDLE {
		logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "timeoutHandle timeoutStep == STEP_IDLE , Step", timeoutStep, COLOR_SHORT_RESET)
		//ignore
		return
	}

	logger.Debug(COLOR_PREFIX+COLOR_FRONT_RED+COLOR_SUFFIX , "timeoutHandle  Step:", timeoutStep, COLOR_SHORT_RESET)
	//fill results in voteRecord
	if sv, ok:= cv.voteRecord[timeoutStep]; !ok {
		svNew := newStepVotes()
		cv.voteRecord[timeoutStep] = svNew
		svNew.isFinish = true
		svNew.value = TimeOut
	} else {
		sv.isFinish = true
		sv.value = TimeOut
	}

	resetTimer := true
	nextTimoutStep := cv.getNextTimerStep(timeoutStep)
	if nextTimoutStep == STEP_IDLE {
		resetTimer = false
	}

	if resetTimer {
		delay := time.Second * time.Duration(Config().delayStep)
		cv.timer.Reset(delay)
	}

	cv.sendVoteResult(timeoutStep, TimeOut)

}

func (cv *countVote) countSuccess(step int, hash types.Hash) {
	//send result
	cv.sendVoteResult(step, hash)

	resetTimer := false
	nextTimoutStep := STEP_IDLE
	if int(cv.timerStep) == step {
		//reset timer
		resetTimer = true
		nextTimoutStep = cv.getNextTimerStep(step)
	}

	if step < int(Config().maxStep) {
		bbaIdex := step % 3
		if bbaIdex == 1 && hash != cv.emptyBlock {
			//bba complete: block hash
			cv.bbaFinish = true
			if cv.timerStep < (Config().maxStep) || cv.timerStep == STEP_IDLE {
				nextTimoutStep = STEP_FINAL
				cv.timerStep = STEP_FINAL
				resetTimer = true
			}
		} else if bbaIdex == 2 && hash == cv.emptyBlock {
			//bba complete: empty block hash
			cv.bbaFinish = true
			if cv.timerStep < (Config().maxStep) || cv.timerStep == STEP_IDLE {
				nextTimoutStep = STEP_FINAL
				cv.timerStep = STEP_FINAL
				resetTimer = true
			}
		}
	}

	if nextTimoutStep == STEP_IDLE {
		resetTimer = false
	}

	if resetTimer {
		delay := time.Second * time.Duration(Config().delayStep)
		cv.timer.Reset(delay)
	}
}

func (cv *countVote) addVotes(ba *ByzantineAgreementStar) (types.Hash, float64) {
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
			sumVote := hashVote + votes
			sv.counts[hash] = sumVote
			return hash, sumVote
		} else {
			sv.counts[hash] = votes
			return hash, votes
		}
	}
}

func (cv *countVote) processMsg(ba *ByzantineAgreementStar) (int, types.Hash, bool) {
	logger.Debug("processMsg", ba.Credential.Step, ba.Hash.String())
	step := int(ba.Credential.Step)
	if step < int(Config().maxStep) && cv.bbaFinish {
		logger.Info("all bba finished. step ", step, "will ignore")
		return step, types.Hash{}, false
	}
	sv, ok := cv.voteRecord[step]
	if ok {
		//check this step whether is finish
		if sv.isFinish {
			logger.Info("step", step, "is finished, ignore vote")
			return step, types.Hash{}, false
		}
	}

	hash, votes := cv.addVotes(ba)

	sv = cv.voteRecord[step]
	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX ,"ProcessMsg getThreshold:" , getThreshold(step) , "  Now Votes:",votes)
	if votes > float64(getThreshold(step)) {
		sv.isFinish = true
		return step, hash, true
	}
	return step, hash, false
}

func (cv *countVote) sendMsg(ba *ByzantineAgreementStar) {
	cv.msgCh<- ba
}

func (cv *countVote) stop() {
	cv.stopCh<- 1
}
