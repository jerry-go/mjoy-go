package apos

import (
	"sync"
	"time"
	"sort"
	"mjoy.io/common/types"
)

type BpWithPriority struct {
	j int   //the priofity
	bp *BlockProposal
}

type BpWithPriorityHeap []*BpWithPriority

func (h BpWithPriorityHeap)Len() int            {return len(h)}
func (h BpWithPriorityHeap)Less(i , j int)bool  {return h[i].j > h[j].j}
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
	existMap map[types.Hash]bool	//todo:here should be change to *bp
	msgChan chan *BlockProposal
	exit    chan interface{}
	ctx     *stepCtx
	priorityBp *BlockProposal
}

func makeBpObj(ctx *stepCtx) *BpObj {
	s := new(BpObj)
	s.ctx = ctx
	s.BpHeap = make(BpWithPriorityHeap , 0)
	s.msgChan = make(chan *BlockProposal)
	s.existMap = make(map[types.Hash]bool)
	s.exit = make(chan interface{})

	return s
}

func (this *BpObj)isExistBlock(blockHash types.Hash)bool{
	this.lock.RLock()
	defer this.lock.RUnlock()

	if _,ok := this.existMap[blockHash];ok{
		return true
	}
	return false
}

func (this *BpObj)addExistBlock(blockHash types.Hash){
	this.lock.Lock()
	defer this.lock.Unlock()

	this.existMap[blockHash] = true
}

func (this *BpObj)CommitteeVote(data *VoteData){

	cret := this.ctx.getCredentialByStep(uint64(data.Step))
	if cret == nil {
		return
	}

	//todo :need pack ba msg

	msgBa := newByzantineAgreementStar()
	//hash
	msgBa.Hash = data.Value
	//Credential
	msgBa.Credential = cret

	//Esig
	msgBa.Esig.round = msgBa.Credential.Round
	msgBa.Esig.step = msgBa.Credential.Step
	msgBa.Esig.val = make([]byte , 0)
	msgBa.Esig.val = append(msgBa.Esig.val , msgBa.Hash.Bytes()...)

	err := this.ctx.esig(msgBa.Esig)
	if err != nil {
		logger.Error("CommitteeVote Esig Err:" , err.Error())
		return
	}


	if cret.votes > 0{
		this.ctx.sendInner(msgBa)
	}
}

func (this *BpObj)makeBlock(){
	bp := newBlockProposal()
	bp.Credential = this.ctx.getCredentialByStep(StepBp)
	if nil == bp.Credential {
		logger.Warn("makeBlock getCredentialByStep--->nil")
		return
	}
	//bcd := &block.ConsensusData{}
	//bcd.Id = ConsensusDataId
	//bcd.Para = bp.Credential.Signature.toBytes()
	bcd := this.ctx.makeBlockConsensusData(bp)

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
	tProposer := int(Config().tProposer)
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
				vd.Step = StepReduction1
				vd.Value = x.bp.Block.Hash()

				this.CommitteeVote(vd)

				this.ctx.startVoteTimer(int(Config().delayStep))
				this.ctx.setBpResult(x.bp.Block.Hash())
				//todo:inform the reduction
				return
			}
		case bp := <-this.msgChan:

				//logic do
				//verify the block
				if !this.ctx.verifyBlock(bp.Block){
					continue
				}
				//check is exist a same block
				if this.isExistBlock(bp.Block.Hash()) {
					continue
				}

				//check the node has the right to produce a block
				pri := bp.Credential.votes

				bpp := new(BpWithPriority)
				bpp.j = int(pri)
				bpp.bp = bp

				this.BpHeap.Push(bpp)
				this.addExistBlock(bp.Block.Hash())
				if this.priorityBp == nil {
					this.ctx.propagateMsg(bp)
				} else if pri > this.priorityBp.Credential.votes {
					this.priorityBp = bp
					this.ctx.propagateMsg(bp)
				}

				if this.BpHeap.Len() > tProposer {
					sort.Sort(&this.BpHeap)
					//get the bigger one
					x := this.BpHeap[0]
					_ = x


					vd := new(VoteData)
					vd.Round = x.bp.Credential.Round
					vd.Step = StepReduction1
					vd.Value = x.bp.Block.Hash()

					this.CommitteeVote(vd)
					this.ctx.startVoteTimer(int(Config().delayStep))
					this.ctx.setBpResult(x.bp.Block.Hash())
					//todo:inform the reduction

					return
				}
		}
	}
}

func (this *BpObj) sendMsg(bp *BlockProposal) {
	this.msgChan <- bp
}
