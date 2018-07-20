package apos

import (
	"mjoy.io/common/types"
	"sync"
	"time"
)

type BpObj struct {
	lock    sync.RWMutex
	allBp   map[types.Hash]*BlockProposal
	msgChan chan *BlockProposal
	ctx     *stepCtx
}

func makeStepBp(ctx *stepCtx) *BpObj {
	s := new(BpObj)
	s.ctx = ctx
	s.allBp = make(map[types.Hash]*BlockProposal)
	s.msgChan = make(chan *BlockProposal)

	return s
}

func (this *BpObj) run() {
	timer := time.Tick(60 * time.Second)

	for {
		select {
		case <-timer:
			//todo :should
		case <-this.msgChan:
				//logic do
		}
	}
}
