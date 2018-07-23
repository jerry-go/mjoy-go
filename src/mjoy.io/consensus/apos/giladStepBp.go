package apos

import (
	"sync"
	"time"
	"sort"
	"mjoy.io/core/blockchain/block"
)

type BpWithPriority struct {
	j int   //the priofity
	bp *BlockProposal
}

type BpWithPriorityHeap []*BpWithPriority

func (h BpWithPriorityHeap)Len() int            {return len(h)}
func (h BpWithPriorityHeap)Less(i , j int)bool  {return h[i].j < h[j].j}
func (h BpWithPriorityHeap)Swap(i , j int)      {h[i],h[j] = h[j],h[i]}

func (h *BpWithPriorityHeap)Push(x interface{}){
	*h = append(*h , x.(*BpWithPriority))
}

func (h *BpWithPriorityHeap)Pop()interface{}{
	old := *h
	n := len(old)
	x := old[n -1]
	*h = old[0:n-1]
	return x
}


type BpObj struct {
	lock    sync.RWMutex
	BpHeap  BpWithPriorityHeap
	msgChan chan *BlockProposal
	exit    chan interface{}
	ctx     *stepCtx
}

func makeBpObj(ctx *stepCtx) *BpObj {
	s := new(BpObj)
	s.ctx = ctx
	s.BpHeap = make(BpWithPriorityHeap , 0)
	s.msgChan = make(chan *BlockProposal)
	s.exit = make(chan interface{})

	return s
}

func (this *BpObj)makeBlock(){
	bp := newBlockProposal()
	bp.Credential = this.ctx.getCredentialByStep(0)

	bcd := &block.ConsensusData{}
	bcd.Id = ConsensusDataId
	bcd.Para = bp.Credential.Signature.toBytes()

	bp.Block = this.ctx.getProducerNewBlock(bcd)

	bp.Esig.round = bp.Credential.Round
	bp.Esig.step = StepBp
	bp.Esig.val = make([]byte , 0)
	h := bp.Block.Hash()
	bp.Esig.val = append(bp.Esig.val , h[:]...)
	err := this.ctx.esig(bp.Esig)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//delay here ?

	this.ctx.sendInner(bp)

	logger.Debug(COLOR_PREFIX+COLOR_FRONT_PINK+COLOR_SUFFIX, "***[A]Out M1 CreHash:", bp.Credential.Signature.Hash().String(), " BlockHash", bp.Block.B_header.Hash().String(), COLOR_SHORT_RESET)
}

func (this *BpObj) run() {
	//make block and send out
	go this.makeBlock()

	timer := time.Tick(60 * time.Second)

	for {
		select {
		case <-this.exit:
			return
		case <-timer:

			if this.BpHeap.Len() == 0 {
				//specialdo
			}else{

				sort.Sort(&this.BpHeap)
				x := this.BpHeap[0]

				//make reduction input data
				vd := new(VoteData)
				vd.Round = x.bp.Credential.Round
				vd.Step = StepBp
				vd.Value = x.bp.Block.Hash()

				this.ctx.sendInner(vd)
				this.ctx.startVoteTimer(int(Config().delayStep))

				//todo:inform the reduction
				return
			}
		case bp := <-this.msgChan:

				//logic do
				//verify the block
				if !this.ctx.verifyBlock(bp.Block){
					continue
				}
				//get the priority
				sender ,err  := bp.Credential.sender()
				if err != nil{
					logger.Error("bp")
					continue
				}
				w := this.ctx.getAccountMonney(sender , bp.Credential.Round - 1)
				W := this.ctx.getTotalMonney(bp.Credential.Round -1 )
				t := this.ctx.getBpThreshold()

				//check the node has the right to produce a block
				pri := this.ctx.verifySort(*bp.Credential , w , W  , t)
				if pri > 0{
					bpp := new(BpWithPriority)
					this.BpHeap.Push(bpp)
					if this.BpHeap.Len() > 26 {
						sort.Sort(&this.BpHeap)
						//get the bigger one
						x := this.BpHeap[0]
						_ = x

						vd := new(VoteData)
						vd.Round = x.bp.Credential.Round
						vd.Step = StepBp
						vd.Value = x.bp.Block.Hash()

						this.ctx.sendInner(vd)
						this.ctx.startVoteTimer(int(Config().delayStep))
						//todo:inform the reduction

						return
					}
				}

		}
	}
}
