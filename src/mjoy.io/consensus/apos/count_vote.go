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

type countVote struct {
	voteRecord  map[int]*stepVotes
	msgCh       chan *ByzantineAgreementStar
	stopCh      chan interface{}
	timer       *time.Timer
	timeStep    uint

}

func newCountVote() *countVote {
	cv := new(countVote)
	cv.init()
	return cv
}

func (cv *countVote) init() {
	cv.voteRecord = make(map[int]*stepVotes)
	cv.msgCh = make(chan *ByzantineAgreementStar, 1)
	cv.stopCh = make(chan interface{}, 1)
}

func (cv *countVote) run() {
	cv.timer = time.NewTimer(time.Second * 1)
	defer cv.timer.Stop()
	for {
		select {
		// receive message
		case voteMsg := <-cv.msgCh:
			fmt.Println(voteMsg)
		//timeout message
		case <-cv.timer.C:
			cv.timeStep++
			fmt.Println("timeout", cv.timeStep)
			cv.timer.Reset(time.Second * 1)
		case <-cv.stopCh:
			fmt.Println("countVote run exit", cv.timeStep)
			return
		}
	}
}
