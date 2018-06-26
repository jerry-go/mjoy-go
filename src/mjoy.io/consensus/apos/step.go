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
	"time"
	"sync"
)

type step interface {
	setCtx(ctx *stepCtx)         // set the context of step
	getTTL() time.Duration       // get the ttl of step
	timerHandle()
	dataHandle(data interface{})
	stopHandle()
}

// the routine of step
type stepRoutine struct {
	inputCh chan interface{}
	stopCh  chan struct{}
	timer   *time.Timer
	s       step
	wg     *sync.WaitGroup
}

func newStepRoutine() *stepRoutine {
	return &stepRoutine{
		make(chan interface{}),
		make(chan struct{}),
		nil,
		nil,
		&sync.WaitGroup{},
	}
}

func (sr *stepRoutine) reset() {
	sr.inputCh = make(chan interface{})
	sr.stopCh = make(chan struct{})
	sr.timer = nil
	sr.s = nil
	sr.wg =  &sync.WaitGroup{}
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
		sr.timer = time.NewTimer(0)
		<-sr.timer.C
		defer sr.timer.Stop()
		if sr.s.getTTL() != 0 {
			sr.timer.Reset(sr.s.getTTL() * time.Second)
		}
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

type stepCtx struct {
	getStep   func() int	// get the number of step in the round
	stopStep  func()        // stop the step
	stopRound func()		// stop all the step in the round, and end the round
	getCredential func() Credential
	getEphemeralSig func(signed []byte) Signature
}



